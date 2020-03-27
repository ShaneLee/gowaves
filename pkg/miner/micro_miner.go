package miner

import (
	"errors"

	"github.com/wavesplatform/gowaves/pkg/node/state_fsm/ng"
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
	blocks ng.Blocks,
	keyPair proto.KeyPair) (*proto.Block, *proto.MicroBlock, proto.MiningLimits, ng.Blocks, error) {

	// way to stop mine microblocks
	if minedBlock == nil {
		return nil, nil, rest, blocks, errors.New("no block provided")
	}

	topBlock := a.state.TopBlock()
	if topBlock.BlockSignature != minedBlock.BlockSignature {
		// block changed, exit
		return nil, nil, rest, blocks, StateChangedErr
	}

	//height, err := a.state.Height()
	//if err != nil {
	//	zap.S().Error(err)
	//	return
	//}
	//
	//topBlock, err := a.state.BlockByHeight(height)
	//if err != nil {
	//	zap.S().Error(err)
	//	return
	//}

	rlocked := a.state.Mutex().RLock()
	height, err := a.state.Height()
	rlocked.Unlock()
	if err != nil {
		return nil, nil, rest, blocks, err
	}

	parentTimestamp := topBlock.Timestamp
	if height > 1 {
		parent, err := a.state.BlockByHeight(height - 1)
		if err != nil {
			return nil, nil, rest, blocks, err
		}
		parentTimestamp = parent.Timestamp
	}

	//
	transactions := make([]proto.Transaction, 0)
	cnt := 0
	binSize := 0

	var unAppliedTransactions []*types.TransactionWithBytes

	mu := a.state.Mutex()
	locked := mu.Lock()

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
	locked.Unlock()

	// return unapplied transactions
	for _, unapplied := range unAppliedTransactions {
		_ = a.utx.AddWithBytes(unapplied.T, unapplied.B)
	}

	// no transactions applied, skip
	if cnt == 0 {
		return nil, nil, rest, blocks, NoTransactionsErr
	}
	row := blocks.Row()
	lastsig := row.LastSignature()

	zap.S().Debugf("micro_miner last sig %s", lastsig)
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
		return nil, nil, rest, blocks, err
	}

	sk := keyPair.Secret
	err = newBlock.Sign(a.scheme, keyPair.Secret)
	if err != nil {
		return nil, nil, rest, blocks, err
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
		VersionField:          3,
		SenderPK:              keyPair.Public,
		Transactions:          transactions,
		TransactionCount:      uint32(cnt),
		PrevResBlockSigField:  lastsig,
		TotalResBlockSigField: newBlock.BlockSignature,
	}

	err = micro.Sign(sk)
	if err != nil {
		return nil, nil, rest, blocks, err
	}

	zap.S().Debugf("micro_miner mined %+v", micro)

	inv := proto.NewUnsignedMicroblockInv(micro.SenderPK, micro.TotalResBlockSigField, micro.PrevResBlockSigField)
	err = inv.Sign(sk, a.scheme)
	if err != nil {
		return nil, nil, rest, blocks, err
	}

	newRest := proto.MiningLimits{
		MaxScriptRunsInBlock:        rest.MaxScriptRunsInBlock,
		MaxScriptsComplexityInBlock: rest.MaxScriptsComplexityInBlock,
		ClassicAmountOfTxsInBlock:   rest.ClassicAmountOfTxsInBlock,
		MaxTxsSizeInBytes:           rest.MaxTxsSizeInBytes - binSize,
	}

	newBlocks, err := blocks.AddMicro(&micro)
	if err != nil {
		return nil, nil, rest, blocks, err
	}

	return newBlock, &micro, newRest, newBlocks, nil

	//go a.mineMicro(ctx, newRest, newBlock, newBlocks, keyPair)
}
