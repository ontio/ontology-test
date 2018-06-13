package datatype

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

func TestString(ctx *testframework.TestFrameworkContext) bool {
	code := "00C56B0B48656C6C6F20576F726C64616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestString GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestString",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestString DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestString WaitForGenerateBlock error:%s", err)
		return false
	}
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMContractWithRes(
		codeAddress,
		[]interface{}{},
		sdkcom.NEOVM_TYPE_STRING,
	)
	if err != nil {
		ctx.LogError("TestString InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToString(res, "Hello World")
	if err != nil {
		ctx.LogError("TestString test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static string Main()
    {
        return "Hello World";
    }
}
*/
