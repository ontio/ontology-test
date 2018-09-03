package runtime

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

/**
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using System;
using System.ComponentModel;
using System.Numerics;

public class HelloWorld : SmartContract
{
    public delegate void PushDelegate(string operation,string msg);

     [DisplayName("notify")]
     public static event PushDelegate Notify;

    public static void Main()
    {
        Notify("hello", "world");
    }
}
*/

func TestRuntimeNotify(ctx *testframework.TestFrameworkContext) bool {
	code := "00c56b61610568656c6c6f05776f726c64617c066e6f7469667953c1681553797374656d2e52756e74696d652e4e6f7469667961616c7566"
	codeAddr, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestRuntimeNotify - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		true,
		code,
		"TestRuntimeNotify",
		"",
		"",
		"",
		"")

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestRuntimeNotify WaitForGenerateBlock error:%s", err)
		return false
	}

	txHash, err := ctx.Ont.NeoVM.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeAddr,
		[]interface{}{},
	)

	if err != nil {
		ctx.LogError("TestRuntimeNotify InvokeSmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestRuntimeNotify WaitForGenerateBlock error:%s", err)
		return false
	}

	events, err := ctx.Ont.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		ctx.LogError("TestRuntimeNotify GetSmartContractEvent error:%s", err)
		return false
	}
	if events.State == 0 {
		ctx.LogError("TestRuntimeNotify contract invoke failed, state:0")
		return false
	}
	notify := events.Notify[0].States.([]interface{})

	name, _ := ctx.ConvertToHexString(notify[0])

	err = ctx.AssertToString(name, "notify")
	if err != nil {
		ctx.LogError("TestRuntimeNotify failed AssertToString:%s", err)
		return false
	}

	key, _ := ctx.ConvertToHexString(notify[1])
	err = ctx.AssertToString(key, "hello")
	if err != nil {
		ctx.LogError("TestRuntimeNotify failed AssertToString %s ", err)
		return false
	}

	value, _ := ctx.ConvertToHexString(notify[2])
	err = ctx.AssertToString(value, "world")
	if err != nil {
		ctx.LogError("TestRuntimeNotify failed AssertToString %s ", err)
		return false
	}

	return true
}
