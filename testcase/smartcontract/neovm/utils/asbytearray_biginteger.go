package utils

import (
	"math/big"
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestAsByteArrayBigInteger(ctx *testframework.TestFrameworkContext) bool {
	code := "51C56B6C766B00527AC46C766B00C3616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.NEOVM,
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
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger WaitForGenerateBlock error:%s", err)
		return false
	}

	input := new(big.Int).SetInt64(-233545554)
	if !testAsArray_BigInteger(ctx, codeAddress, input) {
		return false
	}
	input = new(big.Int).SetInt64(-3434)
	if !testAsArray_BigInteger(ctx, codeAddress, input) {
		return false
	}
	input = new(big.Int).SetInt64(1)
	if !testAsArray_BigInteger(ctx, codeAddress, input) {
		return false
	}
	return true
}

func testAsArray_BigInteger(ctx *testframework.TestFrameworkContext, code common.Address, input *big.Int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContractWithRes(
		0,
		code,
		[]interface{}{input},
		sdkcom.NEOVM_TYPE_BYTE_ARRAY,
	)
	if err != nil {
		ctx.LogError("TestAsByteArrayBigInteger InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToByteArray(res, utils.ConvertBigIntegerToBytes(input))

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
