package state

import (
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type recentBlocks struct {
	height, newestHeight uint64
	rangeSize            int
	// IDs of recent blocks in DB.
	stableIds []crypto.Signature
	isStable  map[crypto.Signature]uint64
	// IDs of recent blocks which have not been flushed to DB yet.
	newestIds []crypto.Signature
	isNewest  map[crypto.Signature]uint64
}

func newRecentBlocks(rangeSize int) (*recentBlocks, error) {
	return &recentBlocks{
		rangeSize: rangeSize,
		isStable:  make(map[crypto.Signature]uint64),
		isNewest:  make(map[crypto.Signature]uint64),
	}, nil
}

func (rb *recentBlocks) setStartHeight(startHeight uint64) {
	rb.height = startHeight
}

// Add to the list of newest IDs.
func (rb *recentBlocks) addNewBlockID(blockID crypto.Signature) error {
	if rb.newestHeight == 0 {
		rb.newestHeight = rb.height
	}
	rb.isNewest[blockID] = rb.newestHeight
	rb.newestIds = append(rb.newestIds, blockID)
	rb.newestHeight++
	return nil
}

// Add directly to the list of stable IDs.
func (rb *recentBlocks) addBlockID(blockID crypto.Signature) error {
	if len(rb.stableIds) < rb.rangeSize {
		rb.isStable[blockID] = rb.height
		rb.stableIds = append(rb.stableIds, blockID)
	} else {
		rb.isStable[blockID] = rb.height
		delete(rb.isStable, rb.stableIds[0])
		rb.stableIds = rb.stableIds[1:]
		rb.stableIds = append(rb.stableIds, blockID)
	}
	rb.height++
	return nil
}

func (rb *recentBlocks) newBlockIDToHeight(blockID crypto.Signature) (uint64, error) {
	height, ok := rb.isNewest[blockID]
	if !ok {
		height, ok = rb.isStable[blockID]
		if !ok {
			return 0, nil
		}
		return height, nil
	}
	return height, nil
}

func (rb *recentBlocks) blockIDToHeight(blockID crypto.Signature) (uint64, error) {
	height, ok := rb.isStable[blockID]
	if !ok {
		return 0, nil
	}
	return height, nil
}

func (rb *recentBlocks) isEmpty() bool {
	return rb.stableIds == nil
}

func (rb *recentBlocks) reset() {
	rb.stableIds = nil
	rb.isStable = make(map[crypto.Signature]uint64)
	rb.newestHeight = 0
	rb.height = 0
	rb.newestIds = nil
	rb.isNewest = make(map[crypto.Signature]uint64)
}

func (rb *recentBlocks) removeOutdated(ids []crypto.Signature) {
	for _, id := range ids {
		delete(rb.isStable, id)
	}
}

func (rb *recentBlocks) addNewIds() {
	for id, height := range rb.isNewest {
		rb.isStable[id] = height
	}
}

// flush() "flushes" newest IDs to stable IDs.
func (rb *recentBlocks) flush() {
	rb.stableIds = append(rb.stableIds, rb.newestIds...)
	rb.addNewIds()
	rb.newestIds = nil
	rb.isNewest = make(map[crypto.Signature]uint64)
	if len(rb.stableIds) > rb.rangeSize {
		rb.removeOutdated(rb.stableIds[:len(rb.stableIds)-rb.rangeSize])
		rb.stableIds = rb.stableIds[len(rb.stableIds)-rb.rangeSize:]
	}
	rb.height = rb.newestHeight
}
