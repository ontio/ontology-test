package operator

import (
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/vm/types"
	"math/big"
	"time"
)

func TestOperationLarger(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C3A0616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationLarger GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		false,
		code,
		"TestOperationLarger",
		"1.0",
		"",
		"",
		"",
	)

	if err != nil {
		ctx.LogError("TestOperationLarger DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationLarger WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationLarger(ctx, codeAddress, 23, 2) {
		return false
	}

	if !testOperationLarger(ctx, codeAddress, -345, 34) {
		return false
	}

	if !testOperationLarger(ctx, codeAddress, -10, -234) {
		return false
	}

	if !testOperationLarger(ctx, codeAddress, 100, 100) {
		return false
	}

	return true
}

func testOperationLarger(ctx *testframework.TestFrameworkContext, codeAddress common.Address, a, b int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		codeAddress,
		[]interface{}{a, b},
		sdkcom.NEOVM_TYPE_BOOL,
	)
	if err != nil {
		ctx.LogError("TestOperationLarger InvokeSmartContract error:%s", err)
		return false
	}

	err = ctx.AssertToBoolean(res, a > b)
	if err != nil {
		ctx.LogError("TestOperationLarger test %d > %d failed %s", a, b, err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static bool Main(int a, int b)
    {
        return a > b;
    }
}
*/
