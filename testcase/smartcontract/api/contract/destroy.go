package contract
import (
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
	"github.com/ontio/ontology-go-sdk/utils"
	"time"
	"math/big"
)

/*
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

*/

func TestContractDestroy(ctx *testframework.TestFrameworkContext) bool {
	code := "00c56b616168144e656f2e436f6e74726163742e44657374726f7961616c7566"
	codeAddr := utils.GetNeoVMContractAddress(code)

	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestGetContract - GetDefaultAccount error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		code,
		"TestContractDestroy",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestContractDestroy DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 2)
	if err != nil {
		ctx.LogError("TestContractDestroy WaitForGenerateBlock error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer,
		new(big.Int),
		codeAddr,
		[]interface{}{0})


	if err != nil {
		ctx.LogError("TestContractDestroy InvokeSmartContract error: %s", err)
		return false
	}


	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer,
		new(big.Int),
		codeAddr,
		[]interface{}{0})

	if err == nil {
		ctx.LogError("TestContractDestroy InvokeSmartContract error: contract call should be failed.")
		return false
	}

	return true
}

