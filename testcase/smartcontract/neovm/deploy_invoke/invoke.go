package deploy_invoke

import (
	"math/big"
	"time"

	"github.com/ontio/ontology-test/testframework"
)

func TestInvokeSmartContract(ctx *testframework.TestFrameworkContext) bool {
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestInvokeSmartContract GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer, new(big.Int), contractCodeAddress, []interface{}{})
	if err != nil {
		ctx.LogError("TestInvokeSmartContract InvokeNeoVMSmartContract error:%s", err)
	}

	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestInvokeSmartContract WaitForGenerateBlock error:%s", err)
		return false
	}

	//Test GetStorageItem api
	skey := "Hello"
	svalue, err := ctx.Ont.Rpc.GetStorage(contractCodeAddress, []byte(skey))
	if err != nil {
		ctx.LogError("TestInvokeSmartContract GetStorageItem key:%s error:%s", skey, err)
		return false
	}
	err = ctx.AssertToString(string(svalue), "World")
	if err != nil {
		ctx.LogError("TestInvokeSmartContract AssertToString error:%s", err)
		return false
	}

	//GetEventLog, to check the result of invoke
	events, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestInvokeSmartContract GetSmartContractEvent error:%s", err)
		return false
	}
	transfer := events[0].States
	ctx.LogInfo("%+v", transfer)

	//Event name
	name, _ := ctx.ConvertToHexString(transfer[0])
	err = ctx.AssertToString(name, "transfer")
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToString:%s", err)
		return false
	}

	//key of event
	key, _ := ctx.ConvertToHexString(transfer[1])
	err = ctx.AssertToString(key, "hello")
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToString %s ", err)
		return false
	}

	//value of event
	value, _ := ctx.ConvertToHexString(transfer[2])
	err = ctx.AssertToString(value, "world")
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToString %s ", err)
		return false
	}

	//amount of event
	amount, _ := ctx.ConvertToBigInt(transfer[3])
	err = ctx.AssertToInt(amount, 123)
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToInt %s ", err)
		return false
	}

	return true
}
