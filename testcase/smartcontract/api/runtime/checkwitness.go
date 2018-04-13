package runtime

import (
	"math/big"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class A : SmartContract
{
    public static bool Main(byte[] input)
    {
        return Runtime.CheckWitness(input);
		Storage.Put(Storage.CurrentContext, "isowner", );
        return null;
    }
}
*/

func TestCheckWitness(ctx *testframework.TestFrameworkContext) bool {
	code := "52c56b6c766b00527ac4616c766b00c36168184e656f2e52756e74696d652e436865636b5769746e6573736c766b51527ac46203006c766b51c3616c7566"
	codeAddr := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestCheckWitness - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		code,
		"TestCheckWitness",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestCheckWitness DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	if err != nil {
		ctx.LogError("TestCheckWitness WaitForGenerateBlock error:%s", err)
		return false
	}

	if !checkWitness(ctx, code, ctx.OntClient.Account1.ProgramHash.ToArray(), true) {
		return false
	}

	if !checkWitness(ctx, code, ctx.OntClient.Account2.ProgramHash.ToArray(), false) {
		return false
	}

	return true
}

//func checkWitness(ctx *testframework.TestFrameworkContext, code string, input []byte, expect bool) bool {
func checkWitness(ctx *testframework.TestFrameworkContext, signer account.Account, codeAddress common.Address, input []byte, expect bool) bool {
	//res, err := ctx.Ont.Rpc.InvokeSmartContract(
	//	ctx.OntClient.Account1,
	//	code,
	//	[]interface{}{input},
	//)
	_, err := ctx.Ont.Rpc.InvokeNeoVMSmartContract(&signer,
		new(big.Int),
		codeAddress,
		[]interface{}{input})

	if err != nil {
		ctx.LogError("CheckWitness error:%s", err)
		return false
	}
	err = ctx.AssertToBoolean(res, expect)
	if err != nil {
		ctx.LogError("CheckWitness AssertToBoolean error:%s", err)
		return false
	}
	return true
}
