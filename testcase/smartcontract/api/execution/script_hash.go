package execution

import (
	"math/big"
	"time"

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
 		byte[] script = ExecutionEngine.ExecutingScriptHash;
        Storage.Put(Storage.CurrentContext, "script", script);
        return null;
    }
}
*/

func TestExecutingScriptHash(ctx *testframework.TestFrameworkContext) bool {
	code := "52c56b6161682d53797374656d2e457865637574696f6e456e67696e652e476574457865637574696e67536372697074486173686c766b00527ac46168164e656f2e53746f726167652e476574436f6e74657874067363726970746c766b00c3615272680f4e656f2e53746f726167652e50757461006c766b51527ac46203006c766b51c3616c7566"
	codeAddr := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestExecutingScriptHash - GetDefaultAccount error: %s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		code,
		"TestContractCreate",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestExecutingScriptHash DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	if err != nil {
		ctx.LogError("TestExecutingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer,
		new(big.Int),
		codeAddr,
		[]interface{}{0})

	if err != nil {
		ctx.LogError("TestExecutingScriptHash error:%s", err)
		return false
	}

	script, err := ctx.Ont.Rpc.GetStorage(codeAddr, []byte("script"))
	if err != nil {
		ctx.LogError("TestExecutingScriptHash - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToByteArray(codeAddr[:], script)
	if err != nil {
		ctx.LogError("TestExecutingScriptHash - AssertToByteArray error:%s", err)
		return false
	}

	return true
}
