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
	txHash,err = invokeTransferOnt(ctx,admin,address,"TA4tBPFEn7Amutm7QWTBYesEHE5sbWZKsB","TA8Xe297g4wGj67maMYZFmdfk9i2riVNrC",4000)
	if err != nil {
		ctx.LogError("TestCallNativeContract invokeTotalSupply error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNativeContract invokeTransferOnt error:%s", err)
		return false
	}

	if len(notifies) != 2{
		ctx.LogError("TestCallNativeContract invokeTransferOnt error:notifies length errors" )
		return false
	}

	ctx.LogInfo(fmt.Sprintf("TestCallNativeContract invokeTransferOnt notify %v\n", notifies))
	ctx.LogInfo("============invokeTotalSupply result is===============")
	ctx.LogInfo(fmt.Sprintf("notifies[0]:%v\n",notifies[0]))
	ctx.LogInfo(fmt.Sprintf("notifies[1]:%v\n",notifies[1]))
	return true
}

func invokeTransferOnt(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,from,to string,amount int64) (common.Uint256, error) {
	method := "transferont"
	params := make([]interface{},3)
	params[0] = from
	params[1] = to
	params[2] = amount

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
