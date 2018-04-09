package utils

import (
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/vm/types"
	"math/big"
	"time"
)

func TestAsString(ctx *testframework.TestFrameworkContext) bool {
	code := "51C56B6C766B00527AC46C766B00C3616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAsString GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		false,
		code,
		"TestAsString",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestAsString DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestAsString WaitForGenerateBlock error:%s", err)
		return false
	}
	input := []byte("Hello World")
	if !testAsString(ctx, codeAddress, input) {
		return false
	}
	//input = []byte("")
	//if !testAsString(ctx, codeAddress, input) {
	//	return false
	//}
	return true
}

func testAsString(ctx *testframework.TestFrameworkContext, code common.Address, input []byte) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		code,
		[]interface{}{input},
		sdkcom.NEOVM_TYPE_STRING,
	)
	if err != nil {
		ctx.LogError("TestAsString InvokeSmartContract error:%s", err)
		return false
	}

	err = ctx.AssertToString(res, string(input))
	if err != nil {
		ctx.LogError("TestAsString test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

public class HelloWorld : SmartContract
{
    public static string Main(byte[] input)
    {
        return input.AsString();
    }
}
*/
