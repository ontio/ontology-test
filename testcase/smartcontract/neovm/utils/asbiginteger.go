package utils

import (
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/vm/types"
	"math/big"
	"time"
)

func TestAsBigInteger(ctx *testframework.TestFrameworkContext) bool {
	code := "52c56b6c766b00527ac4616c766b00c36c766b51527ac46203006c766b51c3616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestAsBigInteger GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		false,
		code,
		"TestAsBigInteger",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestAsBigInteger DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestAsBigInteger WaitForGenerateBlock error:%s", err)
		return false
	}

	b := big.NewInt(1233)
	if !testAsBigInteger(ctx, codeAddress, b) {
		return false
	}
	b = big.NewInt(0)
	if !testAsBigInteger(ctx, codeAddress, b){
		return false
	}
	b = big.NewInt(-1233)
	if !testAsBigInteger(ctx, codeAddress, b) {
		return false
	}
	return true
}

func testAsBigInteger(ctx *testframework.TestFrameworkContext, code common.Address, b *big.Int) bool {
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		code,
		[]interface{}{b},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestAsBigInteger InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertBigInteger(res, b)
	if err != nil {
		ctx.LogError("TestAsBigInteger test failed %s", err)
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
    public static BigInteger Main(byte[] v)
    {
        return v.AsBigInteger();
    }
}
*/
