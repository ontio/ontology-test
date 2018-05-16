package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)

func TestWasmRawContract(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestWasmRawContract wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/rawcontract.wasm", "rwc", "1.0")

	if err != nil {
		ctx.LogError("TestWasmRawContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestWasmRawContract deploy TxHash:%x", txHash)

	address, err := GetWasmContractAddress(filePath + "/rawcontract.wasm")
	if err != nil {
		ctx.LogError("TestWasmRawContract GetWasmContractAddress error:%s", err)
		return false
	}

	txHash, err = callRawContractAdd(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractAdd error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractAdd GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1 {
		ctx.LogError("TestWasmJsonContract callRawContractAdd return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callGetStorage ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callRawContractAddStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractAddStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractAddStorage GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1 {
		ctx.LogError("TestWasmJsonContract callRawContractAddStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callRawContractAddStorage ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callRawContractGetStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractGetStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractGetStorage GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1 {
		ctx.LogError("TestWasmJsonContract callRawContractGetStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callRawContractGetStorage ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callRawContractDeleteStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractDeleteStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractDeleteStorage GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1 {
		ctx.LogError("TestWasmJsonContract callRawContractDeleteStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callRawContractDeleteStorage ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callRawContractGetStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractGetStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmRawContract callRawContractGetStorage GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1 {
		ctx.LogError("TestWasmJsonContract callRawContractGetStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callRawContractGetStorage ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true
}

func callRawContractAdd(ctx *testframework.TestFrameworkContext, acc *account.Account, contractAddress common.Address) (common.Uint256, error) {
	method := "add"
	params := make([]interface{}, 2)
	params[0] = 20
	params[1] = 30
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, contractAddress, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callRawContractAddStorage(ctx *testframework.TestFrameworkContext, acc *account.Account, contractAddress common.Address) (common.Uint256, error) {
	method := "addStorage"
	params := make([]interface{}, 2)
	params[0] = "TestKey"
	params[1] = "Hello World"
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, contractAddress, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callRawContractGetStorage(ctx *testframework.TestFrameworkContext, acc *account.Account, contractAddress common.Address) (common.Uint256, error) {
	method := "getStorage"
	params := make([]interface{}, 1)
	params[0] = "TestKey"
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, contractAddress, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callRawContractDeleteStorage(ctx *testframework.TestFrameworkContext, acc *account.Account, contractAddress common.Address) (common.Uint256, error) {
	method := "deleteStorage"
	params := make([]interface{}, 1)
	params[0] = "TestKey"
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, contractAddress, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
