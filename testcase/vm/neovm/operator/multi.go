package operator

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationMulti(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C395616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationMulti GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationMulti",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationMulti DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationMulti WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationMulti(ctx, codeAddress, 23, 2) {
		return false
	}

	if !testOperationMulti(ctx, codeAddress, -1, 34) {
		return false
	}

	if !testOperationMulti(ctx, codeAddress, -1, -2) {
		return false
	}

	if !testOperationMulti(ctx, codeAddress, 0, 100) {
		return false
	}

	return true
}

func testOperationMulti(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMContractWithRes(
		code,
		[]interface{}{a, b},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestOperationMulti InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToInt(res, a*b)
	if err != nil {
		ctx.LogError("TestOperationMulti test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main(int a, int b)
    {
        return a * b;
    }
}
*/
