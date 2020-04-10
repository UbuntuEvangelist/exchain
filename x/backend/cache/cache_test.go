package cache

import (
	"testing"

	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	require.Equal(t, 0, len(cache.Transactions))
	require.Equal(t, 2000, cap(cache.Transactions))
	require.Equal(t, 0, len(cache.LatestTicker))
	require.Equal(t, 0, len(cache.ProductsBuf))
	require.Equal(t, 200, cap(cache.ProductsBuf))

	txs := []*types.Transaction{
		{TxHash: "hash1", Type: types.TxTypeTransfer, Address: "addr1", Symbol: common.TestToken, Side: types.TxSideFrom, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 100},
		{TxHash: "hash2", Type: types.TxTypeOrderNew, Address: "addr1", Symbol: types.TestTokenPair, Side: types.TxSideBuy, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 300},
		{TxHash: "hash3", Type: types.TxTypeOrderCancel, Address: "addr1", Symbol: types.TestTokenPair, Side: types.TxSideSell, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 200},
		{TxHash: "hash4", Type: types.TxTypeTransfer, Address: "addr2", Symbol: common.TestToken, Side: types.TxSideTo, Quantity: "10.0", Fee: "0.1" + common.NativeToken, Timestamp: 100},
	}

	for _, tx := range txs {
		cache.AddTransaction(tx)
	}

	require.Equal(t, txs, cache.GetTransactions())

	cache.Flush()
	require.Equal(t, 0, len(cache.GetTransactions()))

}
