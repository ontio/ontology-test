package ontid

import (
	"time"

	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

type Contract struct {
	Address common.Address
	Method  string
	Args    []interface{}
}

func getEvent(ctx *testframework.TestFrameworkContext, txHash common.Uint256) bool {
	_, err := ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error: %s", err)
		return false
	}
	events, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("GetSmartContractEvent error: %s", err)
		return false
	}

	if events.State == 0 {
		ctx.LogWarn("ontio contract invoke failed, state:0")
		return false
	}

	if len(events.Notify) > 0 {
		states := events.Notify[0].States
		ctx.LogInfo("result is : %+v", states)
		return true
	} else {
		return false
	}

}

func InvokeContract(ctx *testframework.TestFrameworkContext, contract *Contract, isPreInvoke bool) (bool, []byte) {
	if !isPreInvoke {
		user, _ := ctx.GetDefaultAccount()
		txHash, err := ctx.Ont.Rpc.InvokeNativeContract(
			ctx.GetGasPrice(),
			ctx.GetGasLimit(),
			user,
			0,
			contract.Address,
			contract.Method,
			contract.Args,
		)
		if err != nil {
			ctx.LogError("InvokeNativeContract error: %s", err)
			return false, nil
		}
		return getEvent(ctx, txHash), nil
	} else {
		res, err := ctx.Ont.Rpc.PrepareInvokeNativeContract(
			contract.Address,
			0,
			contract.Method,
			contract.Args,
		)
		if err != nil {
			ctx.LogError("PrepareInvokeNativeContract failed, %s", err)
			return false, nil
		}

		buf, ok := res.Result.([]byte)
		if !ok {
			ctx.LogError("error result value type")
			return false, nil
		}
		return true, buf
	}
}
