package utils

import (
	"math/big"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestAsByteArrayBigInteger(ctx *testframework.TestFrameworkContext) bool {
	code := "52c56b6c766b00527ac4616c766b00c36c766b51527ac46203006c766b51c3616c7566"
	codeAddress, _ := utils.GetContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		false,
		code,
		"TestAsByteArrayBigInteger",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger WaitForGenerateBlock error:%s", err)
		return false
	}

	input := new(big.Int).SetInt64(1)
	if !testAsArray_BigInteger(ctx, codeAddress, input) {
		return false
	}
	input = new(big.Int).SetInt64(0)
	if !testAsArray_BigInteger(ctx, codeAddress, input) {
		return false
	}
	input = new(big.Int).SetInt64(-233545554)
	if !testAsArray_BigInteger(ctx, codeAddress, input) {
		return false
	}
	input = new(big.Int).SetInt64(-3434)
	if !testAsArray_BigInteger(ctx, codeAddress, input) {
		return false
	}
	return true
}

func testAsArray_BigInteger(ctx *testframework.TestFrameworkContext, code common.Address, input *big.Int) bool {
	res, err := ctx.Ont.NeoVM.PreExecInvokeNeoVMContract(
		code,
		[]interface{}{input},
	)
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger InvokeSmartContract error:%s", err)
		return false
	}
	if res.State == 0 {
		ctx.LogError("TestAsByteArrayBigInteger PreExecInvokeNeoVMContract failed. state == 0")
		return false
	}
	resValue, err := res.Result.ToInteger()
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger Result.ToInteger error:%s", err)
		return false
	}
	err = ctx.AssertBigInteger(resValue, input)
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger test failed %s", err)
		return false
	}
	return true
}

/*
using System.Numerics;
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static byte[] Main(BigInteger arg)
    {
        return arg.AsByteArray();
    }
}
*/
