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
*/

func TestStorage(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b6161681953797374656d2e53746f726167652e476574436f6e74657874026b31027631615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e74657874026b31617c681253797374656d2e53746f726167652e4765746c766b00527ac461681953797374656d2e53746f726167652e476574436f6e74657874026b326c766b00c3615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e74657874026b33027633615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e74657874026b33617c681553797374656d2e53746f726167652e44656c65746561616c7566"
	codeAddr, err := utils.GetContractAddress(code)
	if err != nil {
		ctx.LogError("TestStorage GetContractAddress error:%s", err)
		return false
	}
	ctx.LogInfo("TestStorage address:%s", codeAddr.ToHexString())
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestStorage - GetDefaultAccount error: %s", err)
		return false
	}

	tx, err := ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
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

	ctx.LogInfo("TestStorage deploy tx:%s", tx.ToHexString())
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestStorage WaitForGenerateBlock error: %s", err)
		return false
	}

	invoTx, err := ctx.Ont.NeoVM.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeAddr,
		[]interface{}{})

	if err != nil {
		ctx.LogError("TestStorage InvokeSmartContract error: %s", err)
		return false
	}

	ctx.LogInfo("TestStorage invoke tx:%s\n", invoTx.ToHexString())

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)

	if err != nil {
		ctx.LogError("TestStorage WaitForGenerateBlock error: %s", err)
		return false
	}

	v1, err := ctx.Ont.GetStorage(codeAddr.ToHexString(), []byte("k1"))
	if err != nil {
		ctx.LogError("TestStorage GetStorage k1 error:%s", err)
		return false
	}
	v2, err := ctx.Ont.GetStorage(codeAddr.ToHexString(), []byte("k2"))
	if err != nil {
		ctx.LogError("TestStorage GetStorage k2 error:%s", err)
		return false
	}
	v3, err := ctx.Ont.GetStorage(codeAddr.ToHexString(), []byte("k3"))
	if err != nil {
		ctx.LogError("TestStorage GetStorage k3 error:%s", err)
		return false
	}

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
