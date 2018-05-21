package call

import (
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestCallContractStatic(ctx *testframework.TestFrameworkContext) bool {
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestIfElse GetDefaultAccount error:%s", err)
		return false
	}

	codeA := "52c56b6c766b00527ac4616c766b00c36c766b51527ac46203006c766b51c3616c7566"
	codeAddressA := utils.GetNeoVMContractAddress(codeA)

	//Because of compiler will reverse of the address, so the we need to reverse the address of called contract.
	//After fix of compiler, wo won't need reverse.
	ctx.LogInfo("CodeA Address:%x, R:%x", codeAddressA, utils.BytesReverse(codeAddressA[:]))

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		false,
		codeA,
		"TestCallContractStaticA",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestCallContractStatic DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallContractStatic WaitForGenerateBlock error:%s", err)
		return false
	}

	codeB := "52c56b6c766b00527ac4616c766b00c3616780711163a4da8a8e37fd469a37e6cc04d37df3696c766b51527ac46203006c766b51c3616c7566"
	codeAddressB := utils.GetNeoVMContractAddress(codeB)
	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		false,
		codeB,
		"TestCallContractStaticB",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestCallContractStatic DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallContractStatic WaitForGenerateBlock error:%s", err)
		return false
	}

	input := 12
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContractWithRes(
		0,
		codeAddressB,
		[]interface{}{input},
		sdkcom.NEOVM_TYPE_INTEGER,
	)
	if err != nil {
		ctx.LogError("TestCallContractStatic error:%s", err)
		return false
	}
	err = ctx.AssertToInt(res, input)
	if err != nil {
		ctx.LogError("TestCallContractStatic res AssertToInt error:%s", err)
		return false
	}
	return true
}

/*
SmartContractA

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using System.Numerics;

public class A : SmartContract
{
    public static int Main(int arg)
    {
        return arg;
    }
}

Code:52c56b6c766b00527ac4616c766b00c36c766b51527ac46203006c766b51c3616c7566

SmartContractB

using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
usingtestcase/smartcontract/native/transfer.go System.Numerics;

public class B : SmartContract
{
	//Because of compiler will reverse of the address, so the we need to reverse the address of called contract.
	//After fix of compiler, wo won't need reverse.
    [Appcall("69f37dd304cce6379a46fd378e8adaa463117180")]
    public static extern int OtherContract(int input);
    public static int Main(int input)
    {
        return OtherContract(input);
    }
}

Code:52c56b6c766b00527ac4616c766b00c3616780711163a4da8a8e37fd469a37e6cc04d37df3696c766b51527ac46203006c766b51c3616c7566
*/
