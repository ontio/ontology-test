package utils

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

func TestAsByteArrayString(ctx *testframework.TestFrameworkContext) bool {
	code := "51C56B6C766B00527AC46C766B00C3616C756600"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAsByteArrayString GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestAsByteArrayString",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestAsByteArrayString DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestAsByteArrayString WaitForGenerateBlock error:%s", err)
		return false
	}

	input := "Hello World"
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		codeAddress,
		[]interface{}{input},
	)
	if err != nil {
		ctx.LogError("TestAsByteArrayString InvokeSmartContract error:%s", err)
		return false
	}
	resValue,err := res.Result.ToByteArray()
	if err != nil {
		ctx.LogError("TestAsByteArrayString Result.ToByteArray error:%s", err)
		return false
	}
	err = ctx.AssertToByteArray(resValue, []byte(input))
	if err != nil {
		ctx.LogError("TestAsByteArrayString test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static byte[] Main(string arg)
    {
        return arg.AsByteArray();
    }
}
*/
