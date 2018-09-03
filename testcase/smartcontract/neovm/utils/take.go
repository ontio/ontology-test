package utils

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestTake(ctx *testframework.TestFrameworkContext) bool {
	code := "53c56b6c766b00527ac46c766b51527ac4616c766b00c36c766b51c3806c766b52527ac46203006c766b52c3616c7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAsString GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestAsString",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestTake DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestTake WaitForGenerateBlock error:%s", err)
		return false
	}

	input := []byte("Hello World!")
	if !testTake(ctx, codeAddress, input, 0) {
		return false
	}

	if !testTake(ctx, codeAddress, input, len(input)-1) {
		return false
	}

	if !testTake(ctx, codeAddress, input, len(input)) {
		return false
	}
	return true
}

func testTake(ctx *testframework.TestFrameworkContext, code common.Address, b []byte, count int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{b, count},
	)
	if err != nil {
		ctx.LogError("TestTake InvokeSmartContract error:%s", err)
		return false
	}
	r := count
	if count > len(b) {
		r = len(b)
	}
	resValue, err := res.Result.ToByteArray()
	if err != nil {
		ctx.LogError("TestTake Result.ToByteArray error:%s", err)
		return false
	}
	err = ctx.AssertToByteArray(resValue, b[0:r])
	if err != nil {
		ctx.LogError("TestTake test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static byte[] Main(byte[] arg, int count)
    {
        return arg.Take(count);
    }
}
*/
