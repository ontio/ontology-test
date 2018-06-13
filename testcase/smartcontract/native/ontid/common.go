package ontid

import (
	"bytes"
	"errors"
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/states"
)

func makeTx(contract *states.Contract) (*types.Transaction, error) {
	buf := new(bytes.Buffer)
	if err := contract.Serialize(buf); err != nil {
		return nil, errors.New("Serialize contract error: " + err.Error())
	}

	return sdkcom.NewInvokeTransaction(common.DefConfig.GasPrice, common.DefConfig.GasLimit, buf.Bytes()), nil
}

func sendTx(ctx *testframework.TestFrameworkContext, invokeTx *types.Transaction) bool {
	txHash, err := ctx.Ont.Rpc.SendRawTransaction(invokeTx)
	if err != nil {
		ctx.LogError("SendTransaction error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 2)
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
		ctx.LogError("ontio contract invoke failed, state:0")
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

func preSendTx(ctx *testframework.TestFrameworkContext, invokeTx *types.Transaction) []byte {
	var buf bytes.Buffer
	err := invokeTx.Serialize(&buf)
	if err != nil {
		ctx.LogError("serialize tx error: %s", err)
		return nil
	}

	//txData := hex.EncodeToString(buf.Bytes())
	//data, err := ctx.Ont.Rpc.sendRpcRequest(RPC_SEND_TRANSACTION, []interface{}{txData, 1})
	//if err != nil {
	//	ctx.LogError("SendTransaction error: %s", err)
	//	return nil
	//}
	//return data
	return nil
}

func InvokeContract(ctx *testframework.TestFrameworkContext, contract *states.Contract, preExec bool) (bool, []byte) {
	tx, err := makeTx(contract)
	if err != nil {
		ctx.LogError("make transactions error: %s", err)
		return false, nil
	}
	user, _ := ctx.GetDefaultAccount()
	err = sdkcom.SignToTransaction(tx, user)
	if err != nil {
		ctx.LogError("sign transaction error: %s", err)
		return false, nil
	}
	if preExec {
		res := preSendTx(ctx, tx)
		return (res != nil), res
	} else {
		res := sendTx(ctx, tx)
		return res, nil
	}
}
