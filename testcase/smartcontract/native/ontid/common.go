package ontid

import (
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-test/testframework"
	sdk"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
)

type Contract struct {
	Address common.Address
	Method  string
	Args    []interface{}
}

func getEvent(ctx *testframework.TestFrameworkContext, txHash common.Uint256) bool {
	_, err := ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error: %s", err)
		return false
	}
	events, err := ctx.Ont.GetSmartContractEvent(txHash.ToHexString())
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
		txHash, err := ctx.Ont.Native.InvokeNativeContract(
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
		res, err := ctx.Ont.Native.PreExecInvokeNativeContract(
			contract.Address,
			0,
			contract.Method,
			contract.Args,
		)
		if err != nil {
			ctx.LogError("PrepareInvokeNativeContract failed, %s", err)
			return false, nil
		}

		data, err := res.Result.ToByteArray()
		if err != nil {
			ctx.LogError("error result value type")
			return false, nil
		}
		ctx.LogInfo("result: %s", data)
		return true, data
	}
}

func MultiSigInvoke(ctx *testframework.TestFrameworkContext, c *Contract, m uint16, pubs []keypair.PublicKey, user *sdk.Account) bool {
	tx, err := ctx.Ont.Native.NewNativeInvokeTransaction(
		ctx.GetGasPrice(),
		ctx.GetGasLimit(),
		0,
		c.Address,
		c.Method,
		c.Args,
	)
	if err != nil {
		ctx.LogError(err)
		return false
	}

	err = ctx.Ont.MultiSignToTransaction(tx, m, pubs, user)
	if err != nil {
		ctx.LogError("MultiSignToTransaction error: %s", err)
		return false
	}
	txHash, err := ctx.Ont.SendTransaction(tx)
	if err != nil {
		ctx.LogError("SendRawTransaction error: %s", err)
	}
	return getEvent(ctx, txHash)
}
