package state_fsm

import (
	"time"

	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/node/peer_manager"
	"github.com/wavesplatform/gowaves/pkg/node/state_fsm/ng"
	. "github.com/wavesplatform/gowaves/pkg/node/state_fsm/tasks"
	"github.com/wavesplatform/gowaves/pkg/p2p/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/services"
	storage "github.com/wavesplatform/gowaves/pkg/state"
	"github.com/wavesplatform/gowaves/pkg/types"
)

type Async []Task

// TODO send score

type BlocksApplier interface {
	Apply(state storage.State, block []*proto.Block) error
}

type BaseInfo struct {
	// too old state
	//outdated bool

	// peers
	peers peer_manager.PeerManager

	// storage
	storage storage.State

	// ntp time
	tm types.Time

	// outdate period
	outdatePeriod proto.Timestamp

	scheme proto.Scheme

	//
	invRequester  InvRequester
	blockCreater  types.BlockCreator
	blocksApplier BlocksApplier

	// default behaviour
	d Default
}

type FSM interface {
	NewPeer(p peer.Peer) (FSM, Async, error)
	PeerError(peer.Peer, error) (FSM, Async, error)
	Score(p peer.Peer, score *proto.Score) (FSM, Async, error)
	Block(peer peer.Peer, block *proto.Block) (FSM, Async, error)

	// Received signatures after asking by GetSignatures
	Signatures(peer peer.Peer, sigs []crypto.Signature) (FSM, Async, error)
	//GetPeers(peer peer.Peer) (FSM, Async, error)
	Task(task AsyncTask) (FSM, Async, error)

	// micro
	MicroBlock(p peer.Peer, micro *proto.MicroBlock) (FSM, Async, error)
	MicroBlockInv(p peer.Peer, inv *proto.MicroBlockInv) (FSM, Async, error)
}

func NewFsm(
	services services.Services,
	outdatePeriod proto.Timestamp,
	blockCreator types.BlockCreator,

) (FSM, Async, error) {
	b := BaseInfo{
		peers:         services.Peers,
		storage:       services.State,
		tm:            services.Time,
		outdatePeriod: outdatePeriod,
		scheme:        services.Scheme,

		//
		invRequester:  ng.NewInvRequester(),
		blockCreater:  blockCreator,
		blocksApplier: services.BlocksApplier,

		//
		d: DefaultImpl{},
	}

	// default tasks
	tasks := Async{
		// ask about peers for every 5 minutes
		NewAskPeersTask(5 * time.Minute),
		NewPingTask(),
	}

	return NewIdleFsm(b), tasks, nil
}