package state_fsm

import (
	"github.com/wavesplatform/gowaves/pkg/p2p/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/services"
	"go.uber.org/zap"
)

type Actions interface {
	SendScore(curScore *proto.Score)
	SendBlock(block *proto.Block)
}

type ActionsImpl struct {
	services services.Services
}

func (a *ActionsImpl) SendScore(curScore *proto.Score) {
	bts := curScore.Bytes()
	a.services.Peers.EachConnected(func(peer peer.Peer, score *proto.Score) {
		peer.SendMessage(&proto.ScoreMessage{Score: bts})
	})
}

func (a *ActionsImpl) SendBlock(block *proto.Block) {
	zap.S().Info("SendBlock called with")
	bts, err := block.MarshalToProtobuf(a.services.Scheme)
	if err != nil {
		zap.S().Error(err)
		return
	}
	a.services.Peers.EachConnected(func(p peer.Peer, score *proto.Score) {
		p.SendMessage(&proto.PBBlockMessage{PBBlockBytes: bts})
	})
}
