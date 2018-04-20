package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"io/ioutil"
	"errors"
	"github.com/ontio/ontology/smartcontract/types"
	"time"
	"math/big"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
)

func TestCallWasmJsonContract(ctx *testframework.TestFrameworkContext) bool{
	wasmWallet := "/home/zhoupw/work/go/src/github.com/ontio/ontology/wallet.dat"
	wasmWalletPwd := "123456"
	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestCallWasmJsonContract wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := deployCallWasmJsonContract(ctx, admin)
	if err != nil {
		ctx.LogError("TestCallWasmJsonContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestCallWasmJsonContract deploy TxHash:%x", txHash)

	address,err := GetWasmContractAddress(filePath+"/callContract.wasm")
	fmt.Printf("TestCallWasmJsonContract address is %s\n ",address.ToHexString())
	if err != nil {
		ctx.LogError("TestCallWasmJsonContract GetWasmContractAddress error:%s", err)
		return false
	}
	txHash,err = invokeCallContract(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestCallWasmJsonContract invokeCallContract error:%s", err)
		return false
	}
	ctx.LogInfo("invokeContract: %x\n", txHash)
	ctx.LogInfo("TestCallWasmJsonContract invokeCallContract success")
	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallWasmJsonContract init GetSmartContractEvent error:%s", err)
		return false
	}
	fmt.Println("============result is===============")
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))
	fmt.Println("============result is===============")
	bs ,_= common.HexToBytes(notifies[1].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))


	/*txHash,err = invokeCallOffchainContract(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestCallWasmJsonContract invokeCallContract error:%s", err)
		return false
	}
	ctx.LogInfo("invokeContract: %x\n", txHash)
	ctx.LogInfo("TestCallWasmJsonContract invokeCallContract success")

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallWasmJsonContract init GetSmartContractEvent error:%s", err)
		return false
	}
	fmt.Println("============result is===============")
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))
	fmt.Println("============result is===============")
	bs ,_= common.HexToBytes(notifies[1].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))*/


	return true
}

func deployCallWasmJsonContract(ctx *testframework.TestFrameworkContext, signer *account.Account) (common.Uint256, error){

	code, err := ioutil.ReadFile(filePath + "/" + "callContract.wasm")
	if err != nil {
		return common.Uint256{}, errors.New("")
	}

	codeHash := common.ToHexString(code)

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(
		signer,
		types.WASMVM,
		true,
		codeHash,
		"cwjc",
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

func invokeCallContract(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "getValue"
	params := make([]interface{},1)
	params[0] = "TestKey"
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func invokeCallOffchainContract(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "add"
	params := make([]interface{},2)
	params[0] = 40
	params[1] = 50
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Json,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}