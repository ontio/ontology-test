package utils

import (
	"math/big"
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestRange(ctx *testframework.TestFrameworkContext) bool {
	code := "53C56B6C766B00527AC46C766B51527AC46C766B52527AC46C766B00C36C766B51C36C766B52C37F616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestRange GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
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
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
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
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		code,
		[]interface{}{b, start, count},
		sdkcom.NEOVM_TYPE_BYTE_ARRAY,
	)
	if err != nil {
		ctx.LogError("TestRange InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToByteArray(res, b[start:start+count])
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
