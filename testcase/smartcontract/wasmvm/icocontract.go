package wasmvm

import (
	"github.com/ontio/ontology/common"
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)



func TestICOContract(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "wallet.dat"
	wasmWalletPwd := "123456"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestICOContract wallet.GetDefaultAccount error:%s", err)
		return false
	}


	txHash, err := DeployWasmJsonContract(ctx,admin,filePath + "/icotest.wasm","tcoin","1.0")

	if err != nil {
		ctx.LogError("TestICOContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestICOContract deploy TxHash:%x", txHash)

	address ,err := GetWasmContractAddress(filePath + "/icotest.wasm")
	ctx.LogInfo("contract b58address is %s\n",address.ToBase58())
	if err != nil{
		ctx.LogError("TestICOContract GetWasmContractAddress error:%s", err)
		return false
	}

	txHash,err = invokeICOInit(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestICOContract invokeInit error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeInit error:%s", err)
		return false
	}
	if len(notifies) < 1{
		ctx.LogError("TestICOContract return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeInit ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	txHash,err = invokeICOTotalSupply(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestICOContract invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract invokeTotalSupply error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestICOContract invokeTotalSupply return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeTotalSupply ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}



	//collect
	//suppose the address TA4hGJWMawMQKRWFQKGcNs9YFn8Efj8zPq has enough ont
	txHash,err = invokeICOCollect(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB",100)
	if err != nil {
		ctx.LogError("TestICOContract invokeICOCollect error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract invokeICOCollect error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestICOContract invokeICOCollect return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeICOCollect ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = invokeICOBalanceOf(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB")
	if err != nil {
		ctx.LogError("TestICOContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestICOContract invokeBalanceOf return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeBalanceOf TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	txHash,err = invokeICOTransfer(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB","TA4ieHoEDmRmARQo6bVBayqPuvN51rd6wY",20)
	if err != nil {
		ctx.LogError("TestICOContract invokeICOTransfer error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeICOTransfer error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestICOContract invokeICOTransfer return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeICOTransfer ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = invokeICOBalanceOf(ctx,admin,address,"TA4ieHoEDmRmARQo6bVBayqPuvN51rd6wY")
	if err != nil {
		ctx.LogError("TestICOContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}
	if len(notifies) < 1{
		ctx.LogError("TestICOContract invokeBalanceOf return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeBalanceOf TA4ieHoEDmRmARQo6bVBayqPuvN51rd6wY============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	txHash,err = invokeICOBalanceOf(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB")
	if err != nil {
		ctx.LogError("TestICOContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestICOContract invokeBalanceOf return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeBalanceOf TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	return true
}


func invokeICOInit(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "init"
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0,0,acc,1,address,method, wasmvm.Json,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeICOTotalSupply(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "totalSupply"
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0,0,acc,1,address,method, wasmvm.Json,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeICOBalanceOf(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,accountaddress string) (common.Uint256, error) {
	method := "balanceOf"
	params := make([]interface{},1)
	params[0] = accountaddress

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0,0,acc,1,address,method, wasmvm.Json,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeICOTransfer(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,from,to string,amount int64) (common.Uint256, error) {
	method := "transfer"
	params := make([]interface{},3)
	params[0] = from
	params[1] = to
	params[2] = amount

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0,0,acc,1,address,method, wasmvm.Json,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeICOCollect(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,from string,amount int64) (common.Uint256, error) {
	method := "collect"
	params := make([]interface{},2)
	params[0] = from
	params[1] = amount

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0,0,acc,1,address,method, wasmvm.Json,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeICOWithdraw(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,amount int64) (common.Uint256, error) {
	method := "withdraw"
	params := make([]interface{},1)
	params[0] = amount

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0,0,acc,1,address,method, wasmvm.Json,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
