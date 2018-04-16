package runtime

import (
	"math/big"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

/**
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using System;
using System.ComponentModel;
using System.Numerics;

public class HelloWorld : SmartContract
{
    [DisplayName("notify")]
    public static event Action<string, string> Notify;

    public static void Main()
    {
        Notify("hello", "world");
    }
}
*/

func TestRuntimeNotify(ctx *testframework.TestFrameworkContext) bool {
	code := "00c56b61610568656c6c6f05776f726c64617c066e6f7469667953c168124e656f2e52756e74696d652e4e6f7469667961616c7566"
	codeAddr := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestRuntimeNotify - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		code,
		"TestRuntimeNotify",
		"",
		"",
		"",
		"")

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestRuntimeNotify WaitForGenerateBlock error:%s", err)
		return false
	}

	txHash, err := ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer,
		new(big.Int),
		codeAddr,
		[]interface{}{},
	)

	if err != nil {
		ctx.LogError("TestRuntimeNotify InvokeSmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestRuntimeNotify WaitForGenerateBlock error:%s", err)
		return false
	}


	events, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestRuntimeNotify GetSmartContractEvent error:%s", err)
		return false
	}

	notify := events[0].States

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

	return  true
}
