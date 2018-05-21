package operator

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestOperationNegative(ctx *testframework.TestFrameworkContext) bool {
	code := "51C56B6C766B00527AC46C766B00C3009C616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationNegative GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.NEOVM,
		false,
		code,
		"TestOperationNegative",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationNegative DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationNegative WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationNegative(ctx, codeAddress, true) {
		return false
	}

	if !testOperationNegative(ctx, codeAddress, false) {
		return false
	}

	return true
}

func testOperationNegative(ctx *testframework.TestFrameworkContext, code common.Address, a bool) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContractWithRes(
		0,
		code,
		[]interface{}{a},
		sdkcom.NEOVM_TYPE_BOOL,
	)
	if err != nil {
		ctx.LogError("TestOperationNegative InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToBoolean(res, !a)
	if err != nil {
		ctx.LogError("TestOperationNegative test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static bool Main(bool a)
    {
        return !a;
    }
}
*/
