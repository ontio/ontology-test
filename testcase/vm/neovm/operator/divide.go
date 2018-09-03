package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationDivide(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C396616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationDivide GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationDivide",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationDivide DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationDivide WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationDivideFail(ctx, codeAddress, 10, 0) {
		return false
	}

	if !testOperationDivide(ctx, codeAddress, 23, 2) {
		return false
	}

	if !testOperationDivide(ctx, codeAddress, 544, 345) {
		return false
	}

	if !testOperationDivide(ctx, codeAddress, 3456345, 3545) {
		return false
	}

	if !testOperationDivide(ctx, codeAddress, -10, -234) {
		return false
	}

	if !testOperationDivide(ctx, codeAddress, -345, 34) {
		return false
	}

	if !testOperationDivide(ctx, codeAddress, 0, 100) {
		return false
	}

	return true
}

func testOperationDivide(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{a, b},
	)
	if err != nil {
		ctx.LogError("TestOperationDivide InvokeSmartContract error:%s", err)
		return false
	}
	resValue,err := res.Result.ToInteger()
	if err != nil {
		ctx.LogError("TestOperationDivide Result.ToInteger error:%s", err)
		return false
	}
	err = ctx.AssertToInt(resValue, a/b)
	if err != nil {
		ctx.LogError("TestOperationDivide test %d / %d failed %s", a, b, err)
		return false
	}
	return true
}

func testOperationDivideFail(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	_, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{a, b},
	)
	if err == nil {
		ctx.LogError("testOperationDivideFail %v / %v should failed", a, b)
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
        return a / b;
    }
}
*/
