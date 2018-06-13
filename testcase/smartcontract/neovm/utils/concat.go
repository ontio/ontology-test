package utils

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

func TestConcat(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C37E616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestConcat GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestConcat",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestConcat DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestConcat WaitForGenerateBlock error:%s", err)
		return false
	}
	input1 := "Hello"
	input2 := "World"
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMContractWithRes(
		codeAddress,
		[]interface{}{input1, input2},
		sdkcom.NEOVM_TYPE_BYTE_ARRAY,
	)
	if err != nil {
		ctx.LogError("TestConcat InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToByteArray(res, []byte(string(input1)+string(input2)))
	if err != nil {
		ctx.LogError("TestConcat test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

public class HelloWorld : SmartContract
{
    public static byte[] Main(byte[] arg1, byte[] arg2)
    {
        return arg1.Concat(arg2);
    }
}
*/
