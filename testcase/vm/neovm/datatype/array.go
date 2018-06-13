package datatype

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestArray(ctx *testframework.TestFrameworkContext) bool {
	code := "52c56b6c766b00527ac4616c766b00c3c06c766b51527ac46203006c766b51c3616c7566"
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
		ctx.LogError("TestArray DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestArray WaitForGenerateBlock error:%s", err)
		return false
	}
	params := []interface{}{[]byte("Hello"), []byte("world")}
	if !testArray(ctx, codeAddress, params) {
		return false
	}
	params = []interface{}{[]byte("Hello"), []byte("world"), "123456", 8}
	if !testArray(ctx, codeAddress, params) {
		return false
	}
	return true
}

func testArray(ctx *testframework.TestFrameworkContext, code common.Address, params []interface{}) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMContractWithRes(
		code,
		[]interface{}{params},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestArray InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToInt(res, len(params))
	if err != nil {
		ctx.LogError("TestArray test failed %s", err)
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
    public class Hello : SmartContract
    {
        public static int Main(params object[] args)
        {
            return args.Length;
        }
    }
}
*/
