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



func TestDomainContract(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "wallet.dat"
	wasmWalletPwd := "123456"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestDomainContract wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/domain.wasm", "domain", "1.0")

	if err != nil {
		ctx.LogError("TestDomainContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestDomainContract deploy TxHash:%x", txHash)

	address, err := GetWasmContractAddress(filePath + "/domain.wasm")
	ctx.LogInfo("contract b58address is %s\n", address.ToBase58())
	if err != nil {
		ctx.LogError("TestDomainContract GetWasmContractAddress error:%s", err)
		return false
	}

	txHash,err = invokeDomainRegister(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB","www.letsrock.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainRegister error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainRegister error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainRegister return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainRegister ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = invokeDomainRegister(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB","www.letsrock2.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainRegister error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainRegister error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainRegister return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainRegister ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}



	txHash,err = invokeDomainQuery(ctx,admin,address,"www.letsrock.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainQuery return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeDomainQuery ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = invokeDomainTransfer(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB","TA8Xe297g4wGj67maMYZFmdfk9i2riVNrC","www.letsrock.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainTransfer return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeDomainTransfer ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = invokeDomainQuery(ctx,admin,address,"www.letsrock.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainQuery return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeDomainQuery ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	txHash,err = invokeDomainDelete(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB","www.letsrock2.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainDelete error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainDelete error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainDelete return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeDomainDelete ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = invokeDomainQuery(ctx,admin,address,"www.letsrock2.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainQuery return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeDomainQuery ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	return true
}

func invokeDomainRegister(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,regAddress string,domain string) (common.Uint256, error) {
	method := "register"

	params := make([]interface{},2)
	params[0] = regAddress
	params[1] = domain

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeDomainQuery(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,domain string) (common.Uint256, error) {
	method := "query"

	params := make([]interface{},1)
	params[0] = domain

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeDomainTransfer(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,from,to string,domain string) (common.Uint256, error) {
	method := "transfer"

	params := make([]interface{},3)
	params[0] = from
	params[1] = to
	params[2] = domain

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeDomainDelete(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,addr,domain string) (common.Uint256, error) {
	method := "delete"

	params := make([]interface{},2)
	params[0] = addr
	params[1] = domain

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}