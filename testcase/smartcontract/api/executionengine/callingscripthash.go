package executionengine

import (
	"time"

	"math/big"

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
        return ExecutionEngine.CallingScriptHash;
    }
}
Code := 51c56b6161682b53797374656d2e457865637574696f6e456e67696e652e47657443616c6c696e67536372697074486173686c766b00527ac46203006c766b00c3616c7566
---------------------------------------------------------------

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class B : SmartContract
{
    [Appcall("7d0cd19e13a388af45af797fff87894cecb6d480")]
    public static extern byte[] CallContract();
    public static void Main()
    {
        byte[] callScript = CallContract();
        Storage.Put(Storage.CurrentContext, "callScript", callScript);
    }
}
Code := 51c56b616167624c902239566a0b7dd4a59dcf38ab57aa36a6706c766b00527ac46168164e656f2e53746f726167652e476574436f6e746578740a63616c6c5363726970746c766b00c3615272680f4e656f2e53746f726167652e50757461616c7566
*/

func TestCallingScriptHash(ctx *testframework.TestFrameworkContext) bool {
	codeA := "51c56b6161682b53797374656d2e457865637574696f6e456e67696e652e47657443616c6c696e67536372697074486173686c766b00527ac46203006c766b00c3616c7566"
	codeAddressA := utils.GetNeoVMContractAddress(codeA)
	signer, err := ctx.Wallet.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestCallingScriptHash - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		codeA,
		"TestCallingScriptHash",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestCallingScriptHash DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	codeB := "51c56b61616780d4b6ec4c8987ff7f79af45af88a3139ed10c7d6c766b00527ac46168164e656f2e53746f726167652e476574436f6e746578740a63616c6c5363726970746c766b00c3615272680f4e656f2e53746f726167652e50757461616c7566"
	codeAddressB := utils.GetNeoVMContractAddress(codeB)

	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		codeB,
		"TestCallingScriptHash",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestCallingScriptHash DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer,
		new(big.Int),
		codeAddressB,
		[]interface{}{0})

	if err != nil {
		ctx.LogError("TestCallingScriptHash error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	callScript, err := ctx.Ont.Rpc.GetStorage(codeAddressB, []byte("callScript"))
	if err != nil {
		ctx.LogError("TestCallingScriptHash - GetStorage error: %s", err)
		return false
	}

	ctx.LogInfo("CodeA Address:%x, R:%x", codeAddressA, utils.BytesReverse(codeAddressA[:]))

	err = ctx.AssertToByteArray(callScript, codeAddressA[:])
	if err != nil {
		ctx.LogError("TestCallingScriptHash AssertToByteArray error:%s", err)
		return false
	}
	return true
}
