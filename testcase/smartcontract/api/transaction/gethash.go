package transaction

import (
	"math/big"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
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
*/

func TestGetTxHash(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566"
	signer, err := ctx.Wallet.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestGetTxHash - GetDefaultAccount error: %s", err)
		return false
	}

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
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

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxHash WaitForGenerateBlock error: %s", err)
		return false
	}

	code = "52c56b6c766b00527ac4616c766b00c361681d4e656f2e426c6f636b636861696e2e4765745472616e73616374696f6e6c766b51527ac46168164e656f2e53746f726167652e476574436f6e74657874067478486173686c766b51c36168174e656f2e5472616e73616374696f6e2e47657448617368615272680f4e656f2e53746f726167652e50757461616c7566"
	codeAddr := utils.GetNeoVMContractAddress(code)
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		code,
		"TestGetTxHash",
		"",
		"",
		"",
		"")

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxHash WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer,
		new(big.Int),
		codeAddr,
		[]interface{}{txHash.ToArray()})

	if err != nil {
		ctx.LogError("TestGetTxHash InvokeSmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxHash WaitForGenerateBlock error: %s", err)
		return false
	}

	hash, err := ctx.Ont.Rpc.GetStorage(codeAddr, []byte("txHash"))
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
