package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"github.com/ontio/ontology/smartcontract/types"
	"time"
)

func TestCallNeoContract(ctx *testframework.TestFrameworkContext) bool {
	//TODO
	fileName := ""
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestBlockApi wallet.GetDefaultAccount error:%s", err)
		return false
	}

	res := deployNeoContract(ctx)
	if !res {
		ctx.LogError("TestCallNeoContract deployNeoContract failed")
		return false
	}

	txHash, err := DeployWasmJsonContract(ctx, admin, filePath+"/"+fileName, "CNC", "1.0")

	if err != nil {
		ctx.LogError("TestCallNeoContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestCallNeoContract deploy TxHash:%x", txHash)

	address, err := GetWasmContractAddress(filePath + "/" + fileName)
	if err != nil {
		ctx.LogError("TestCallNeoContract GetWasmContractAddress error:%s", err)
		return false
	}

	txHash, err = invokePut(ctx, admin, address, "TestKey", "Hello World")
	if err != nil {
		ctx.LogError("TestCallNeoContract invokePut error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNeoContract invokePut error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallNeoContract contract invoke failed state:0")
		return false
	}
	ctx.LogInfo("==========TestCallNeoContract invokePut ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = invokeGet(ctx, admin, address, "TestKey")
	if err != nil {
		ctx.LogError("TestCallNeoContract invokeGet error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestCallNeoContract invokeGet error:%s", err)
		return false
	}
	if notifies.State == 0 {
		ctx.LogError("TestCallNeoContract contract invoke failed state:0")
		return false
	}
	ctx.LogInfo("==========TestCallNeoContract invokeGet ============")
	for i, n := range notifies.Notify {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true

}

func invokePut(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, key, value string) (common.Uint256, error) {
	method := "putValue"
	params := make([]interface{}, 2)
	params[0] = key
	params[1] = value

	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(), acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
func invokeGet(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address, key string) (common.Uint256, error) {
	method := "getValue"
	params := make([]interface{}, 1)
	params[0] = key

	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(), acc, 1, address, method, wasmvm.Json, params)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func deployNeoContract(ctx *testframework.TestFrameworkContext) bool {
	neocontractCode := "5ac56b6c766b00527ac46c766b51527ac4616c766b00c303507574876c766b52527ac46c766b52c3645d00616c766b51c3c0529c009c6c766b55527ac46c766b55c3640e00006c766b56527ac462a2006c766b51c300c36c766b53527ac46c766b51c351c36c766b54527ac46c766b53c36c766b54c3617c6580006c766b56527ac4626d006c766b00c303476574876c766b57527ac46c766b57c3644900616c766b51c3c0519c009c6c766b59527ac46c766b59c3640e00006c766b56527ac4622f006c766b51c300c36c766b58527ac46c766b58c36165c9006c766b56527ac4620e00006c766b56527ac46203006c766b56c3616c756653c56b6c766b00527ac46c766b51527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c36c766b51c3615272680f4e656f2e53746f726167652e5075746161035075746c766b00c36c766b51c3615272097075745265636f726454c168124e656f2e52756e74696d652e4e6f74696679610350757461680f4e656f2e52756e74696d652e4c6f6761516c766b52527ac46203006c766b52c3616c756652c56b6c766b00527ac46161034765746c766b00c3617c096765745265636f726453c168124e656f2e52756e74696d652e4e6f74696679610347657461680f4e656f2e52756e74696d652e4c6f67616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3617c680f4e656f2e53746f726167652e4765746c766b51527ac46203006c766b51c3616c7566"

	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestDeploySmartContract GetDefaultAccount erro`r:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(), signer,
		types.NEOVM,
		true,
		neocontractCode,
		"TestDeploySmartContract",
		"1.0",
		"",
		"",
		"",
	)

	if err != nil {
		ctx.LogError("TestDeploySmartContract DeploySmartContract error:%s", err)
		return false
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestDeploySmartContract WaitForGenerateBlock error:%s", err)
		return false
	}
	return true
}
