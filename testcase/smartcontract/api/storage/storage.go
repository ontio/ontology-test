package storage

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

/*

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static void Main()
    {
        Storage.Put(Storage.CurrentContext, "k1", "v1");
        byte[] temp = Storage.Get(Storage.CurrentContext, "k1");
        Storage.Put(Storage.CurrentContext, "k2", temp);
        Storage.Put(Storage.CurrentContext, "k3", "v3");
        Storage.Delete(Storage.CurrentContext, "k3");
    }
}

code 51c56b616168164e656f2e53746f726167652e476574436f6e74657874026b31027631615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e74657874026b31617c680f4e656f2e53746f726167652e4765746c766b00527ac46168164e656f2e53746f726167652e476574436f6e74657874026b326c766b00c3615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e74657874026b33027633615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e74657874026b33617c68124e656f2e53746f726167652e44656c65746561616c7566
*/

func TestStorage(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b616168164e656f2e53746f726167652e476574436f6e74657874026b31027631615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e74657874026b31617c680f4e656f2e53746f726167652e4765746c766b00527ac46168164e656f2e53746f726167652e476574436f6e74657874026b326c766b00c3615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e74657874026b33027633615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652e476574436f6e74657874026b33617c68124e656f2e53746f726167652e44656c65746561616c7566"
	codeAddr, _ := utils.GetContractAddress(code)

	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestStorage - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		true,
		code,
		"TestStorage",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestStorage DeploySmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestStorage WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeAddr,
		[]interface{}{0})

	if err != nil {
		ctx.LogError("TestStorage InvokeSmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	if err != nil {
		ctx.LogError("TestStorage WaitForGenerateBlock error: %s", err)
		return false
	}

	v1, _ := ctx.Ont.Rpc.GetStorage(codeAddr, []byte("k1"))
	v2, _ := ctx.Ont.Rpc.GetStorage(codeAddr, []byte("k2"))
	v3, _ := ctx.Ont.Rpc.GetStorage(codeAddr, []byte("k3"))

	err = ctx.AssertToByteArray(v1, []byte("v1"))
	if err != nil {
		ctx.LogError("TestStorage - AssertToByteArray error: %s", err)
		return false
	}

	err = ctx.AssertToByteArray(v2, []byte("v1"))
	if err != nil {
		ctx.LogError("TestStorage - AssertToByteArray error: %s", err)
		return false
	}

	err = ctx.AssertToByteArray(v3, []byte(""))
	if err != nil {
		ctx.LogError("TestStorage - AssertToByteArray error: %s", err)
		return false
	}

	return true
}
