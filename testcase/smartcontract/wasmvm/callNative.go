package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)

func TestCallNativeContract(ctx *testframework.TestFrameworkContext) bool {
	//TODO
	fileName := ""
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestBlockApi wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/"+fileName, "CNC", "1.0")

	if err != nil {
		ctx.LogError("TestCallNativeContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestCallNativeContract deploy TxHash:%x", txHash)

	address, err := GetWasmContractAddress(filePath + "/" + fileName)
	if err != nil {
		ctx.LogError("TestCallNativeContract GetWasmContractAddress error:%s", err)
		return false
	}
	//trans ont
	txHash, err = invokeTransferOnt(ctx, admin, address, "TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB", "TA8Xe297g4wGj67maMYZFmdfk9i2riVNrC", 5000)
	if err != nil {
		ctx.LogError("TestCallNativeContract invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNativeContract invokeTransferOnt error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallNativeContract contract invoke failed state:0")
		return false
	}
	ctx.LogInfo("==========TestCallNativeContract TestCallNativeContract ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	//trans ont
	txHash, err = invokeTransferOnt(ctx, admin, address, "TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB", "TA8aqS3PyDcFG567qa2qJuufHH1M82zVig", 5000)
	if err != nil {
		ctx.LogError("TestCallNativeContract invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNativeContract invokeTransferOnt error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallNativeContract contract invoke failed state:0")
		return false
	}
	ctx.LogInfo("==========TestCallNativeContract TestCallNativeContract ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true
}

func invokeTransferOnt(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, from, to string, amount int64) (common.Uint256, error) {
	method := "transferont"
	params := make([]interface{}, 3)
	params[0] = from
	params[1] = to
	params[2] = amount

	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(), acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
