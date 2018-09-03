package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationSelfAdd(ctx *testframework.TestFrameworkContext) bool {
	code := "51C56B6C766B00527AC46C766B00C35193766A00527AC4616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationSelfAdd GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationSelfAdd",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationSelfAdd DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationSelfAdd WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationSelfAdd(ctx, codeAddress, 1) {
		return false
	}

	if !testOperationSelfAdd(ctx, codeAddress, -1) {
		return false
	}

	if !testOperationSelfAdd(ctx, codeAddress, 0) {
		return false
	}

	return true
}

func testOperationSelfAdd(ctx *testframework.TestFrameworkContext, code common.Address, a int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{a},
	)
	if err != nil {
		ctx.LogError("TestOperationSelfAdd InvokeSmartContract error:%s", err)
		return false
	}
	resValue,err := res.Result.ToInteger()
	if err != nil {
		ctx.LogError("TestOperationSelfAdd Result.ToInteger error:%s", err)
		return false
	}
	a++
	err = ctx.AssertToInt(resValue, a)
	if err != nil {
		ctx.LogError("TestOperationSelfAdd test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main(int a)
    {
        return ++a;
    }
}
*/
