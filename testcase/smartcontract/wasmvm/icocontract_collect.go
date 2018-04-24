package wasmvm

import (
	"github.com/ontio/ontology/common"
	"fmt"
	"github.com/ontio/ontology-test/testframework"
)



func TestICOContractCollect(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "testwallet.dat"
	wasmWalletPwd := "123456"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestICOContractCollect wallet.GetDefaultAccount error:%s", err)
		return false
	}
	address ,err := GetWasmContractAddress(filePath + "/icotest.wasm")
	if err != nil{
		ctx.LogError("TestICOContract GetWasmContractAddress error:%s", err)
		return false
	}



	//collect
	txHash,err := invokeICOCollect(ctx,admin,address,"TA4hGJWMawMQKRWFQKGcNs9YFn8Efj8zPq",500)
	if err != nil {
		ctx.LogError("TestICOContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}

	bs ,_ := common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}

	fmt.Printf("collect result is %s\n",bs)


	txHash,err = invokeICOBalanceOf(ctx,admin,address,"TA4hGJWMawMQKRWFQKGcNs9YFn8Efj8zPq")
	if err != nil {
		ctx.LogError("TestICOContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}

	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}

	fmt.Printf("balance of %s is %s\n","TA4hGJWMawMQKRWFQKGcNs9YFn8Efj8zPq",bs)


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

	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}

	fmt.Printf("balance of %s is %s\n","TA4ieHoEDmRmARQo6bVBayqPuvN51rd6wY",bs)


	txHash,err = invokeICOWithdraw(ctx,admin,address,100)
	if err != nil {
		ctx.LogError("TestICOContract invokeICOWithdraw error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeICOWithdraw error:%s", err)
		return false
	}

	bs ,_= common.HexToBytes(notifies[0].States[0].(string))
	if bs == nil{
		ctx.LogError("TestICOContract init invokeICOWithdraw error:%s", err)
		return false
	}

	fmt.Printf("invokeICOWithdraw  is %s\n",bs)


	return true
}


