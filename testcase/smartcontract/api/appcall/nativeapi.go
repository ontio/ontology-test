package appcall

import (
	"fmt"
	"time"

	"math/big"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

//package main
//import (
//"bytes"
//"math/big"
//
//"encoding/hex"
//"fmt"
//
//"github.com/ontio/ontology/common"
//"github.com/ontio/ontology/smartcontract/service/native/states"
//scstates "github.com/ontio/ontology/smartcontract/states"
//)

//func main() {
//  //from is ontology bookkeeper wallet base58 address
//  //to is any other wallet base58 address

//	from, _ := common.AddressFromBase58("TA51jyVtocZ3aZjCvxCpBABP7u1HGaF7VR")
//	to, _ := common.AddressFromBase58("TA9AFxE1aW6YX4nSSxymr6ghi2HyPwhQ3Z")
//
//	state := &states.State{
//		Version: 1,
//		From:    from,
//		To:      to,
//		Value:   big.NewInt(1234),
//	}
//
//	transfer := &states.Transfers{
//		Version: 1,
//		States:  []*states.State{state},
//	}
//
//	bf := new(bytes.Buffer)
//	transfer.Serialize(bf)
//	args := bf.Bytes()
//
//	ontAddress, _ := common.AddressParseFromBytes([]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
//
//	c := scstates.Contract{
//		Version: 1,
//		Code:    nil,
//		Address: ontAddress,
//		Method:  "transfer",
//		Args:    args,
//	}
//
//	bf = new(bytes.Buffer)
//	err := c.Serialize(bf)
//
//	if err != nil {
//		fmt.Println("contract serialize error")
//		return
//	}
//
//	bytes_contract := bf.Bytes()
//
//	fmt.Println(hex.EncodeToString(bytes_contract))
//}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class A : SmartContract
{
    [Appcall("ff00000000000000000000000000000000000001")]
    public static extern byte[] CallContract();

    public static void Main()
    {
        byte[] result =  CallContract();
        Storage.Put(Storage.CurrentContext, "result", result);
    }
}


Code := 51c56b6161670100ff000000000000000000000000000000000000010b746f74616c537570706c79006c766b00527ac46168164e656f2e53746f726167652e476574436f6e7465787406726573756c746c766b00c3615272680f4e656f2e53746f726167652e50757461616c7566

*/

func TestNativeTotalSupply(ctx *testframework.TestFrameworkContext) bool {
	code := "51c56b6161670100ff000000000000000000000000000000000000010b746f74616c537570706c79006c766b00527ac46168164e656f2e53746f726167652e476574436f6e7465787406726573756c746c766b00c3615272680f4e656f2e53746f726167652e50757461616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()

	fmt.Printf("code Address:%v\n", codeAddress)

	if err != nil {
		ctx.LogError("TestNativeApiContract - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestNativeApiContract DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestNativeApiContract WaitForGenerateBlock error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		0,
		codeAddress,
		[]interface{}{},
	)

	if err != nil {
		ctx.LogError("TestNativeApiContract InvokeSmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestNativeApiContract WaitForGenerateBlock error:%s", err)
		return false
	}

	result, err := ctx.Ont.Rpc.GetStorage(codeAddress, []byte("result"))

	if err != nil {
		ctx.LogError("TestNativeApiContract getstorage error: %s", err)
		return false
	}

	supply := new(big.Int).SetBytes(result)

	err = ctx.AssertToInt(supply, 1000000000)
	if err != nil {
		ctx.LogError(" TestNativeApiContract - AssertToInt error: %s", err)
		return false
	}

	return true
}
