package operator

import (
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
	"math/big"
	"time"
)

func TestOperationSelfSub(ctx *testframework.TestFrameworkContext) bool {
	code := "51C56B6C766B00527AC46C766B00C35194766A00527AC4616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationSelfAdd GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		false,
		code,
		"TestOperationSelfAdd",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationSelfSub DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationSelfSub WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationSelfSub(ctx, codeAddress, 1) {
		return false
	}

	if !testOperationSelfSub(ctx, codeAddress, -1) {
		return false
	}

	return true
}

func testOperationSelfSub(ctx *testframework.TestFrameworkContext, code common.Address, a int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		code,
		[]interface{}{a},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestOperationSelfSub InvokeSmartContract error:%s", err)
		return false
	}
	a--
	err = ctx.AssertToInt(res, a)
	if err != nil {
		ctx.LogError("TestOperationSelfSub test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main(int a)
    {
        return --a;
    }
}
*/
