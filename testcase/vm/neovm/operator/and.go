package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationAnd(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C3640C006C766B51C3616C756600616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationAnd GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationAnd",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationAnd DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationAnd WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationAnd(ctx, codeAddress, true, true) {
		return false
	}

	if !testOperationAnd(ctx, codeAddress, true, false) {
		return false
	}

	if !testOperationAnd(ctx, codeAddress, false, true) {
		return false
	}

	if !testOperationAnd(ctx, codeAddress, false, false) {
		return false
	}

	return true
}

func testOperationAnd(ctx *testframework.TestFrameworkContext, codeAddress common.Address, a, b bool) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		codeAddress,
		[]interface{}{a, b},
	)
	if err != nil {
		ctx.LogError("TestOperationAnd InvokeSmartContract error:%s", err)
		return false
	}
	resValue , err := res.Result.ToBool()
	if err != nil {
		ctx.LogError("TestOperationAnd Result.ToBool error:%s", err)
		return false
	}
	err = ctx.AssertToBoolean(resValue, a && b)
	if err != nil {
		ctx.LogError("TestOperationAnd test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static bool Main(bool a, bool b)
    {
        return a && b;
    }
}
*/
