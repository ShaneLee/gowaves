package node

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/wavesplatform/gowaves/pkg/ng"
	"github.com/wavesplatform/gowaves/pkg/node/peer_manager"
	"github.com/wavesplatform/gowaves/pkg/p2p/peer"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/services"
	"github.com/wavesplatform/gowaves/pkg/state"
	"github.com/wavesplatform/gowaves/pkg/types"
	"github.com/wavesplatform/gowaves/pkg/util/common"
	"go.uber.org/zap"
)

type Config struct {
	AppName  string
	NodeName string
	Listen   string
	DeclAddr string
}

type Node struct {
	peers     peer_manager.PeerManager
	state     state.State
	subscribe types.Subscribe
	sync      types.StateSync
	declAddr  proto.TCPAddr
	bindAddr  proto.TCPAddr
	scheduler types.Scheduler
	utx       types.UtxPool
	ng        *ng.RuntimeImpl
	services  services.Services
}

func NewNode(services services.Services, declAddr proto.TCPAddr, bindAddr proto.TCPAddr, ng *ng.RuntimeImpl, sync types.StateSync) *Node {
	if bindAddr.Empty() {
		bindAddr = declAddr
	}
	return &Node{
		state:     services.State,
		peers:     services.Peers,
		subscribe: services.Subscribe,
		sync:      sync,
		declAddr:  declAddr,
		bindAddr:  bindAddr,
		scheduler: services.Scheduler,
		utx:       services.UtxPool,
		ng:        ng,
		services:  services,
	}
}

func (a *Node) State() state.State {
	return a.state
}

func (a *Node) PeerManager() peer_manager.PeerManager {
	return a.peers
}

func (a *Node) HandleProtoMessage(mess peer.ProtoMessage) {
	switch t := mess.Message.(type) {
	case *proto.PeersMessage:
		a.handlePeersMessage(mess.ID, t)
	case *proto.GetPeersMessage:
		a.handleGetPeersMessage(mess.ID, t)
	case *proto.ScoreMessage:
		a.handleScoreMessage(mess.ID, t.Score)
	case *proto.BlockMessage:
		a.handleBlockMessage(mess.ID, t)
	case *proto.GetBlockMessage:
		a.handleBlockBySignatureMessage(mess.ID, t.BlockID)
	case *proto.SignaturesMessage:
		a.handleSignaturesMessage(mess.ID, t)
	case *proto.GetSignaturesMessage:
		a.handleGetSignaturesMessage(mess.ID, t)
	case *proto.TransactionMessage:
		a.handleTransactionMessage(mess.ID, t)
	case *proto.MicroBlockInvMessage:
		a.handleMicroblockInvMessage(mess.ID, t)
	case *proto.MicroBlockRequestMessage:
		a.handleMicroBlockRequestMessage(mess.ID, t)
	case *proto.MicroBlockMessage:
		a.handleMicroBlockMessage(mess.ID, t)
	case *proto.PBBlockMessage:
		a.handlePBBlockMessage(mess.ID, t)
	case *proto.PBMicroBlockMessage:
		a.handlePBMicroBlockMessage(mess.ID, t)
	case *proto.PBTransactionMessage:
		a.handlePBTransactionMessage(mess.ID, t)
	case *proto.GetBlockIdsMessage:
		a.handleGetBlockIdsMessage(mess.ID, t)
	case *proto.BlockIdsMessage:
		a.handleBlockIdsMessage(mess.ID, t)

	default:
		zap.S().Errorf("unknown proto Message %T", mess.Message)
	}
}

func (a *Node) handlePBBlockMessage(p peer.Peer, mess *proto.PBBlockMessage) {
	if !a.subscribe.Receive(p, mess) {
		b := &proto.Block{}
		err := b.UnmarshalFromProtobuf(mess.PBBlockBytes)
		if err != nil {
			zap.S().Debug(err)
			return
		}
		a.ng.HandleBlockMessage(p, b)
	}
}

func (a *Node) handlePBMicroBlockMessage(p peer.Peer, mess *proto.PBMicroBlockMessage) {
	a.ng.HandlePBMicroBlockMessage(p, mess)
}

func (a *Node) handlePBTransactionMessage(_ peer.Peer, mess *proto.PBTransactionMessage) {
	t, err := proto.SignedTxFromProtobuf(mess.Transaction)
	if err != nil {
		zap.S().Debug(err)
		return
	}
	_ = a.utx.AddWithBytes(t, common.Dup(mess.Transaction))
}

func (a *Node) handleTransactionMessage(_ peer.Peer, mess *proto.TransactionMessage) {
	t, err := proto.BytesToTransaction(mess.Transaction, a.services.Scheme)
	if err != nil {
		zap.S().Debug(err)
		return
	}
	_ = a.utx.AddWithBytes(t, common.Dup(mess.Transaction))
}

func (a *Node) handlePeersMessage(_ peer.Peer, peers *proto.PeersMessage) {
	var prs []proto.TCPAddr
	for _, p := range peers.Peers {
		prs = append(prs, proto.NewTCPAddr(p.Addr, int(p.Port)))
	}
	err := a.peers.UpdateKnownPeers(prs)
	if err != nil {
		zap.S().Error(err)
	}
}

func (a *Node) handleGetPeersMessage(p peer.Peer, m *proto.GetPeersMessage) {
	rs, err := a.peers.KnownPeers()
	if err != nil {
		zap.L().Error("failed got known peers", zap.Error(err))
		return
	}
	_, ok := a.peers.Connected(p)
	if !ok {
		// peer gone offline, skip
		return
	}

	var out []proto.PeerInfo
	for _, r := range rs {
		out = append(out, proto.PeerInfo{
			Addr: net.IP(r.IP[:]),
			Port: uint16(r.Port),
		})
	}

	p.SendMessage(&proto.PeersMessage{Peers: out})
}

func (a *Node) HandleInfoMessage(m peer.InfoMessage) {
	switch t := m.Value.(type) {
	case *peer.Connected:
		a.handleNewConnection(t.Peer)
	case error:
		a.handlePeerError(m.Peer, t)
	}
}

func (a *Node) AskPeers() {
	a.peers.AskPeers()
}

func (a *Node) handlePeerError(p peer.Peer, err error) {
	zap.S().Debug(err)
	a.peers.Suspend(p, err.Error())
}

func (a *Node) Close() {
	a.peers.Close()
	a.sync.Close()
	locked := a.state.Mutex().Lock()
	a.state.Close()
	locked.Unlock()
}

func (a *Node) handleNewConnection(p peer.Peer) {
	err := a.peers.NewConnection(p)
	if err != nil {
		return
	}

	// send score to new connected
	go func() {
		locked := a.state.Mutex().RLock()
		score, err := a.state.CurrentScore()
		locked.Unlock()
		if err != nil {
			zap.S().Error(err)
			return
		}
		p.SendMessage(&proto.ScoreMessage{
			Score: score.Bytes(),
		})
	}()
}

func (a *Node) handleBlockBySignatureMessage(p peer.Peer, id proto.BlockID) {
	block, err := a.state.Block(id)
	if err != nil {
		zap.S().Error(err)
		return
	}
	bm, err := proto.MessageByBlock(block, a.services.Scheme)
	if err != nil {
		zap.S().Error(err)
		return
	}
	p.SendMessage(bm)
}

func (a *Node) handleScoreMessage(p peer.Peer, score []byte) {
	b := new(big.Int)
	b.SetBytes(score)
	a.peers.UpdateScore(p, b)

	go func() {
		<-time.After(4 * time.Second)
		a.sync.Sync()
	}()

}

func (a *Node) handleBlockMessage(p peer.Peer, mess *proto.BlockMessage) {
	if !a.subscribe.Receive(p, mess) {
		b := &proto.Block{}
		err := b.UnmarshalBinary(mess.BlockBytes, a.services.Scheme)
		if err != nil {
			zap.S().Debug(err)
			return
		}
		a.ng.HandleBlockMessage(p, b)
	}
}

func (a *Node) handleGetSignaturesMessage(p peer.Peer, mess *proto.GetSignaturesMessage) {
	for _, sig := range mess.Blocks {
		id := proto.NewBlockIDFromSignature(sig)
		block, err := a.state.Block(id)
		if err != nil {
			continue
		}
		if block.BlockID() != id {
			panic("id error")
		}
		sendBlockIdsFromBlock(block, a.state, p)
		return
	}
}

func (a *Node) handleGetBlockIdsMessage(p peer.Peer, mess *proto.GetBlockIdsMessage) {
	for _, id := range mess.Blocks {
		block, err := a.state.Block(id)
		if err != nil {
			continue
		}
		if block.BlockID() != id {
			panic("id error")
		}
		sendBlockIdsFromBlock(block, a.state, p)
		return
	}
}

func (a *Node) handleMicroblockInvMessage(p peer.Peer, mess *proto.MicroBlockInvMessage) {
	a.ng.HandleInvMessage(p, mess)
}

func (a *Node) handleMicroBlockRequestMessage(p peer.Peer, mess *proto.MicroBlockRequestMessage) {
	a.ng.HandleMicroBlockRequestMessage(p, mess)
}

func (a *Node) SpawnOutgoingConnections(ctx context.Context) {
	a.peers.SpawnOutgoingConnections(ctx)
}

func (a *Node) SpawnOutgoingConnection(ctx context.Context, addr proto.TCPAddr) error {
	return a.peers.Connect(ctx, addr)
}

func (a *Node) Serve(ctx context.Context) error {
	// if empty declared address, listen on port doesn't make sense
	if a.declAddr.Empty() {
		return nil
	}

	if a.bindAddr.Empty() {
		return nil
	}

	zap.S().Info("start listening on ", a.bindAddr.String())
	l, err := net.Listen("tcp", a.bindAddr.String())
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			zap.S().Error(err)
			continue
		}

		go func() {
			if err := a.peers.SpawnIncomingConnection(ctx, conn); err != nil {
				zap.S().Error(err)
				return
			}
		}()
	}
}

func (a *Node) handleMicroBlockMessage(p peer.Peer, message *proto.MicroBlockMessage) {
	a.ng.HandleMicroBlockMessage(p, message)
}

func (a *Node) handleSignaturesMessage(p peer.Peer, message *proto.SignaturesMessage) {
	a.subscribe.Receive(p, message)
}

func (a *Node) handleBlockIdsMessage(p peer.Peer, message *proto.BlockIdsMessage) {
	a.subscribe.Receive(p, message)
}

func (n *Node) Run(ctx context.Context, p peer.Parent) {
	go n.sync.Run(ctx)

	go func() {
		for {
			n.SpawnOutgoingConnections(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Minute):
			}
		}
	}()

	go func() {
		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
			return
		}

		n.AskPeers()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(4 * time.Minute):
				n.AskPeers()
			}
		}
	}()

	// info messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-p.InfoCh:
				n.HandleInfoMessage(m)
			}
		}
	}()

	go func() {
		if err := n.Serve(ctx); err != nil {
			return
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-p.MessageCh:
			n.services.LoggableRunner.Named(fmt.Sprintf("Node.Run.Handler.%T", m.Message), func() {
				n.HandleProtoMessage(m)
			})
		}
	}
}

type BlockIds struct {
	ids    []proto.BlockID
	unique map[proto.BlockID]struct{}
}

func (a *BlockIds) Ids() []proto.BlockID {
	return a.ids
}

func NewBlockIds(ids ...proto.BlockID) *BlockIds {
	unique := make(map[proto.BlockID]struct{})
	for _, v := range ids {
		unique[v] = struct{}{}
	}

	return &BlockIds{
		ids:    ids,
		unique: unique,
	}
}

func (a *BlockIds) Exists(id proto.BlockID) bool {
	_, ok := a.unique[id]
	return ok
}

func (a *BlockIds) Revert() *BlockIds {
	out := make([]proto.BlockID, len(a.ids))
	for k, v := range a.ids {
		out[len(a.ids)-1-k] = v
	}
	return NewBlockIds(out...)
}
