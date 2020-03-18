package state_fsm

import (
	"net"

	"github.com/wavesplatform/gowaves/pkg/node/peer_manager"
	. "github.com/wavesplatform/gowaves/pkg/p2p/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/types"
	"go.uber.org/zap"
)

func sendPeers(fsm FSM, p Peer, peers peer_manager.PeerManager) (FSM, Async, error) {
	rs, err := peers.KnownPeers()
	if err != nil {
		zap.L().Error("failed got known peers", zap.Error(err))
		return fsm, nil, err
	}

	var out []proto.PeerInfo
	for _, r := range rs {
		out = append(out, proto.PeerInfo{
			Addr: net.IP(r.IP[:]),
			Port: uint16(r.Port),
		})
	}

	p.SendMessage(&proto.PeersMessage{Peers: out})
	return fsm, nil, nil
}

func newPeer(fsm FSM, p Peer, peers peer_manager.PeerManager) (FSM, Async, error) {
	err := peers.NewConnection(p)
	return fsm, nil, err
}

// TODO handle no peers
func peerError(fsm FSM, p Peer, peers peer_manager.PeerManager, _ error) (FSM, Async, error) {
	peers.Disconnect(p)
	return fsm, nil, nil
}

func noop(fsm FSM) (FSM, Async, error) {
	return fsm, nil, nil
}

func IsOutdate(period proto.Timestamp, lastBlock *proto.Block, tm types.Time) bool {
	curTime := proto.NewTimestampFromTime(tm.Now())
	return curTime-lastBlock.Timestamp > period
}

func handleScore(fsm FSM, info BaseInfo, p Peer, score *proto.Score) (FSM, Async, error) {
	err := info.peers.UpdateScore(p, score)
	if err != nil {
		return fsm, nil, err
	}

	locked := info.storage.Mutex().Lock()
	myScore, err := info.storage.CurrentScore()
	locked.Unlock()
	if err != nil {
		return NewIdleFsm(info), nil, err
	}

	if score.Cmp(myScore) == 1 { // remote score > my score
		return NewIdleToSyncTransition(info, p)
	}
	return fsm, nil, nil
}