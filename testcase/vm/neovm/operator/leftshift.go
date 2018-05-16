package operator

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestOperationLeftShift(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C3011F8498616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationLeftShift GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.NEOVM,
		false,
		code,
		"TestOperationLeftShift",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationLeftShift DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationLeftShift WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationLeftShift(ctx, codeAddress, 1, 2) {
		return false
	}

	if !testOperationLeftShift(ctx, codeAddress, 1326567565434, 2) {
		return false
	}

	if !testOperationLeftShift(ctx, codeAddress, 2, 3) {
		return false
	}

	if !testOperationLeftShift(ctx, codeAddress, -1, 2) {
		return false
	}

	if !testOperationLeftShift(ctx, codeAddress, 1, -1) {
		return false
	}

	return true
}

func testOperationLeftShift(ctx *testframework.TestFrameworkContext, code common.Address, a int, b int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		0,
		0,
		0,
		code,
		[]interface{}{a, b},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestOperationLeftShift InvokeSmartContract error:%s", err)
		return false
	}
	expect := 0
	if b >= 0 {
		expect = a << uint(b)
	}
	err = ctx.AssertToInt(res, expect)
	if err != nil {
		ctx.LogError("TestOperationLeftShift test %d << %d failed %s", a, b, err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main(int a, int b)
    {
        return a << b;
    }
}
*/
