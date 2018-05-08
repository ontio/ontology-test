package wasmvm

import (
	"github.com/ontio/ontology-test/testframework"
	"fmt"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"math/big"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
	"encoding/json"
	"strconv"
)

type Result struct {
	Ptype    string `json:"type"`
	Pval     string `json:"value"`
	Psucceed int    `json:"succeed"`
}

var blockhash string
var height = 1
func TestBlockApi(ctx *testframework.TestFrameworkContext) bool {
	wasmWallet := "wallet.dat"
	wasmWalletPwd := "123456"

	wallet, err := ctx.Ont.OpenWallet(wasmWallet, wasmWalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", wasmWallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestBlockApi wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/blockapi.wasm", "testblockapi", "1.0")

	if err != nil {
		ctx.LogError("TestBlockApi deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestBlockApi deploy TxHash:%x", txHash)

	address, err := GetWasmContractAddress(filePath + "/blockapi.wasm")
	fmt.Println(address.ToHexString())
	if err != nil {
		ctx.LogError("TestBlockApi GetWasmContractAddress error:%s", err)
		return false
	}

	txHash,err = callGetHeaderHeight(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestBlockApi callGetHeaderHeight error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetHeaderHeight GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetHeaderHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetHeaderHeight ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	txHash,err = callGetHeaderHash(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestBlockApi callGetHeaderHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetHeaderHash GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetHeaderHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetHeaderHash ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	txHash,err = callGetBlockHeight(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestBlockApi callGetBlockHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetBlockHeight GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetBlockHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetBlockHeight ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	ret := &Result{}
	s := notifies[0].States[0].(string)
	err = json.Unmarshal([]byte(s),ret)
	if err!= nil{
		fmt.Printf("error is %s\n",err.Error())
	}
	height,err = strconv.Atoi(ret.Pval)
	if err!= nil{
		fmt.Printf("error is %s\n",err.Error())
	}

	txHash,err = callGetBlockHash(ctx,admin,address)
	if err != nil {
		ctx.LogError("TestBlockApi callGetBlockHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetBlockHash GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetBlockHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetBlockHash ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}
	ret = &Result{}
	s = notifies[0].States[0].(string)
	err = json.Unmarshal([]byte(s),ret)
	if err!= nil{
		fmt.Printf("error is %s\n",err.Error())
	}
	blockhash = ret.Pval

	txHash,err = callGetTransByHash(ctx,admin,address,common.ToHexString(txHash.ToArray()))
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransByHash GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetTransByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetTransByHash ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}





	txHash,err = callGetTransCount(ctx,admin,address,blockhash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransByHash GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetTransByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetTransByHash ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}


	txHash,err = callGetTransCountByHeight(ctx,admin,address,height)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransCountByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransCountByHeight GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetTransCountByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetTransCountByHeight ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = callGetTransactionsByHash(ctx,admin,address,blockhash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransactionsByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransactionsByHash GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetTransactionsByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetTransactionsByHash ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}



	txHash,err = callGetTransCountByHeight(ctx,admin,address,height)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransCountByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransCountByHeight GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetTransCountByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetTransCountByHeight ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	txHash,err = callGetTransactionsByHeight(ctx,admin,address,height)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransactionsByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestBlockApi callGetTransactionsByHeight GetSmartContractEvent error:%s", err)
		return false
	}

	if len(notifies) < 1{
		ctx.LogError("TestBlockApi callGetTransactionsByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestBlockApi callGetTransactionsByHeight ============")
	for i ,n := range notifies{
		ctx.LogInfo(fmt.Sprintf("notify %d is %v",i, n))
	}

	return true
}

func callGetHeaderHash(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "getCurrentHeadHash"
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callGetHeaderHeight(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "getCurrentHeaderHeight"

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
func callGetBlockHash(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "getCurrentBlockHash"

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callGetBlockHeight(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address) (common.Uint256, error) {
	method := "getCurrentBlockHeight"

	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,nil)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
func callGetTransByHash(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,hash string) (common.Uint256, error) {
	method := "getTransactionByHash"
	params := make([]interface{},1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
func callGetTransCount(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,hash string) (common.Uint256, error) {
	method := "getTransactionCountByHash"
	params := make([]interface{},1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callGetTransCountByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,height int) (common.Uint256, error) {
	method := "getTransactionCountByHeight"
	params := make([]interface{},1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callGetTransactionsByHash(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,hash string) (common.Uint256, error) {
	method := "getTransactions"
	params := make([]interface{},1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func callGetTransactionsByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account,address common.Address,height int) (common.Uint256, error) {
	method := "getTransactionsByHeight"
	params := make([]interface{},1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash,err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(acc,new(big.Int),address,method, wasmvm.Raw,1,params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}