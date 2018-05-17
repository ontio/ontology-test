package wasmvm

import (
	"errors"
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"github.com/ontio/ontology/smartcontract/types"
	"io/ioutil"
	"time"
)

const (
	filePath = "test_data"
)

func TestWasmJsonContract(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestBlockApi wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := deployWasmJsonContract(ctx, admin)
	if err != nil {
		ctx.LogError("TestWasmJsonContract deploy error:%s", err)
		return false
	}

	address, err := GetWasmContractAddress(filePath + "/contract.wasm")
	ctx.LogInfo(fmt.Sprintf("address is %s\n", address.ToHexString()))
	if err != nil {
		ctx.LogError("TestWasmJsonContract GetWasmContractAddress error:%s", err)
		return false
	}
	txHash, err = callAdd(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract invokeContract error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract init GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallWasmJsonContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestWasmJsonContract callAdd return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callAdd ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callconcat(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callconcat error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callconcat GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallWasmJsonContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestWasmJsonContract callconcat return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callconcat ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callAddStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract invokeAddStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract invokeAddStorage GetSmartContractEvent error:%s", err)
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestWasmJsonContract invokeAddStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract invokeAddStorage ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callGetStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callGetStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callGetStorage GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallWasmJsonContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestWasmJsonContract callGetStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callGetStorage ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callDeleteStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callDeleteStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callDeleteStorage GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallWasmJsonContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestWasmJsonContract callDeleteStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callDeleteStorage ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callGetStorage(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callGetStorage error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callGetStorage GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallWasmJsonContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestWasmJsonContract callGetStorage return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callGetStorage ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = callSumarray(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callSumarray error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract callSumarray GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallWasmJsonContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestWasmJsonContract callSumarray return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestWasmJsonContract callSumarray ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true
}

func callAdd(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address) (common.Uint256, error) {
	method := "add"
	params := make([]interface{}, 2)
	params[0] = 20
	params[1] = 30
	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callAddStorage(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address) (common.Uint256, error) {
	method := "addStorage"
	params := make([]interface{}, 2)
	params[0] = "TestKey"
	params[1] = "Hello World"
	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callconcat(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address) (common.Uint256, error) {
	method := "concat"
	params := make([]interface{}, 2)
	params[0] = "TestKey"
	params[1] = "Hello World"
	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callGetStorage(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address) (common.Uint256, error) {
	method := "getStorage"
	params := make([]interface{}, 1)
	params[0] = "TestKey"
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
func callDeleteStorage(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address) (common.Uint256, error) {
	method := "deleteStorage"
	params := make([]interface{}, 1)
	params[0] = "TestKey"
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callSumarray(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address) (common.Uint256, error) {
	method := "sumArray"
	params := make([]interface{}, 2)
	params[0] = []int{1, 2, 3, 4}
	params[1] = []int{5, 6, 7, 8}

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func deployWasmJsonContract(ctx *testframework.TestFrameworkContext, signer *account.Account) (common.Uint256, error) {

	code, err := ioutil.ReadFile(filePath + "/" + "contract.wasm")
	if err != nil {
		return common.Uint256{}, errors.New("")
	}

	codeHash := common.ToHexString(code)

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(0, 0,
		signer,
		types.WASMVM,
		true,
		codeHash,
		"wjc",
		"1.0",
		"test",
		"",
		"",
	)

	if err != nil {
		return common.Uint256{}, fmt.Errorf("DeploySmartContract error:%s", err)
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
