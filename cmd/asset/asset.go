package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func main() {
	var (
		blockchainFile = flag.String("blockchain-file", "", "Path to binary blockchain file to collect RIDE scripts from.")
		asset          = flag.String("asset", "", "Base58 encoded asset's ID.")
	)
	flag.Parse()

	if *blockchainFile == "" || *asset == "" {
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
	assetID, err := crypto.NewDigestFromBase58(*asset)
	if err != nil {
		log.Fatalf("Invalid asset: %v", err)
	}
	quantity, cropped, err := assetQuantity(in, assetID)
	if err != nil {
		log.Fatalf("Failed to count transactions in given file: %v", err)
	}
	log.Printf("Calculated asset quantity: %s", quantity.String())
	log.Printf("Calculated asset cropped: %d (%d)", cropped, int(quantity.Int64()))
}

func assetQuantity(f *os.File, asset crypto.Digest) (*big.Int, int, error) {
	quantity := big.NewInt(0)
	cropped := 0
	h := 1
	sizeBuf := make([]byte, 4)
	blockBuf := make([]byte, 2*1024*1024)
	for {
		h++
		_, err := io.ReadFull(f, sizeBuf)
		if err != nil {
			if err != io.EOF {
				return nil, 0, fmt.Errorf("unable to read block size: %v", err)
			}
			return quantity, cropped, nil
		}
		size := binary.BigEndian.Uint32(sizeBuf)
		bb := blockBuf[:size]
		_, err = io.ReadFull(f, bb)
		if err != nil {
			if err != io.EOF {
				return nil, 0, fmt.Errorf("unable to read block: %v", err)
			}
			return nil, 0, errors.New("unexpected EOF")
		}
		var block proto.Block
		err = block.UnmarshalBinary(bb)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to parse block at height %d: %v", h, err)
		}
		txs, err := block.Transactions.Transactions()
		if err != nil {
			return nil, 0, fmt.Errorf("failed to extract transaction from block at height %d: %v", h, err)
		}
		for _, tx := range txs {
			switch atx := tx.(type) {
			case *proto.IssueV1:
				if bytes.Equal(asset.Bytes(), atx.ID.Bytes()) {
					quantity.Add(quantity, big.NewInt(int64(atx.Quantity)))
					cropped += int(atx.Quantity)
					log.Printf("ISSUEv1\t%s\t%d\t%s\t(%d)", atx.ID.String(), int(atx.Quantity), quantity.String(), cropped)
				}
			case *proto.IssueV2:
				if bytes.Equal(asset.Bytes(), atx.ID.Bytes()) {
					quantity.Add(quantity, big.NewInt(int64(atx.Quantity)))
					cropped += int(atx.Quantity)
					log.Printf("ISSUEv2\t%s\t%d\t%s\t(%d)", atx.ID.String(), int(atx.Quantity), quantity.String(), cropped)
				}
			case *proto.ReissueV1:
				if bytes.Equal(asset.Bytes(), atx.AssetID.Bytes()) {
					quantity.Add(quantity, big.NewInt(int64(atx.Quantity)))
					cropped += int(atx.Quantity)
					log.Printf("REISSUEv1\t%s\t%d\t%s\t(%d)", atx.ID.String(), int(atx.Quantity), quantity.String(), cropped)
				}
			case *proto.ReissueV2:
				if bytes.Equal(asset.Bytes(), atx.AssetID.Bytes()) {
					quantity.Add(quantity, big.NewInt(int64(atx.Quantity)))
					cropped += int(atx.Quantity)
					log.Printf("REISSUEv2\t%s\t%d\t%s\t(%d)", atx.ID.String(), int(atx.Quantity), quantity.String(), cropped)
				}
			case *proto.BurnV1:
				if bytes.Equal(asset.Bytes(), atx.AssetID.Bytes()) {
					quantity.Sub(quantity, big.NewInt(int64(atx.Amount)))
					cropped -= int(atx.Amount)
					log.Printf("BURNv1\t%s\t%d\t%s\t(%d)", atx.ID.String(), int(atx.Amount), quantity.String(), cropped)
				}
			case *proto.BurnV2:
				if bytes.Equal(asset.Bytes(), atx.AssetID.Bytes()) {
					quantity.Sub(quantity, big.NewInt(int64(atx.Amount)))
					cropped -= int(atx.Amount)
					log.Printf("BURNv2\t%s\t%d\t%s\t(%d)", atx.ID.String(), int(atx.Amount), quantity.String(), cropped)
				}
			}
		}
	}
}
