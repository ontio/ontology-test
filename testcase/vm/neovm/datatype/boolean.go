package datatype

import (
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
	"math/big"
	"time"
)

func TestBoolean(ctx *testframework.TestFrameworkContext)bool{
	code := "00C56B51616C7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestReturnType GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		false,
		code,
		"TestReturnType",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestBoolean DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestBoolean WaitForGenerateBlock error:%s", err)
		return false
	}
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		codeAddress,
		[]interface{}{},
		sdkcom.NEOVM_TYPE_BOOL,
	)
	if err != nil {
		ctx.LogError("TestBoolean InvokeSmartContract error:%s", err)
		return false
	}
	err = ctx.AssertToBoolean(res, true)
	if err != nil {
		ctx.LogError("TestBoolean test failed %s", err)
		return false
	}
	return true
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

class A : SmartContract
{
    public static bool Main()
    {
        return true;
    }
}
*/