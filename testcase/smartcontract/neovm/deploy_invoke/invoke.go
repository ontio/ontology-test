package deploy_invoke

import (
	"github.com/ontio/ontology-test/testframework"
	"time"
)

func TestInvokeSmartContract(ctx *testframework.TestFrameworkContext) bool {
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestInvokeSmartContract GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(), signer, contractCodeAddress, []interface{}{})
	if err != nil {
		ctx.LogError("TestInvokeSmartContract InvokeNeoVMSmartContract error:%s", err)
		return false
	}

	ctx.LogInfo("TestInvokeSmartContract txHash:%s\n", txHash.ToHexString())
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
	if events.State == 0 {
		ctx.LogError("TestInvokeSmartContract failed invoked exec state return 0")
		return false
	}
	notify := events.Notify[0]
	ctx.LogInfo("%+v", notify)

	invokeState := notify.States.([]interface{})
	//Event name
	name, _ := ctx.ConvertToHexString(invokeState[0])
	err = ctx.AssertToString(name, "transfer")
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToString:%s", err)
		return false
	}

	//key of event
	key, _ := ctx.ConvertToHexString(invokeState[1])
	err = ctx.AssertToString(key, "hello")
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToString %s ", err)
		return false
	}

	//value of event
	value, _ := ctx.ConvertToHexString(invokeState[2])
	err = ctx.AssertToString(value, "world")
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToString %s ", err)
		return false
	}

	//amount of event
	amount, _ := ctx.ConvertToBigInt(invokeState[3])
	err = ctx.AssertToInt(amount, 123)
	if err != nil {
		ctx.LogError("TestInvokeSmartContract failed AssertToInt %s ", err)
		return false
	}

	return true
}
