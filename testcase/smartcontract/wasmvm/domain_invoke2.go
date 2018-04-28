package wasmvm

import (
	"github.com/ontio/ontology-test/testframework"
	"fmt"
)



func TestDomainContract_Invoke2(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "./test_data/testwallet2.dat"
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
	address, err := GetWasmContractAddress(filePath + "/domain.wasm")


	//current Price
	txHash,err := invokeDomainCurrentPrice(ctx,admin,address,"www.goodthings.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainCurrentPrice ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	//query
	txHash,err = invokeDomainQuery(ctx,admin,address,"www.goodthings.com")
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

	//Buy
	txHash,err = invokeDomainBuy(ctx,admin,address,"TA8aqS3PyDcFG567qa2qJuufHH1M82zVig","www.goodthings.com",200)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainBuy error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainBuy error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainBuy return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainBuy ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	//current Price
	txHash,err = invokeDomainCurrentPrice(ctx,admin,address,"www.goodthings.com")
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestDomainContract invokeDomainCurrentPrice return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestDomainContract invokeDomainCurrentPrice ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	return true
}

