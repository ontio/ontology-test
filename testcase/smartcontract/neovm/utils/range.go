package utils

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestRange(ctx *testframework.TestFrameworkContext) bool {
	code := "53C56B6C766B00527AC46C766B51527AC46C766B52527AC46C766B00C36C766B51C36C766B52C37F616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestRange GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestRange",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestRange DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestRange WaitForGenerateBlock error:%s", err)
		return false
	}

	input := []byte("Hello World!")
	if !testRange(ctx, codeAddress, input, 0, len(input)) {
		return false
	}

	if !testRange(ctx, codeAddress, input, 1, len(input)-2) {
		return false
	}

	if !testRange(ctx, codeAddress, input, 2, len(input)-3) {
		return false
	}
	return true
}

func testRange(ctx *testframework.TestFrameworkContext, code common.Address, b []byte, start, count int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{b, start, count},
	)
	if err != nil {
		ctx.LogError("TestRange InvokeSmartContract error:%s", err)
		return false
	}
	resValue, err := res.Result.ToByteArray()
	if err !=nil {
		ctx.LogError("TestRange Result.ToByteArray error:%s", err)
		return false
	}
	err = ctx.AssertToByteArray(resValue, b[start:start+count])
	if err != nil {
		ctx.LogError("TestRange test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static byte[] Main(byte[] arg, int start, int count)
    {
        return arg.Range(start, count);
    }
}
*/
