package wasmvm

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"time"
)

func TestHeaderApi(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestDomainContract wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/headerapi.wasm", "testheaderapi", "1.0")

	if err != nil {
		ctx.LogError("TestHeaderApi deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestHeaderApi deploy TxHash:%x", txHash)

	address, err := GetWasmContractAddress(filePath + "/headerapi.wasm")
	fmt.Println(address.ToHexString())
	if err != nil {
		ctx.LogError("TestHeaderApi GetWasmContractAddress error:%s", err)
		return false
	}

	txHash, err = getHeaderHashByHeight(ctx, admin, address, 1)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderHashByHeight error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderHashByHeight GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getHeaderHashByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getHeaderHashByHeight ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}
	ret := &Result{}
	s := notifies.Notify[0].States.([]interface{})[0].(string)
	err = json.Unmarshal([]byte(s), ret)
	if err != nil {
		fmt.Printf("error is %s\n", err.Error())
	}
	headerhash := ret.Pval
	fmt.Printf("header hash is %s\n", headerhash)

	txHash, err = getHeaderVersionByHeight(ctx, admin, address, 1)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderVersionByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderVersionByHeight GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getHeaderVersionByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getHeaderVersionByHeight ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getHeaderVersionByHash(ctx, admin, address, headerhash)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderVersionByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderVersionByHash GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getHeaderVersionByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getHeaderVersionByHash ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getHeaderPrevHashByHeight(ctx, admin, address, 1)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderPrevHashByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderPrevHashByHeight GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getHeaderPrevHashByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getHeaderPrevHashByHeight ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getHeaderPrevHashByHash(ctx, admin, address, headerhash)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderPrevHashByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getHeaderPrevHashByHash GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getHeaderPrevHashByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getHeaderPrevHashByHash ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getMerkelRootByHeight(ctx, admin, address, 1)
	if err != nil {
		ctx.LogError("TestHeaderApi getMerkelRootByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getMerkelRootByHeight GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getMerkelRootByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getMerkelRootByHeight ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getMerkelRootByHash(ctx, admin, address, headerhash)
	if err != nil {
		ctx.LogError("TestHeaderApi getMerkelRootByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getMerkelRootByHash GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getMerkelRootByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getMerkelRootByHash ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getTimestampByHeight(ctx, admin, address, 1)
	if err != nil {
		ctx.LogError("TestHeaderApi getTimestampByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getTimestampByHeight GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getTimestampByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getTimestampByHeight ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getTimestampByHash(ctx, admin, address, headerhash)
	if err != nil {
		ctx.LogError("TestHeaderApi getTimestampByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getTimestampByHash GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getTimestampByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getTimestampByHash ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getIndexByHash(ctx, admin, address, headerhash)
	if err != nil {
		ctx.LogError("TestHeaderApi getIndexByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getIndexByHash GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getIndexByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getIndexByHash ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getConsensusDataByHeight(ctx, admin, address, 1)
	if err != nil {
		ctx.LogError("TestHeaderApi getConsensusDataByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getConsensusDataByHeight GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi contract invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getConsensusDataByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getConsensusDataByHeight ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getConsensusDataByHash(ctx, admin, address, headerhash)
	if err != nil {
		ctx.LogError("TestHeaderApi getConsensusDataByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getConsensusDataByHash GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi getConsensusDataByHash invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getConsensusDataByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getConsensusDataByHash ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getNextConsensusByHeight(ctx, admin, address, 1)
	if err != nil {
		ctx.LogError("TestHeaderApi getNextConsensusByHeight error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getNextConsensusByHeight GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi getNextConsensusByHeight invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getNextConsensusByHeight return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getNextConsensusByHeight ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = getNextConsensusByHash(ctx, admin, address, headerhash)
	if err != nil {
		ctx.LogError("TestHeaderApi getNextConsensusByHash error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestHeaderApi getNextConsensusByHash GetSmartContractEvent error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestHeaderApi getNextConsensusByHash invoke failed state:0")
		return false
	}
	if len(notifies.Notify) < 1 {
		ctx.LogError("TestHeaderApi getNextConsensusByHash return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestHeaderApi getNextConsensusByHash ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true
}

func getHeaderHashByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, height int) (common.Uint256, error) {

	method := "getHeaderHashByHeight"
	params := make([]interface{}, 1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getHeaderVersionByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, height int) (common.Uint256, error) {

	method := "getHeaderVersionByHeight"
	params := make([]interface{}, 1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getHeaderVersionByHash(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getHeaderVersionByHash"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getHeaderPrevHashByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, height int) (common.Uint256, error) {

	method := "getPrevHashByHeight"
	params := make([]interface{}, 1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getHeaderPrevHashByHash(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getPrevHashByHash"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getMerkelRootByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, height int) (common.Uint256, error) {

	method := "getMerkelRootByHeight"
	params := make([]interface{}, 1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getMerkelRootByHash(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getMerkelRootByHash"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getTimestampByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, height int) (common.Uint256, error) {

	method := "getTimestampByHeight"
	params := make([]interface{}, 1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getTimestampByHash(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getTimestampByHash"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getIndexByHash(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getIndexByHash"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getConsensusDataByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, height int) (common.Uint256, error) {

	method := "getConsensusDataByHeight"
	params := make([]interface{}, 1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getConsensusDataByHash(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getConsensusDataByHash"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getNextConsensusByHeight(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, height int) (common.Uint256, error) {

	method := "getNextConsensusByHeight"
	params := make([]interface{}, 1)
	params[0] = height

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}

func getNextConsensusByHash(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, hash string) (common.Uint256, error) {

	method := "getNextConsensusByHash"
	params := make([]interface{}, 1)
	params[0] = hash

	//txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),address,method,wasm.Json,params,1,false)
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(0, 0, acc, 1, address, method, wasmvm.Raw, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil

}
