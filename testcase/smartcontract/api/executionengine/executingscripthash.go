package executionengine

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
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
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestExecutingScriptHash - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
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

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestExecutingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		codeAddress,
		[]interface{}{},
	)
	if err != nil {
		ctx.LogError("TestExecutingScriptHash error: %s", err)
		return false
	}
	resValue, err := res.Result.ToByteArray()
	if err != nil {
		ctx.LogError("TestExecutingScriptHash Result.ToByteArray error: %s", err)
		return false
	}
	ctx.LogInfo("TestExecutingScriptHash res: %s", resValue)

	err = ctx.AssertToByteArray(resValue, codeAddress[:])
	if err != nil {
		ctx.LogError("AssertToByteArray error:%s", err)
		return false
	}

	return true
}
