package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	g "github.com/wavesplatform/gowaves/pkg/grpc/generated"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"google.golang.org/grpc"
)

func TestTransactionsAPIClient(t *testing.T) {
	t.SkipNow()
	conn := connect(t)
	defer conn.Close()

	c := g.NewTransactionsApiClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := g.TransactionsRequest{}
	var err error
	uc, err := c.GetUnconfirmed(ctx, &req, grpc.EmptyCallOption{})
	require.NoError(t, err)
	var msg g.TransactionResponse
	for err = uc.RecvMsg(&msg); err == nil; err = uc.RecvMsg(&msg) {
		c := SafeConverter{}
		tx, err := c.SignedTransaction(msg.Transaction)
		require.NoError(t, err)
		js, err := json.Marshal(tx)
		require.NoError(t, err)
		fmt.Println(string(js))
	}
	assert.Equal(t, io.EOF, err)
}

func TestBlocksAPIClient(t *testing.T) {
	t.SkipNow()
	conn := connect(t)
	defer conn.Close()

	c := g.NewBlocksApiClient(conn)

	getBlock := func(h int) (*g.BlockWithHeight, error) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		req := &g.BlockRequest{IncludeTransactions: true, Request: &g.BlockRequest_Height{Height: int32(h)}}
		return c.GetBlock(ctx, req, grpc.EmptyCallOption{})
	}

	var err error
	var b *g.BlockWithHeight
	cnv := SafeConverter{}
	h := 1
	for b, err = getBlock(h); err == nil; b, err = getBlock(h) {
		cnv.Reset()
		txs, err := cnv.BlockTransactions(b)
		require.NoError(t, err)
		sb := strings.Builder{}
		sb.WriteRune('[')
		sb.WriteString(strconv.Itoa(len(txs)))
		sb.WriteRune(']')
		sb.WriteRune(' ')
		for _, tx := range txs {
			js, err := json.Marshal(tx)
			require.NoError(t, err)
			sb.WriteString(string(js))
			sb.WriteRune(',')
		}
		header, err := cnv.BlockHeader(b)
		require.NoError(t, err)
		bjs, err := json.Marshal(header)
		require.NoError(t, err)
		fmt.Println("HEIGHT:", b.Height, "BLOCK:", string(bjs), "TXS:", sb.String())
		h++
	}
	assert.Equal(t, io.EOF, err)
}

func TestAccountData(t *testing.T) {
	//t.SkipNow()
	conn := connect(t)
	defer conn.Close()

	c := g.NewAccountsApiClient(conn)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))
	defer cancel()
	addr, err := proto.NewAddressFromString("3NCzApG3ka4tvDuKX8HxnJtXe7eJw5PmdVt")
	require.NoError(t, err)
	b, err := addr.Body()
	require.NoError(t, err)
	asset, err := crypto.NewDigestFromBase58("Gf9t8FA4H3ssoZPCwrg3KwUFCci8zuUFP9ssRsUY3s6a")
	require.NoError(t, err)
	a := make([]byte, len(asset.Bytes()))
	copy(a, asset.Bytes())
	req := &g.BalancesRequest{
		Address: b,
		Assets:  [][]byte{a},
	}
	bc, err := c.GetBalances(ctx, req)
	require.NoError(t, err)
	var msg g.BalanceResponse
	for err = bc.RecvMsg(&msg); err == nil; err = bc.RecvMsg(&msg) {
		br, ok := msg.Balance.(*g.BalanceResponse_Asset)
		require.True(t, ok)
		fmt.Printf("%s: %d\n", base58.Encode(br.Asset.AssetId), int(br.Asset.Amount))
	}
	assert.Equal(t, io.EOF, err)
}

func connect(t *testing.T) *grpc.ClientConn {
	conn, err := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	require.NoError(t, err)
	return conn
}
