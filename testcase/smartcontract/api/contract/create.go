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

class A : SmartContract
{
    public static string Main()
    {
        return "Hello World!";
    }
}

code 51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566
*/

func TestContractCreate(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b610c48656c6c6f20576f726c64216c766b00527ac46203006c766b00c3616c7566"
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
		"TestContractCreate",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetContract DeploySmartContract call code error:%s", err)
		return false
	}

	if err != nil {
		ctx.LogError("TestContractCreate DeploySmartContract error: %s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestContractCreate WaitForGenerateBlock error: %s", err)
		return false
	}


	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(signer,
		new(big.Int),
		codeAddr,
		[]interface{}{0})

	if err != nil {
		ctx.LogError("TestContractCreate InvokeSmartContract error: %s", err)
		return false
	}

	return true
}
