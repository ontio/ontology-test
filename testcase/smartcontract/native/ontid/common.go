package ontid

import (
	"encoding/hex"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
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

		str, ok := res.Result.(string)
		if !ok {
			ctx.LogError("error result value type")
			return false, nil
		}
		ctx.LogInfo("result: %s", str)
		buf, err := hex.DecodeString(str)
		if err != nil {
			ctx.LogError("error hex code")
			return false, nil
		}
		return true, buf
	}
}

func MultiSigInvoke(ctx *testframework.TestFrameworkContext, c *Contract, m uint16, pubs []keypair.PublicKey, user *account.Account) bool {
	tx, err := ctx.Ont.Rpc.NewNativeInvokeTransaction(
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

	err = sdkcom.MultiSignToTransaction(tx, m, pubs, user)
	if err != nil {
		ctx.LogError("MultiSignToTransaction error: %s", err)
		return false
	}
	txHash, err := ctx.Ont.Rpc.SendRawTransaction(tx)
	if err != nil {
		ctx.LogError("SendRawTransaction error: %s", err)
	}
	return getEvent(ctx, txHash)
}
