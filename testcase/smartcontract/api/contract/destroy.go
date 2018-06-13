package contract

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

/*
contract A

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;

public class Contract1:SmartContract
{
    public static void Main()
    {
        Neo.SmartContract.Framework.Services.Neo.Contract.Destroy();
    }
}

code = 00c56b616168144e656f2e436f6e74726163742e44657374726f7961616c7566

------------------------------------------------------------------------
contract B

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System;
using System.ComponentModel;
using System.Numerics;

class OnTest : SmartContract
{
    public static bool Main(byte[] codeHash)
    {
        byte[] script = Blockchain.GetContract(codeHash).Script;
        if (script == null || script.Length == 0)
        {
            return false;
        }
        return true;
    }
}

code = 54c56b6c766b00527ac4616c766b00c361681a4e656f2e426c6f636b636861696e2e476574436f6e74726163746168164e656f2e436f6e74726163742e4765745363726970746c766b51527ac46c766b51c3640e006c766b51c3c0009c620400516c766b52527ac46c766b52c3640f0061006c766b53527ac4620e00516c766b53527ac46203006c766b53c3616c7566
*/

func TestContractDestroy(ctx *testframework.TestFrameworkContext) bool {
	code := "00c56b616168144e656f2e436f6e74726163742e44657374726f7961616c7566"
	codeAddressA := utils.GetNeoVMContractAddress(code)

	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestGetContract - GetDefaultAccount error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"TestContractDestroy",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestContractDestroy DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 2)
	if err != nil {
		ctx.LogError("TestContractDestroy WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		0,
		codeAddressA,
		[]interface{}{0})

	if err != nil {
		ctx.LogError("TestContractDestroy InvokeSmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 2)
	if err != nil {
		ctx.LogError("TestContractDestroy WaitForGenerateBlock error: %s", err)
		return false
	}

	code = "54c56b6c766b00527ac4616c766b00c361681a4e656f2e426c6f636b636861696e2e476574436f6e74726163746168164e656f2e436f6e74726163742e4765745363726970746c766b51527ac46c766b51c3640e006c766b51c3c0009c620400516c766b52527ac46c766b52c3640f0061006c766b53527ac4620e00516c766b53527ac46203006c766b53c3616c7566"
	codeAddressB := utils.GetNeoVMContractAddress(code)

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"TestContractDestroy",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestContractDestroy DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 2)
	if err != nil {
		ctx.LogError("TestContractDestroy WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		0,
		codeAddressB,
		[]interface{}{codeAddressA[:]})

	if err != nil {
		ctx.LogError("TestContractDestroy contract should be destroyed。")
		return false
	}

	return true
}
