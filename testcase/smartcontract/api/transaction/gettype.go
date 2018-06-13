package transaction

import (
	"time"

	"bytes"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	ctypes "github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/types"
)

/*

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class A : SmartContract
{
    public static BigInteger Main(byte[] txHash)
    {
        Transaction tx = Blockchain.GetTransaction(txHash);
        Storage.Put(Storage.CurrentContext, "txType", tx.Type);
        return 0;
    }
}

*/

func TestGetTxType(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566"
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestGetTxType - GetDefaultAccount error: %s", err)
		return false
	}

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"TestGetTxType",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetTxType DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"TestGetTxType",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetTxType DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxType WaitForGenerateBlock error: %s", err)
		return false
	}

	code = "53c56b6c766b00527ac4616c766b00c361681d4e656f2e426c6f636b636861696e2e4765745472616e73616374696f6e6c766b51527ac46168164e656f2e53746f726167652e476574436f6e74657874067478547970656c766b51c36168174e656f2e5472616e73616374696f6e2e47657454797065615272680f4e656f2e53746f726167652e50757461006c766b52527ac46203006c766b52c3616c7566"
	codeAddr := utils.GetNeoVMContractAddress(code)
	txHash, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"TestGetTxType",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetTxType DeploySmartContract error: %s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetTxType WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		0,
		codeAddr,
		[]interface{}{txHash.ToArray()})

	if err != nil {
		ctx.LogError("TestGetTxType InvokeSmartContract error: %s", err)
		return false
	}

	txType, err := ctx.Ont.Rpc.GetStorage(codeAddr, []byte("txType"))
	if err != nil {
		ctx.LogError("TestGetTxType - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToInt(txType, int(0xd0))
	if bytes.Equal(txType, []byte{byte(ctypes.Deploy)}) {
		ctx.LogError("TestGetTxType AssertToInt %s", err)
	}
	return true
}
