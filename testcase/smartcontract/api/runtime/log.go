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
    public static void Main(string msg)
    {
        Runtime.Log(msg);
    }
}
*/

func TestRuntimLog(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b6c766b00527ac4616c766b00c361681253797374656d2e52756e74696d652e4c6f6761616c7566"
	codeAddr, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestRuntimLog - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		true,
		code,
		"TestRuntimLog",
		"",
		"",
		"",
		"")

	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 2)

	if err != nil {
		ctx.LogError("TestRuntimLog WaitForGenerateBlock error:%s", err)
		return false
	}

	txHash, err := ctx.Ont.NeoVM.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeAddr,
		[]interface{}{"ontology"})

	if err != nil {
		ctx.LogError("TestRuntimLog InvokeSmartContract error:%s", err)
		return false
	}
	ctx.LogInfo("TestRuntimLog invoke Tx:%s\n", txHash.ToHexString())
	return true
}
