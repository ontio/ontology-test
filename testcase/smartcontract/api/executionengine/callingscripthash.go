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
    [Appcall("7d0cd19e13a388af45af797fff87894cecb6d4ae")]
    public static extern byte[] CallContract();
    public static void Main()
    {
        byte[] callScript = CallContract();
        Storage.Put(Storage.CurrentContext, "callScript", callScript);
    }
}
Code := 51c56b616167aed4b6ec4c8987ff7f79af45af88a3139ed10c7d6c766b00527ac461681953797374656d2e53746f726167652e476574436f6e746578740a63616c6c5363726970746c766b00c3615272681253797374656d2e53746f726167652e50757461616c7566
*/

func TestCallingScriptHash(ctx *testframework.TestFrameworkContext) bool {
	codeA := "51c56b6161682b53797374656d2e457865637574696f6e456e67696e652e47657443616c6c696e67536372697074486173686c766b00527ac46203006c766b00c3616c7566"
	codeAddressA, _ := utils.GetContractAddress(codeA)
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestCallingScriptHash - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
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

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	codeB := "51c56b616167aed4b6ec4c8987ff7f79af45af88a3139ed10c7d6c766b00527ac461681953797374656d2e53746f726167652e476574436f6e746578740a63616c6c5363726970746c766b00c3615272681253797374656d2e53746f726167652e50757461616c7566"
	codeAddressB, _ := utils.GetContractAddress(codeB)

	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
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

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	_, err = ctx.Ont.NeoVM.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeAddressB,
		[]interface{}{})

	if err != nil {
		ctx.LogError("TestCallingScriptHash error:%s", err)
		return false
	}

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingScriptHash WaitForGenerateBlock error:%s", err)
		return false
	}

	callScript, err := ctx.Ont.GetStorage(codeAddressB.ToHexString(), []byte("callScript"))
	if err != nil {
		ctx.LogError("TestCallingScriptHash - GetStorage error: %s", err)
		return false
	}

	ctx.LogInfo("CodeA Address:%s", codeAddressA.ToHexString())
	ctx.LogInfo("CodeB Address:%s", codeAddressB.ToHexString())

	err = ctx.AssertToByteArray(callScript, codeAddressB[:])
	if err != nil {
		ctx.LogError("TestCallingScriptHash AssertToByteArray error:%s", err)
		return false
	}
	return true
}
