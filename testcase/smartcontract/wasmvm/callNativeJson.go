package wasmvm

import (
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"math/big"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
	"fmt"
)

func TestCallNativeContractJson(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "wallet.dat"
	wasmWalletPwd := "123456"
	fileName := "callNativeJson.wasm"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestCallNativeContractJson wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx,admin,filePath + "/" + fileName,"CNC","1.0")

	if err != nil {
		ctx.LogError("TestCallNativeContractJson deploy error:%s", err)
		return false
	}

	address ,err := GetWasmContractAddress(filePath + "/" + fileName)
	if err != nil{
		ctx.LogError("TestCallNativeContractJson GetWasmContractAddress error:%s", err)
		return false
	}
	txHash,err = invokeTransferOntJson(ctx,admin,address,"TA4hGJWMawMQKRWFQKGcNs9YFn8Efj8zPq",40000)
	if err != nil {
		ctx.LogError("TestCallNativeContractJson invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNativeContract init invokeTransferOnt error:%s", err)
		return false
	}
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestAssetContract init invokeTotalSupply error:%s", err)
		return false
	}



	return true
}

func invokeTransferOntJson(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,to string,amount int64) (common.Uint256, error) {
	method := "transferont"
	params := make([]interface{},2)
	params[0] = to
	params[1] = amount

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
