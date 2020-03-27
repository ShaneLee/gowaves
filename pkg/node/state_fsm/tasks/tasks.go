package tasks

import (
	"context"
	"time"

	"github.com/wavesplatform/gowaves/pkg/proto"
)

const (
	PING = iota + 1
	ASK_PEERS

	MINE_MICRO
)

type TaskType int

//type Async []Task

type AsyncTask struct {
	TaskType int
	Data     interface{}
}

type Task interface {
	Run(ctx context.Context, output chan AsyncTask) error
	Type() int
}

func Tasks(tasks ...Task) []Task {
	return tasks
}

type AskPeersTask struct {
	type_ int
	d     time.Duration
}

func NewAskPeersTask(d time.Duration) AskPeersTask {
	return AskPeersTask{
		type_: ASK_PEERS,
		d:     d,
	}
}

func (a AskPeersTask) Type() int {
	return a.type_
}

func (a AskPeersTask) Run(ctx context.Context, output chan AsyncTask) error {
	<-time.After(5 * time.Second)
	output <- AsyncTask{
		TaskType: a.type_,
	}

	in := time.NewTicker(a.d)
	defer in.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-in.C:
			output <- AsyncTask{
				TaskType: a.type_,
			}
		}
	}
}

type PingTask struct {
}

func NewPingTask() Task {
	return PingTask{}
}

func (PingTask) Type() int {
	return PING
}

func (PingTask) Run(ctx context.Context, output chan AsyncTask) error {
	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			output <- AsyncTask{
				TaskType: PING,
				Data:     nil,
			}
		}
	}
}

type MineMicroTaskData struct {
	Block   *proto.Block
	Limits  proto.MiningLimits
	KeyPair proto.KeyPair
}

type MineMicroTask struct {
	timeout           time.Duration
	MineMicroTaskData MineMicroTaskData
}

func NewMineMicroTask(timeout time.Duration, block *proto.Block, limits proto.MiningLimits, keyPair proto.KeyPair) MineMicroTask {
	if block == nil {
		panic("NewMineMicroTask block eq nil")
	}
	return MineMicroTask{
		timeout: timeout,
		MineMicroTaskData: MineMicroTaskData{
			Block:   block,
			Limits:  limits,
			KeyPair: keyPair,
		},
	}
}

func (MineMicroTask) Type() int {
	return MINE_MICRO
}

func (a MineMicroTask) Run(ctx context.Context, output chan AsyncTask) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(a.timeout):
		output <- AsyncTask{
			TaskType: a.Type(),
			Data:     a.MineMicroTaskData,
		}
	}
	return nil
}

//type GetSignaturesTimoutTask struct {
//	duration time.Duration
//}
//
//func NewGetSignaturesTimoutTask(duration time.Duration) GetSignaturesTimoutTask {
//	return GetSignaturesTimoutTask{duration: duration}
//}
//
//func (a GetSignaturesTimoutTask) Run(ctx context.Context, output chan AsyncTask) error {
//	select {
//	case <-ctx.Done():
//		return ctx.Err()
//	case <-time.After(a.duration):
//		output <- AsyncTask{
//			TaskType: SYNC_WAIT_SIGNATURES_TIMEOUT,
//			Data:     nil,
//		}
//	}
//	return nil
//}
