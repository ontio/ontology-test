package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationRightShift(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C3011F8499616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationRightShift GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestOperationRightShift",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationRightShift DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 2)
	if err != nil {
		ctx.LogError("TestOperationRightShift WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, 1, 2) {
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, 34252452, 3) {
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, -1, 2) {
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, 1, -1) {
		return false
	}

	return true
}

func testOperationRightShift(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{a, b},
	)
	if err != nil {
		ctx.LogError("TestOperationRightShift InvokeSmartContract error:%s", err)
		return false
	}
	resValue,err := res.Result.ToInteger()
	if err != nil {
		ctx.LogError("TestOperationRightShift Result.ToInteger error:%s", err)
		return false
	}
	expect := 0
	if b >= 0 {
		expect = a >> uint(b)
	}
	err = ctx.AssertToInt(resValue, expect)
	if err != nil {
		ctx.LogError("TestOperationRightShift test %d >> %d failed %s", a, b, err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using System.Numerics;

public class HelloWorld : SmartContract
{
    public static BigInteger Main(BigInteger a, int b)
    {
        return a >> b;
    }
}
*/
