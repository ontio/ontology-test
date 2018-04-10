package wasmvm

import (
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"math/big"
	"github.com/ontio/ontology/smartcontract/service/wasm"
	"time"
	"fmt"
)

func TestWasmRawContract(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "/home/zhoupw/work/go/src/github.com/ontio/ontology/wallet.dat"
	wasmWalletPwd := "123456"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestWasmRawContract wallet.GetDefaultAccount error:%s", err)
		return false
	}


	txHash, err := DeployWasmJsonContract(ctx,admin,filePath + "/rawcontract.wasm","rwc","1.0")

	if err != nil {
		ctx.LogError("TestWasmRawContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestWasmRawContract deploy TxHash:%x", txHash)

	address ,err := GetWasmContractAddress(filePath + "/rawcontract.wasm")
	if err != nil{
		ctx.LogError("TestWasmRawContract GetWasmContractAddress error:%s", err)
		return false
	}

	txHash,err = invokeRawContract(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestWasmRawContract invokeContract error:%s", err)
		return false
	}

	ctx.LogInfo("invokeContract: %x\n", txHash)
	ctx.LogInfo("TestWasmRawContract invokeContract success")
	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmRawContract init GetSmartContractEvent error:%s", err)
		return false
	}
	ctx.LogInfo("TestWasmRawContract invoke notify %s", notifies)
	fmt.Println("============result is===============")
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))


	return true
}

func invokeRawContract(ctx *testframework.TestFrameworkContext, acc *account.Account,contractAddress common.Address) (common.Uint256, error) {
	method := "add"
	params := make([]interface{},2)
	params[0] = 20
	params[1] = 30
	txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),contractAddress,method,wasm.Raw,params,1,false)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
