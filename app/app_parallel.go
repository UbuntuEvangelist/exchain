package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// feeCollectorHandler set or get the value of feeCollectorAcc
func updateFeeCollectorHandler(bk bank.Keeper, sk supply.Keeper) sdk.UpdateFeeCollectorAccHandler {
	return func(ctx sdk.Context, balance sdk.Coins) error {
		return bk.SetCoins(ctx, sk.GetModuleAddress(auth.FeeCollectorName), balance)
	}
}

// evmTxFeeHandler get tx fee for evm tx
func evmTxFeeHandler(ak auth.AccountKeeper) sdk.GetTxFeeHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (fee sdk.Coins, isEvm bool, signerCache sdk.SigCache, from ethcmn.Address, to *ethcmn.Address) {
		if evmTx, ok := tx.(evmtypes.MsgEthereumTx); ok {
			isEvm = true
			signerCache, _ = evmTx.VerifySig(evmTx.ChainID(), ctx.BlockHeight(), ctx.SigCache())
			from = ethcmn.BytesToAddress(evmTx.From())
			to = evmTx.To()
			ak.GetAccount(ctx, evmTx.From())
		}
		if feeTx, ok := tx.(authante.FeeTx); ok {
			fee = feeTx.GetFee()
		}
		return
	}
}

// fixLogForParallelTxHandler fix log for parallel tx
func fixLogForParallelTxHandler(ek *evm.Keeper) sdk.LogFix {
	return func(execResults [][]string) (logs [][]byte) {
		return ek.FixLog(execResults)
	}
}
