package cond_loop

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestIfElse(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C3A163080051616C75666C766B00C36C766B51C3A26308004F616C756600616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestIfElse GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.NEOVM,
		false,
		code,
		"TestIfElse",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestIfElse DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestIfElse WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testIfElse(ctx, codeAddress, 23, 2) {
		return false
	}

	if !testIfElse(ctx, codeAddress, 2, 23) {
		return false
	}

	if !testIfElse(ctx, codeAddress, 0, 0) {
		return false
	}

	return true
}

func testIfElse(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContractWithRes(
		0,
		code,
		[]interface{}{a, b},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestIfElse InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToInt(res, condIfElse(a, b))
	if err != nil {
		ctx.LogError("TestIfElse test %d ifelse %d failed %s", a, b, err)
		return false
	}
	return true
}

func condIfElse(a, b int) int {
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main(int a, int b)
    {
        if(a > b)
        {
            return 1;
        }
        else if(a < b)
        {
            return -1;
        }
        else{
            return 0;
        }
    }
}
*/
