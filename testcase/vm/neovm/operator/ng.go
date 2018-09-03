package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationNegative(ctx *testframework.TestFrameworkContext) bool {
	code := "51C56B6C766B00527AC46C766B00C3009C616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationNegative GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationNegative",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationNegative DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationNegative WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationNegative(ctx, codeAddress, true) {
		return false
	}

	if !testOperationNegative(ctx, codeAddress, false) {
		return false
	}

	return true
}

func testOperationNegative(ctx *testframework.TestFrameworkContext, code common.Address, a bool) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{a},
	)
	if err != nil {
		ctx.LogError("TestOperationNegative InvokeSmartContract error:%s", err)
		return false
	}
	resValue, err := res.Result.ToBool()
	if err != nil {
		ctx.LogError("TestOperationNegative Result.ToBool error:%s", err)
		return false
	}
	err = ctx.AssertToBoolean(resValue, !a)
	if err != nil {
		ctx.LogError("TestOperationNegative test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static bool Main(bool a)
    {
        return !a;
    }
}
*/
