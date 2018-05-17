package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)

func TestDomainContract_Invoke(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestDomainContract wallet.GetDefaultAccount error:%s", err)
		return false
	}
	address, err := GetWasmContractAddress(filePath + "/domain.wasm")

	//current Price
	txHash, err := invokeDomainCurrentPrice(ctx, admin, address, "www.goodthings.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestDomainContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainCurrentPrice ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	//query
	txHash, err = invokeDomainQuery(ctx, admin, address, "www.goodthings.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainQuery error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestDomainContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestDomainContract invokeDomainQuery return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeDomainQuery ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	//Buy
	txHash, err = invokeDomainBuy(ctx, admin, address, "TA8Xe297g4wGj67maMYZFmdfk9i2riVNrC", "www.goodthings.com", 150)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainBuy error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainBuy error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestDomainContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestDomainContract invokeDomainBuy return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainBuy ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	//current Price
	txHash, err = invokeDomainCurrentPrice(ctx, admin, address, "www.goodthings.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestDomainContract contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainCurrentPrice ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true
}

func invokeDomainBuy(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, from string, domain string, basePrice int64) (common.Uint256, error) {
	method := "buy"

	params := make([]interface{}, 3)
	params[0] = from
	params[1] = domain
	params[2] = basePrice

	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeDomainCurrentPrice(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, domain string) (common.Uint256, error) {
	method := "currentPrice"

	params := make([]interface{}, 1)
	params[0] = domain

	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
