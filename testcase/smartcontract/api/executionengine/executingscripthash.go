package executionengine

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class A : SmartContract
{
    public static byte[] Main()
    {
        return ExecutionEngine.ExecutingScriptHash;
    }
}
*/

func TestExecutingScriptHash(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b6161682b53797374656d2e457865637574696f6e456e67696e652e47657443616c6c696e67536372697074486173686c766b00527ac46203006c766b00c3616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestExecutingScriptHash - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"TestExecutingScriptHash",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestExecutingScriptHash DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestExecutingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContractWithRes(
		0,
		codeAddress,
		[]interface{}{},
		sdkcom.NEOVM_TYPE_BYTE_ARRAY,
	)
	if err != nil {
		ctx.LogError("TestExecutingScriptHash error: %s", err)
		return false
	}

	ctx.LogInfo("TestExecutingScriptHash res: %s", res)

	err = ctx.AssertToByteArray(res, codeAddress[:])
	if err != nil {
		ctx.LogError("AssertToByteArray error:%s", err)
		return false
	}

	return true
}
