package wasmvm

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
)

func TestICOContractCollect(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestICOContractCollect wallet.GetDefaultAccount error:%s", err)
		return false
	}
	address, err := GetWasmContractAddress(filePath + "/icotest.wasm")
	if err != nil {
		ctx.LogError("TestICOContract GetWasmContractAddress error:%s", err)
		return false
	}

	//collect
	txHash, err := invokeICOCollect(ctx, admin, address, "TA8Xe297g4wGj67maMYZFmdfk9i2riVNrC", 500)
	if err != nil {
		ctx.LogError("TestICOContract invokeICOCollect error:%s", err)
		return false
	}

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract invokeICOCollect error:%s", err)
		return false
	}

	if len(notifies) < 1 {
		ctx.LogError("TestICOContract invokeICOCollect return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeICOCollect ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = invokeICOBalanceOf(ctx, admin, address, "TA8Xe297g4wGj67maMYZFmdfk9i2riVNrC")
	if err != nil {
		ctx.LogError("TestICOContract invokeBalanceOf error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract init invokeBalanceOf error:%s", err)
		return false
	}
	if len(notifies) < 1 {
		ctx.LogError("TestICOContract invokeICOCollect return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeBalanceOf TA8Xe297g4wGj67maMYZFmdfk9i2riVNrC ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	txHash, err = invokeICOWithdraw(ctx, admin, address, 10)
	if err != nil {
		ctx.LogError("TestICOContract invokeICOWithdraw error:%s", err)
		return false
	}

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestICOContract  invokeICOWithdraw error:%s", err)
		return false
	}

	if len(notifies) < 1 {
		ctx.LogError("TestICOContract invokeICOWithdraw return notifies count error!")
		return false
	}
	ctx.LogInfo("==========TestICOContract invokeICOWithdraw ============")
	for i, n := range notifies {
		ctx.LogInfo(fmt.Sprintf("notify %d is %v", i, n))
	}

	return true
}
