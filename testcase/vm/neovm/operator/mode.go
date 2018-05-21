package operator

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestOperationMode(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C397616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationMode GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.NEOVM,
		false,
		code,
		"TestOperationMode",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationMode DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationMode WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationMode(ctx, codeAddress, 23, 2) {
		return false
	}

	if !testOperationMode(ctx, codeAddress, -345, 34) {
		return false
	}

	if !testOperationMode(ctx, codeAddress, -10, -234) {
		return false
	}

	if !testOperationMode(ctx, codeAddress, 0, 100) {
		return false
	}

	return true
}

func testOperationMode(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContractWithRes(
		0,
		code,
		[]interface{}{a, b},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestOperationMode InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToInt(res, a%b)
	if err != nil {
		ctx.LogError("TestOperationMode test failed %s", err)
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
        return a * b;
    }
}
*/
