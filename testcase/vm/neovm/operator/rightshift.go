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

func TestOperationRightShift(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C3011F8499616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationRightShift GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		false,
		code,
		"TestOperationRightShift",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationRightShift DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 2)
	if err != nil {
		ctx.LogError("TestOperationRightShift WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, 1, 2) {
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, 34252452, 3) {
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, -1, 2) {
		return false
	}

	if !testOperationRightShift(ctx, codeAddress, 1, -1) {
		return false
	}

	return true
}

func testOperationRightShift(ctx *testframework.TestFrameworkContext, code common.Address, a, b int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		code,
		[]interface{}{a, b},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestOperationRightShift InvokeSmartContract error:%s", err)
		return false
	}
	expect := 0
	if b >= 0 {
		expect = a>>uint(b)
	}
	err = ctx.AssertToInt(res, expect)
	if err != nil {
		ctx.LogError("TestOperationRightShift test %d >> %d failed %s", a, b, err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using System.Numerics;

public class HelloWorld : SmartContract
{
    public static BigInteger Main(BigInteger a, int b)
    {
        return a >> b;
    }
}
*/
