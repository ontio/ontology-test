package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationSmallerEqual(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C3A0009C616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationSmallerEqual GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationSmallerEqual",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationSmallerEqual DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationSmallerEqual WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationSmallerEqual(ctx, codeAddress, 23, 2) {
		return false
	}

	if !testOperationSmallerEqual(ctx, codeAddress, -345, 34) {
		return false
	}

	if !testOperationSmallerEqual(ctx, codeAddress, -10, -234) {
		return false
	}

	if !testOperationSmallerEqual(ctx, codeAddress, 100, 100) {
		return false
	}

	return true
}

func testOperationSmallerEqual(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{a, b},
	)
	if err != nil {
		ctx.LogError("TestOperationSmallerEqual InvokeSmartContract error:%s", err)
		return false
	}
	resValue, err := res.Result.ToBool()
	if err != nil {
		ctx.LogError("TestOperationLarger Result.ToBool error:%s", err)
		return false
	}
	err = ctx.AssertToBoolean(resValue, a <= b)
	if err != nil {
		ctx.LogError("TestOperationLarger test %d <= %d failed %s", a, b, err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static bool Main(int a, int b)
    {
        return a <= b;
    }
}
*/
