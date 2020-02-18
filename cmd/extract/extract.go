package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func main() {
	var (
		blockchainFile = flag.String("blockchain-file", "", "Path to binary blockchain file to collect RIDE scripts from.")
	)
	flag.Parse()

	if *blockchainFile == "" {
		flag.Usage()
		log.Fatalln("No blockchain file")
	}
	in, err := os.Open(*blockchainFile)
	if err != nil {
		log.Fatalf("Failed to open blockchain file '%s': %v", *blockchainFile, err)
	}
	defer func() {
		err = in.Close()
		if err != nil {
			log.Printf("Failed to close blockchain file: %v", err)
		}
	}()
	count, assets, err := countInvokesWithAssetFee(in)
	if err != nil {
		log.Fatalf("Failed to count transactions in given file: %v", err)
	}
	log.Printf("Found %d assets in %d transactions", len(assets), count)
	for k, v := range assets {
		log.Printf("\t%s: %d", k.String(), v)
	}
}

func countInvokesWithAssetFee(f *os.File) (int, map[crypto.Digest]int, error) {
	count := 0
	assets := make(map[crypto.Digest]int)

	h := 1
	sizeBuf := make([]byte, 4)
	blockBuf := make([]byte, 2*1024*1024)
	for {
		h++
		_, err := io.ReadFull(f, sizeBuf)
		if err != nil {
			if err != io.EOF {
				return 0, nil, fmt.Errorf("unable to read block size: %v", err)
			}
			return count, assets, nil
		}
		size := binary.BigEndian.Uint32(sizeBuf)
		bb := blockBuf[:size]
		_, err = io.ReadFull(f, bb)
		if err != nil {
			if err != io.EOF {
				return 0, nil, fmt.Errorf("unable to read block: %v", err)
			}
			return 0, nil, errors.New("unexpected EOF")
		}
		var header proto.BlockHeader
		err = header.UnmarshalHeaderFromBinary(bb)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to parse block's header at height %d: %v", h, err)
		}
		if header.Version >= 3 {
			var block proto.Block
			err = block.UnmarshalBinary(bb)
			if err != nil {
				return 0, nil, fmt.Errorf("failed to parse block at height %d: %v", h, err)
			}
			txs, err := block.Transactions.Transactions()
			if err != nil {
				return 0, nil, fmt.Errorf("failed to extract transaction from block at height %d: %v", h, err)
			}
			for _, tx := range txs {
				switch atx := tx.(type) {
				case *proto.InvokeScriptV1:
					if atx.FeeAsset.Present {
						count++
						assets[atx.FeeAsset.ID]++
					}
				}
			}
		}
	}
}
