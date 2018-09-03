package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationSub(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C394616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationSub GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationSub",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationSub DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationSub WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationSub(ctx, codeAddress, 1, 2) {
		return false
	}

	if !testOperationSub(ctx, codeAddress, -1, 1) {
		return false
	}

	if !testOperationSub(ctx, codeAddress, -1, -2) {
		return false
	}

	if !testOperationSub(ctx, codeAddress, 0, 0) {
		return false
	}

	return true
}

func testOperationSub(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{a, b},
	)
	if err != nil {
		ctx.LogError("TestOperationSub InvokeSmartContract error:%s", err)
		return false
	}
	resValue, err := res.Result.ToInteger()
	if err != nil {
		ctx.LogError("TestOperationSub Result.ToInteger error:%s", err)
		return false
	}
	err = ctx.AssertToInt(resValue, a-b)
	if err != nil {
		ctx.LogError("TestOperationSub test failed %s", err)
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
        return a - b;
    }
}
*/
