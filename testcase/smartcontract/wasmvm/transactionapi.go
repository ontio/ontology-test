package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)

func TestTransactionApi(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestHeaderApi wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/transactionapi.wasm", "testheaderapi", "1.0")

	if err != nil {
		ctx.LogError("TestTransactionApi deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestTransactionApi deploy TxHash:%x", txHash)

	address, err := GetWasmContractAddress(filePath + "/transactionapi.wasm")
	fmt.Println(address.ToHexString())
	if err != nil {
		ctx.LogError("TestHeaderApi GetWasmContractAddress error:%s", err)
		return false
	}

	txHash, err = getTransactionType(ctx, admin, address, common.ToHexString(txHash.ToArray()))
	if err != nil {
		ctx.LogError("TestTransactionApi getTransactionType error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestTransactionApi getTransactionType GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestTransactionApi getTransactionType invoke failed, state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestTransactionApi getTransactionType return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestTransactionApi getTransactionType ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getTransactionAttributes(ctx, admin, address, common.ToHexString(txHash.ToArray()))
	if err != nil {
		ctx.LogError("TestTransactionApi getTransactionAttributes error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestTransactionApi getTransactionAttributes GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestTransactionApi getTransactionAttributes invoke failed, state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestTransactionApi getTransactionAttributes return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestTransactionApi getTransactionAttributes ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true
}

func getTransactionType(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getTransactionType"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getTransactionAttributes(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getTransactionAttributes"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}
