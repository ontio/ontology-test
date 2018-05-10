package contract

import (
	"time"

	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"

	"encoding/hex"

	"github.com/ontio/ontology-go-sdk/utils"
)

/**

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static string Main()
    {
        return "Hello World!";
    }
}

code = 51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566
------------------------------------------------------------

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System;
using System.ComponentModel;
using System.Numerics;

class OnTest : SmartContract
{
    public static byte[] Main(byte[] codeHash)
    {
        byte[] script = Blockchain.GetContract(codeHash).Script;
        Storage.Put(Storage.CurrentContext, "script", script);
        return null;
    }
}

code = 53c56b6c766b00527ac4616c766b00c361681a4e656f2e426c6f636b636861696e2e476574436f6e74726163746168164e656f2e436f6e74726163742e4765745363726970746c766b51527ac46168164e656f2e53746f726167652e476574436f6e74657874067363726970746c766b51c3615272680f4e656f2e53746f726167652e50757461006c766b52527ac46203006c766b52c3616c7566
*/

func TestGetContract(ctx *testframework.TestFrameworkContext) bool {
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestGetContract - GetDefaultAccount error: %s", err)
		return false
	}

	codeA := "53c56b6c766b00527ac4616c766b00c361681a4e656f2e426c6f636b636861696e2e476574436f6e74726163746168164e656f2e436f6e74726163742e4765745363726970746c766b51527ac46168164e656f2e53746f726167652e476574436f6e74657874067363726970746c766b51c3615272680f4e656f2e53746f726167652e50757461006c766b52527ac46203006c766b52c3616c7566"
	codeAAddr := utils.GetNeoVMContractAddress(codeA)
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.NEOVM,
		true,
		codeA,
		"TestGetContract",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetContract DeploySmartContract error: %s", err)
		return false
	}

	codeB := "51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566"
	codeBAddr := utils.GetNeoVMContractAddress(codeB)
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.NEOVM,
		true,
		codeB,
		"TestGetContract",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetContract DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetContract - WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(
		0,
		0,
		signer,
		0,
		codeAAddr,
		[]interface{}{codeBAddr[:]})

	if err != nil {
		ctx.LogError("TestGetContract InvokeSmartContract error: %s", err)
		return false
	}

	script, err := ctx.Ont.Rpc.GetStorage(codeAAddr, []byte("script"))
	if err != nil {
		ctx.LogError("TestGetContract - GetStorage error: %s", err)
		return false
	}

	codeBHash, _ := hex.DecodeString(codeB)

	err = ctx.AssertToByteArray(codeBHash, script)
	if err != nil {
		ctx.LogError("TestGetContract - AssertToByteArray error: %s", err)
		return false
	}

	return true
}
