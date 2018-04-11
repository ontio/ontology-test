package hash

import (
	"crypto/sha256"
	"math/big"
	"time"

	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestSha256(ctx *testframework.TestFrameworkContext) bool {
	code := "52c56b6c766b00527ac4616c766b00c3a86c766b51527ac46203006c766b51c3616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestSha256 GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		false,
		code,
		"TestSha256",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("TestSha256 DeploySmartContract error:%s", err)
		return false
	}
	//等待出块
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestSha256 WaitForGenerateBlock error:%s", err)
		return false
	}
	input := []byte("Hello World")
	res, err := ctx.Ont.Rpc.PrepareInvokeNeoVMSmartContract(
		new(big.Int),
		codeAddress,
		[]interface{}{input},
		sdkcom.NEOVM_TYPE_BYTE_ARRAY,
	)
	if err != nil {
		ctx.LogError("TestSha256 InvokeSmartContract error:%s", err)
		return false
	}
	data := csha256(input)
	temp := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		temp[i] = data[i]
	}
	err = ctx.AssertToByteArray(res, temp)
	if err != nil {
		ctx.LogError("TestSha256 test failed %s", err)
		return false
	}
	return true
}

func csha256(input []byte) [32]byte {
	return sha256.Sum256(input)
}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;

public class HelloWorld : SmartContract
{
    public static byte[] Main(byte[] input)
    {
        return Sha256(input);
    }
}
*/
