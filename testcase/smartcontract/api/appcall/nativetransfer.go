package appcall

import (
	"fmt"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

//package main
//
//import (
//"bytes"
//
//"encoding/hex"
//"fmt"
//
//"github.com/ontio/ontology/common"
//"github.com/ontio/ontology/smartcontract/service/native/ont"
//scstates "github.com/ontio/ontology/smartcontract/states"
//)
//
//func main() {
//	//from is ontology bookkeeper wallet base58 address
//	//to is any other wallet base58 address
//
//	from, _ := common.AddressFromBase58("TA51jyVtocZ3aZjCvxCpBABP7u1HGaF7VR")
//	to, _ := common.AddressFromBase58("TA9AFxE1aW6YX4nSSxymr6ghi2HyPwhQ3Z")
//
//	state := &ont.State{
//		Version: 1,
//		From:    from,
//		To:      to,
//		Value:   1234,
//	}
//
//	transfer := &ont.Transfers{
//		Version: 1,
//		States:  []*ont.State{state},
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
    public static extern bool CallContract();

    public static bool Main()
    {
        bool b =  CallContract();
        if (b == true)
        {
            Storage.Put(Storage.CurrentContext, "result", "true");
        }
        else {
            Storage.Put(Storage.CurrentContext, "result", "false");
        }
        return b;
    }
}


Code := 53c56b6161670100ff00000000000000000000000000000000000001087472616e736665722e010101011918a99197a5afef816bdc357fd00c6b4a9a8901e1dea4e3ec2f0bd1d1584b7979b3f973ed65520204d26c766b00527ac46c766b00c36c766b51527ac46c766b51c3644200616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740474727565615272680f4e656f2e53746f726167652e5075746161624000616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740566616c7365615272680f4e656f2e53746f726167652e50757461616c766b00c36c766b52527ac46203006c766b52c3616c7566

*/

func TestCallingContract(ctx *testframework.TestFrameworkContext) bool {
	code := "53c56b6161670100ff00000000000000000000000000000000000001087472616e736665722e010101011918a99197a5afef816bdc357fd00c6b4a9a8901e1dea4e3ec2f0bd1d1584b7979b3f973ed65520204d26c766b00527ac46c766b00c36c766b51527ac46c766b51c3644200616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740474727565615272680f4e656f2e53746f726167652e5075746161624000616168164e656f2e53746f726167652e476574436f6e7465787406726573756c740566616c7365615272680f4e656f2e53746f726167652e50757461616c766b00c36c766b52527ac46203006c766b52c3616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.GetDefaultAccount()

	fmt.Printf("code Address:%v\n", codeAddress)

	if err != nil {
		ctx.LogError("TestCallingContract - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.NEOVM,
		true,
		code,
		"TestCallingContract",
		"",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestCallingContract DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingContract WaitForGenerateBlock error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		0,
		codeAddress,
		[]interface{}{},
	)

	if err != nil {
		ctx.LogError("TestCallingContract InvokeSmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallingContract WaitForGenerateBlock error:%s", err)
		return false
	}

	result, err := ctx.Ont.Rpc.GetStorage(codeAddress, []byte("result"))

	err = ctx.AssertToByteArray(result, []byte("true"))
	if err != nil {
		ctx.LogError(" TestCallingContract - AssertToByteArray error: %s", err)
		return false
	}

	return true
}
