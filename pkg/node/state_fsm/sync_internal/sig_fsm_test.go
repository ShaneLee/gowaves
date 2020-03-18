package sync_internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/libs/ordered_blocks"
	"github.com/wavesplatform/gowaves/pkg/libs/signatures"
	"github.com/wavesplatform/gowaves/pkg/mock"
	. "github.com/wavesplatform/gowaves/pkg/node/state_fsm/sync_internal"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

var sig1 = crypto.MustSignatureFromBase58("5syuWANDSgk8KyPxq2yQs2CYV23QfnrBoZMSv2LaciycxDYfBw6cLA2SqVnonnh1nFiFumzTgy2cPETnE7ZaZg5P")
var sig2 = crypto.MustSignatureFromBase58("3kvbjSovZWLg1zdMyW5vGsCj1DR1jkHY3ALtu5VxoqscrXQq3nH2vS2V5dhVo6ff9bxtbFAkUkVQQqCFUAHmwnpX")

func TestSigFSM_Signatures(t *testing.T) {
	or := ordered_blocks.NewOrderedBlocks()
	sigs := signatures.NewSignatures()

	t.Run("error on receive unexpected signatures", func(t *testing.T) {
		fsm := NewSigFSM(or, sigs, NoSignaturesExpected, false)
		rs2, err := fsm.Signatures(nil, []crypto.Signature{sig1, sig2})
		require.Equal(t, NoSignaturesExpectedErr, err)
		require.NotNil(t, rs2)
	})

	t.Run("successful receive signatures", func(t *testing.T) {
		fsm := NewSigFSM(or, sigs, WaitingForSignatures, false)
		rs2, err := fsm.Signatures(mock.NoOpPeer{}, []crypto.Signature{sig1, sig2})
		require.NoError(t, err)
		require.NotNil(t, rs2)
		require.True(t, rs2.NearEnd())
		require.False(t, rs2.WaitingForSignatures())
	})
}

func block(sig crypto.Signature) *proto.Block {
	return &proto.Block{
		BlockHeader: proto.BlockHeader{
			BlockSignature: sig,
		},
	}
}

func TestSigFSM_Block(t *testing.T) {
	or := ordered_blocks.NewOrderedBlocks()
	sigs := signatures.NewSignatures()
	fsm := NewSigFSM(or, sigs, WaitingForSignatures, false)
	fsm, _ = fsm.Signatures(mock.NoOpPeer{}, []crypto.Signature{sig1, sig2})

	fsm, _ = fsm.Block(block(sig1))
	fsm, _ = fsm.Block(block(sig2))
	require.Equal(t, 2, fsm.AvailableCount())

	// no panic, cause `nearEnd` is True
	fsm, blocks := fsm.Blocks(nil)
	require.Equal(t, 2, len(blocks))

}

func TestSigFSM_BlockGetSignatures(t *testing.T) {
	or := ordered_blocks.NewOrderedBlocks()
	sigs := signatures.NewSignatures()
	require.Panics(t, func() {
		NewSigFSM(or, sigs, NoSignaturesExpected, false).Blocks(nil)
	})
}
