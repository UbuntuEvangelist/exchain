package keeper

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/okex/exchain/x/evm/types"
)

func (k *Keeper) FixLog(execResults [][]string) [][]byte {
	res := make([][]byte, len(execResults), len(execResults))
	logSize := uint(0)
	txInBlock := int(-1)
	k.Bloom = new(big.Int)
	ts := time.Now()
	ttB := time.Duration(0)
	ttE := time.Duration(0)

	for index := 0; index < len(execResults); index++ {
		rs, ok := k.LogsManages.Get(execResults[index][0])
		if !ok || execResults[index][1] != "" {
			continue
		}
		txInBlock++
		if rs.ResultData == nil {
			continue
		}

		for _, v := range rs.ResultData.Logs {
			v.Index = logSize
			v.TxIndex = uint(txInBlock)
			logSize++
		}

		ttBloom := time.Now()
		k.Bloom = k.Bloom.Or(k.Bloom, rs.ResultData.Bloom.Big())
		ttB += time.Now().Sub(ttBloom)

		ttEncode := time.Now()
		data, err := types.EncodeResultData(*rs.ResultData)
		ttE += time.Now().Sub(ttEncode)
		if err != nil {
			panic(err)
		}

		res[index] = data
	}
	k.LogsManages.Reset()
	fmt.Println("fixxxx (second)", ttB.Seconds(), ttE.Seconds(), time.Now().Sub(ts).Seconds())
	return res
}

type LogsManager struct {
	mu      sync.RWMutex
	Results map[string]TxResult
}

func NewLogManager() *LogsManager {
	return &LogsManager{
		mu:      sync.RWMutex{},
		Results: make(map[string]TxResult),
	}
}

func (l *LogsManager) Set(txBytes string, value TxResult) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Results[txBytes] = value
}

func (l *LogsManager) Get(txBytes string) (TxResult, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	data, ok := l.Results[txBytes]
	return data, ok
}

func (l *LogsManager) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.Results)
}

func (l *LogsManager) Reset() {
	l.Results = make(map[string]TxResult)
}

type TxResult struct {
	ResultData *types.ResultData
	Err        error
}
