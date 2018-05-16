package utils

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestTake(ctx *testframework.TestFrameworkContext) bool {
	code := "53c56b6c766b00527ac46c766b51527ac4616c766b00c36c766b51c3806c766b52527ac46203006c766b52c3616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAsString GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
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
		ctx.LogError("TestTake DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
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
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		0,
		0,
		0,
		code,
		[]interface{}{b, count},
		sdkcom.NEOVM_TYPE_BYTE_ARRAY,
	)
	if err != nil {
		ctx.LogError("TestTake InvokeSmartContract error:%s", err)
		return false
	}
	r := count
	if count > len(b) {
		r = len(b)
	}
	err = ctx.AssertToByteArray(res, b[0:r])
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
