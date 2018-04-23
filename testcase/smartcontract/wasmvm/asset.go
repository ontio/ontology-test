package wasmvm

import (
	"github.com/ontio/ontology/common"
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"math/big"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)



func TestAssetContract(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "wallet.dat"
	wasmWalletPwd := "123456"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAssetContract wallet.GetDefaultAccount error:%s", err)
		return false
	}


	txHash, err := DeployWasmJsonContract(ctx,admin,filePath + "/asset.wasm","tcoin","1.0")

	if err != nil {
		ctx.LogError("TestAssetContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestAssetContract deploy TxHash:%x", txHash)

	address ,err := GetWasmContractAddress(filePath + "/asset.wasm")
	if err != nil{
		ctx.LogError("TestAssetContract GetWasmContractAddress error:%s", err)
		return false
	}

	txHash,err = invokeInit(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestAssetContract invokeInit error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestAssetContract init invokeInit error:%s", err)
		return false
	}
	ctx.LogInfo("TestAssetContract invoke notify %s", notifies)
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestAssetContract init invokeInit error:%s", err)
		return false
	}

	txHash,err = invokeTotalSupply(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestAssetContract invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestAssetContract init invokeTotalSupply error:%s", err)
		return false
	}
	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestAssetContract init invokeTotalSupply error:%s", err)
		return false
	}

	txHash,err = invokeBalanceOf(ctx,admin,address,"00000001")
	if err != nil {
		ctx.LogError("TestAssetContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}

	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}
	txHash,err = invokeTransfer(ctx,admin,address,"00000001","00000002",20000)
	if err != nil {
		ctx.LogError("TestAssetContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}

	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}

	txHash,err = invokeBalanceOf(ctx,admin,address,"00000001")
	if err != nil {
		ctx.LogError("TestAssetContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}

	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}


	txHash,err = invokeBalanceOf(ctx,admin,address,"00000002")
	if err != nil {
		ctx.LogError("TestAssetContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}

	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestAssetContract init invokeBalanceOf error:%s", err)
		return false
	}

	return true
}


func invokeInit(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "init"
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeTotalSupply(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "totalSupply"
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeBalanceOf(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,accountaddress string) (common.Uint256, error) {
	method := "balanceOf"
	params := make([]interface{},1)
	params[0] = accountaddress

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeTransfer(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,from,to string,amount int) (common.Uint256, error) {
	method := "transfer"
	params := make([]interface{},3)
	params[0] = from
	params[1] = to
	params[2] = amount

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
