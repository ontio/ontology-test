package cond_loop

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestSwitch(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B00C36C766B51527AC46C766B51C351907C907C9E63080051616C756600616C7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestArray GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		false,
		code,
		"TestArray",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestSwitch DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestSwitch WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testSwitch(ctx, codeAddress, 23) {
		return false
	}

	if !testSwitch(ctx, codeAddress, 1) {
		return false
	}

	if !testSwitch(ctx, codeAddress, 0) {
		return false
	}

	return true
}

func testSwitch(ctx *testframework.TestFrameworkContext, code common.Address, a int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMContractWithRes(
		code,
		[]interface{}{a},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestSwitch InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToInt(res, tswitch(a))
	if err != nil {
		ctx.LogError("TestSwitch test switch %d failed %s", a, err)
		return false
	}
	return true
}

func tswitch(a int) int {
	switch a {
	case 1:
		return 1
	default:
		return 0
	}
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main(int a)
    {
        switch(a)
        {
            case 1:
                return 1;
            default:
                return 0;
        }
    }
}
*/
