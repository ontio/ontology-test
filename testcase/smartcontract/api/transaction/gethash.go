package transaction

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

/*

contract A:

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static string Main()
    {
        return "Hello World!";
    }
}

Code:51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566

contract B:

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class A : SmartContract
{
    public static void Main(byte[] txHash)
    {
        Transaction tx = Blockchain.GetTransaction(txHash);
        Storage.Put(Storage.CurrentContext, "txHash", tx.Hash);
    }
}

Code:52c56b6c766b00527ac4616c766b00c361682053797374656d2e426c6f636b636861696e2e4765745472616e73616374696f6e6c766b51527ac461681953797374656d2e53746f726167652e476574436f6e74657874067478486173686c766b51c361681a53797374656d2e5472616e73616374696f6e2e47657448617368615272681253797374656d2e53746f726167652e50757461616c7566
*/

func TestGetTxHash(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566"
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestGetTxHash - GetDefaultAccount error: %s", err)
		return false
	}

	txHash, err := ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		true,
		code,
		"TestGetTxHash",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetTxHash DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxHash WaitForGenerateBlock error: %s", err)
		return false
	}

	code = "52c56b6c766b00527ac4616c766b00c361682053797374656d2e426c6f636b636861696e2e4765745472616e73616374696f6e6c766b51527ac461681953797374656d2e53746f726167652e476574436f6e74657874067478486173686c766b51c361681a53797374656d2e5472616e73616374696f6e2e47657448617368615272681253797374656d2e53746f726167652e50757461616c7566"
	codeAddr, _ := utils.GetContractAddress(code)
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,

		true,
		code,
		"TestGetTxHash",
		"",
		"",
		"",
		"")

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxHash WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.NeoVM.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeAddr,
		[]interface{}{txHash.ToArray()})

	if err != nil {
		ctx.LogError("TestGetTxHash InvokeSmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxHash WaitForGenerateBlock error: %s", err)
		return false
	}

	hash, err := ctx.Ont.GetStorage(codeAddr.ToHexString(), []byte("txHash"))
	if err != nil {
		ctx.LogError("TestGetTxHash - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToByteArray(hash, txHash.ToArray())
	if err != nil {
		ctx.LogError("TestGetTxHash test failed %s", err)
		return false
	}
	return true
}
