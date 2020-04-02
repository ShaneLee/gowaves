package miner

import (
	"errors"

	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/services"
	"github.com/wavesplatform/gowaves/pkg/state"
	"github.com/wavesplatform/gowaves/pkg/types"
	"go.uber.org/zap"
)

var NoTransactionsErr = errors.New("no transactions")
var StateChangedErr = errors.New("state changed")

type MicroMiner struct {
	state  state.State
	utx    types.UtxPool
	scheme proto.Scheme
}

func NewMicroMiner(services services.Services) *MicroMiner {
	return &MicroMiner{
		state:  services.State,
		utx:    services.UtxPool,
		scheme: services.Scheme,
	}
}

func (a *MicroMiner) Micro(
	minedBlock *proto.Block,
	rest proto.MiningLimits,
	keyPair proto.KeyPair) (*proto.Block, *proto.MicroBlock, proto.MiningLimits, error) {

	// way to stop mine microblocks
	if minedBlock == nil {
		return nil, nil, rest, errors.New("no block provided")
	}

	topBlock := a.state.TopBlock()
	if topBlock.BlockSignature != minedBlock.BlockSignature {
		// block changed, exit
		return nil, nil, rest, StateChangedErr
	}

	height, err := a.state.Height()
	if err != nil {
		return nil, nil, rest, err
	}

	parentTimestamp := topBlock.Timestamp
	if height > 1 {
		parent, err := a.state.BlockByHeight(height - 1)
		if err != nil {
			return nil, nil, rest, err
		}
		parentTimestamp = parent.Timestamp
	}

	//
	transactions := make([]proto.Transaction, 0)
	cnt := 0
	binSize := 0

	var unAppliedTransactions []*types.TransactionWithBytes

	// 255 is max transactions count in microblock
	for i := 0; i < 255; i++ {
		t := a.utx.Pop()
		if t == nil {
			break
		}
		binTr := t.B
		transactionLenBytes := 4
		if binSize+len(binTr)+transactionLenBytes > rest.MaxTxsSizeInBytes {
			unAppliedTransactions = append(unAppliedTransactions, t)
			continue
		}

		err = a.state.ValidateNextTx(t.T, minedBlock.Timestamp, parentTimestamp, minedBlock.Version)
		if err != nil {
			unAppliedTransactions = append(unAppliedTransactions, t)
			continue
		}

		cnt += 1
		binSize += len(binTr) + transactionLenBytes
		transactions = append(transactions, t.T)
	}

	a.state.ResetValidationList()

	// return unapplied transactions
	for _, unapplied := range unAppliedTransactions {
		_ = a.utx.AddWithBytes(unapplied.T, unapplied.B)
	}

	// no transactions applied, skip
	if cnt == 0 {
		return nil, nil, rest, NoTransactionsErr
	}

	zap.S().Debugf("micro_miner top block sig %s", a.state.TopBlock().BlockSignature)

	newTransactions := minedBlock.Transactions.Join(transactions)

	newBlock, err := proto.CreateBlock(
		newTransactions,
		minedBlock.Timestamp,
		minedBlock.Parent,
		minedBlock.GenPublicKey,
		minedBlock.NxtConsensus,
		minedBlock.Version,
		minedBlock.Features,
		minedBlock.RewardVote,
		a.scheme,
	)
	if err != nil {
		return nil, nil, rest, err
	}
	sk := keyPair.Secret

	err = newBlock.SetTransactionsRoot(a.scheme)
	if err != nil {
		return nil, nil, rest, err
	}
	err = newBlock.Sign(a.scheme, keyPair.Secret)
	if err != nil {
		return nil, nil, rest, err
	}
	err = newBlock.GenerateBlockID(a.scheme)
	if err != nil {
		return nil, nil, rest, err
	}

	//locked = mu.Lock()
	//_ = a.state.RollbackTo(minedBlock.Parent)
	//locked.Unlock()

	//err = a.services.BlocksApplier.Apply(a.state, []*proto.Block{newBlock})
	//if err != nil {
	//	zap.S().Error(err)
	//	return
	//}

	micro := proto.MicroBlock{
		VersionField:          5, // protobuf version
		SenderPK:              keyPair.Public,
		Transactions:          transactions,
		TransactionCount:      uint32(cnt),
		Reference:             a.state.TopBlock().BlockID(),
		TotalResBlockSigField: newBlock.BlockSignature,
	}

	err = micro.Sign(sk)
	if err != nil {
		return nil, nil, rest, err
	}

	zap.S().Debugf("micro_miner mined %+v", micro)

	inv := proto.NewUnsignedMicroblockInv(micro.SenderPK, newBlock.BlockID(), micro.Reference)
	err = inv.Sign(sk, a.scheme)
	if err != nil {
		return nil, nil, rest, err
	}

	newRest := proto.MiningLimits{
		MaxScriptRunsInBlock:        rest.MaxScriptRunsInBlock,
		MaxScriptsComplexityInBlock: rest.MaxScriptsComplexityInBlock,
		ClassicAmountOfTxsInBlock:   rest.ClassicAmountOfTxsInBlock,
		MaxTxsSizeInBytes:           rest.MaxTxsSizeInBytes - binSize,
	}

	//newBlocks, err := blocks.AddMicro(&micro)
	//if err != nil {
	//	return nil, nil, rest, err
	//}

	return newBlock, &micro, newRest, nil
}