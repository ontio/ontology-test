package operator

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestOperationAdd(ctx *testframework.TestFrameworkContext) bool {
	code := "52C56B6C766B00527AC46C766B51527AC46C766B00C36C766B51C393616C7566"
	codeAddress, err := utils.GetContractAddress(code)
	if err != nil {
		ctx.LogError("TestOperationAdd GetContractAddress error:%s", err)
		return false
	}
	ctx.LogInfo("TestOperationAdd contact address:%s", codeAddress.ToHexString())
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOperationAdd GetDefaultAccount error:%s", err)
		return false
	}
	tx, err := ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestOperationAdd",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestOperationAdd DeploySmartContract error:%s", err)
		return false
	}
	ctx.LogInfo("DeployContract TxHash:%s", tx.ToHexString())
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOperationAdd WaitForGenerateBlock error:%s", err)
		return false
	}

	if !testOperationAdd(ctx, codeAddress, 1, 2) {
		return false
	}

	if !testOperationAdd(ctx, codeAddress, -1, 1) {
		return false
	}

	if !testOperationAdd(ctx, codeAddress, -1, -2) {
		return false
	}

	if !testOperationAdd(ctx, codeAddress, 0, 0) {
		return false
	}

	return true
}

func testOperationAdd(ctx *testframework.TestFrameworkContext, codeAddress common.Address, a, b int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		codeAddress,
		[]interface{}{a, b},
	)
	if err != nil {
		ctx.LogError("TestOperationAdd InvokeSmartContract error:%s", err)
		return false
	}
	resValue,err := res.Result.ToInteger()
	if err != nil {
		ctx.LogError("TestOperationAdd Result.ToInteger error:%s", err)
		return false
	}
	err = ctx.AssertToInt(resValue, a+b)
	if err != nil {
		ctx.LogError("TestOperationAdd test failed %s , %d, %d", err, a, b)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static int Main(int a, int b)
    {
        return a + b;
    }
}
*/
