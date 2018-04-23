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

func TestCallNativeContract(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "wallet.dat"
	wasmWalletPwd := "123456"
	fileName := "callNative.wasm"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestCallNativeContract wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx,admin,filePath + "/" + fileName,"CNC","1.0")

	if err != nil {
		ctx.LogError("TestCallNativeContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestCallNativeContract deploy TxHash:%x", txHash)

	address ,err := GetWasmContractAddress(filePath + "/" + fileName)
	if err != nil{
		ctx.LogError("TestCallNativeContract GetWasmContractAddress error:%s", err)
		return false
	}
	txHash,err = invokeTransferOnt(ctx,admin,address,"TA7xfQvv3h6eGicv4VE4FU6NjBLpFNB9jr",400)
	if err != nil {
		ctx.LogError("TestCallNativeContract invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNativeContract init invokeTransferOnt error:%s", err)
		return false
	}
	fmt.Printf("TestCallNativeContract invokeTransferOnt notify %v\n", notifies)
	fmt.Println("============invokeTotalSupply result is===============")
	fmt.Printf("notifies[0]:%v\n",notifies[0])
	fmt.Printf("States[0]:%v\n",notifies[0].States[0])
	fmt.Printf("notifies[1]:%v\n",notifies[1])
	fmt.Printf("States[0]:%v\n",notifies[1].States[0])
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))

	return true
}

func invokeTransferOnt(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,to string,amount int64) (common.Uint256, error) {
	method := "transferont"
	params := make([]interface{},2)
	params[0] = to
	params[1] = amount

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
