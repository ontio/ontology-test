package wasmvm

import (
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/vm/types"
	"io/ioutil"
	"errors"
	"fmt"
	"time"
	"math/big"
	"github.com/ontio/ontology/smartcontract/service/wasm"
)

const (
	 filePath = "/home/zhoupw/work/go/src/github.com/ontio/ontology/vm/wasmvm/exec/test_data2"
)


func TestWasmJsonContract(ctx *testframework.TestFrameworkContext) bool{
	wasmWallet := "/home/zhoupw/work/go/src/github.com/ontio/ontology/wallet.dat"
	wasmWalletPwd := "123456"
	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestWasmJsonContract wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := deployWasmJsonContract(ctx, admin)
	if err != nil {
		ctx.LogError("TestWasmJsonContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestWasmJsonContract deploy TxHash:%x", txHash)

	address,err := GetWasmContractAddress(filePath+"/contract.wasm")
	if err != nil {
		ctx.LogError("TestWasmJsonContract GetWasmContractAddress error:%s", err)
		return false
	}
	txHash,err = invokeContract(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestWasmJsonContract invokeContract error:%s", err)
		return false
	}
	ctx.LogInfo("invokeContract: %x\n", txHash)
	ctx.LogInfo("TestWasmJsonContract invokeContract success")
	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract init GetSmartContractEvent error:%s", err)
		return false
	}
	ctx.LogInfo("TestNep5Contract init notify %s", notifies)
	fmt.Println("============result is===============")
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))
	return true
}


func invokeContract(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "add"
	params := make([]interface{},2)
	params[0] = 20
	params[1] = 30
	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method,wasm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func deployWasmJsonContract(ctx *testframework.TestFrameworkContext, signer *account.Account) (common.Uint256, error){

	code, err := ioutil.ReadFile(filePath + "/" + "contract.wasm")
	if err != nil {
		return common.Uint256{}, errors.New("")
	}

	codeHash := common.ToHexString(code)

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(
		signer,
		types.WASMVM,
		true,
		codeHash,
		"wjc",
		"1.0",
		"test",
		"",
		"",
	)

	if err != nil {
		return common.Uint256{}, fmt.Errorf("TestNep5Contract DeploySmartContract error:%s", err)
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}


