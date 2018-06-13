package datatype

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestByteArray(ctx *testframework.TestFrameworkContext) bool {
	code := "53c56b6c766b00527ac46c766b51527ac4616c766b00c36c766b51c39c6c766b52527ac46203006c766b52c3616c7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestByteArray GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestByteArray",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestArray DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestArray WaitForGenerateBlock error:%s", err)
		return false
	}

	arg1 := []byte("Hello")
	arg2 := []byte("World")

	if !testByteArray(ctx, codeAddress, arg1, arg1, true) {
		return false
	}
	if !testByteArray(ctx, codeAddress, arg1, arg2, false) {
		return false
	}
	return true
}

func testByteArray(ctx *testframework.TestFrameworkContext, code common.Address, arg1, arg2 []byte, expect bool) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMContractWithRes(
		code,
		[]interface{}{arg1, arg2},
		sdkcom.NEOVM_TYPE_BOOL,
	)
	if err != nil {
		ctx.LogError("testByteArray InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToBoolean(res, expect)
	if err != nil {
		ctx.LogError("testByteArray test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System;
using System.Numerics;

namespace Hello
{
    public class A : SmartContract
    {
        public static bool Main(byte[] arg1, byte[] arg2)
        {
            return arg1 == arg2;
        }
    }
}
*/
