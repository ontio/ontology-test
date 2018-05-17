package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)

func TestCallNativeContractJson(ctx *testframework.TestFrameworkContext) bool {
	//TODO
	fileName := ""
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestBlockApi wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/"+fileName, "CNC", "1.0")

	if err != nil {
		ctx.LogError("TestCallNativeContractJson deploy error:%s", err)
		return false
	}

	address, err := GetWasmContractAddress(filePath + "/" + fileName)
	if err != nil {
		ctx.LogError("TestCallNativeContractJson GetWasmContractAddress error:%s", err)
		return false
	}
	txHash, err = invokeTransferOntJson(ctx, admin, address, "TA4hGJWMawMQKRWFQKGcNs9YFn8Efj8zPq", 40000)
	if err != nil {
		ctx.LogError("TestCallNativeContractJson invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNativeContract init invokeTransferOnt error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallNativeContract contract invoke failed state:0")
		return false
	}
	bs, _ := common.HexToBytes(notifies.Notify[0].States[0].(string))
	if bs == nil {
		ctx.LogError("TestAssetContract init invokeTotalSupply error:%s", err)
		return false
	}

	return true
}

func invokeTransferOntJson(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, to string, amount int64) (common.Uint256, error) {
	method := "transferont"
	params := make([]interface{}, 2)
	params[0] = to
	params[1] = amount

	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
