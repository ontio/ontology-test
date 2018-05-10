package appcall

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	testwasam "github.com/ontio/ontology-test/testcase/smartcontract/wasmvm"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/wasmvm"
	"github.com/ontio/ontology/smartcontract/types"
)

//
//package main
//
//import (
//"bytes"
//
//"encoding/hex"
//"fmt"
//
//"encoding/json"
//"strconv"
//
//"github.com/ontio/ontology/common"
//scstates "github.com/ontio/ontology/smartcontract/states"
//)
//
//func main() {
//	params := make([]interface{}, 2)
//	params[0] = 20
//	params[1] = 30
//	args, err := BuildWasmContractParam(params)
//
//	if err != nil {
//		fmt.Printf("Build wasam parameter error %s\n", err)
//		return
//	}
//  // 9004e629d5df8306405ce2c9074feb9a2d8d47ef is the wasam contract hash
//	hexBytes, _ := hex.DecodeString("9004e629d5df8306405ce2c9074feb9a2d8d47ef")
//	wasam_addr, _ := common.AddressParseFromBytes(hexBytes)
//	c := scstates.Contract{
//		Version: 1,
//		Code:    nil,
//		Address: wasam_addr,
//		Method:  "add",
//		Args:    args,
//	}
//
//	bf := new(bytes.Buffer)
//	err = c.Serialize(bf)
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
//
//type Args struct {
//	Params []Param `json:"Params"`
//}
//
//type Param struct {
//	Ptype string `json:"type"`
//	Pval  string `json:"value"`
//}
//
//func BuildWasmContractParam(params []interface{}) ([]byte, error) {
//	args := make([]Param, len(params))
//
//	for i, param := range params {
//		switch param.(type) {
//		case string:
//			arg := Param{Ptype: "string", Pval: param.(string)}
//			args[i] = arg
//		case int:
//			arg := Param{Ptype: "int", Pval: strconv.Itoa(param.(int))}
//			args[i] = arg
//		case int64:
//			arg := Param{Ptype: "int64", Pval: strconv.FormatInt(param.(int64), 10)}
//			args[i] = arg
//		case []int:
//			bf := bytes.NewBuffer(nil)
//			array := param.([]int)
//			for i, tmp := range array {
//				bf.WriteString(strconv.Itoa(tmp))
//				if i != len(array)-1 {
//					bf.WriteString(",")
//				}
//			}
//			arg := Param{Ptype: "int_array", Pval: bf.String()}
//			args[i] = arg
//		case []int64:
//			bf := bytes.NewBuffer(nil)
//			array := param.([]int64)
//			for i, tmp := range array {
//				bf.WriteString(strconv.FormatInt(tmp, 10))
//				if i != len(array)-1 {
//					bf.WriteString(",")
//				}
//			}
//			arg := Param{Ptype: "int_array", Pval: bf.String()}
//			args[i] = arg
//		default:
//			return nil, fmt.Errorf("not a supported type :%v\n", param)
//		}
//	}
//
//	bs, err := json.Marshal(Args{args})
//	if err != nil {
//		return nil, err
//	}
//	return bs, nil
//}

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System.Numerics;

public class A : SmartContract
{
    [Appcall("ff00000000000000000000000000000000000001")]
    public static extern string CallContract();

    public static void Main()
    {
        string result =  CallContract();
        Storage.Put(Storage.CurrentContext, "result", result);
    }
}


code = 51c56b61616701009004e629d5df8306405ce2c9074feb9a2d8d47ef03616464447b22506172616d73223a5b7b2274797065223a22696e74222c2276616c7565223a223230227d2c7b2274797065223a22696e74222c2276616c7565223a223330227d5d7d6c766b00527ac46168164e656f2e53746f726167652e476574436f6e7465787406726573756c746c766b00c3615272680f4e656f2e53746f726167652e50757461616c7566
*/

const (
	WASAM_APPCALL_PATH = "testcase/smartcontract/api/appcall"
)

func TestCallWasamContract(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.Wallet.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestCallWasamContract - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = deployWasmJsonContract(ctx, admin)
	if err != nil {
		ctx.LogError("TestCallWasamContract deploy error:%s", err)
		return false
	}

	address, err := testwasam.GetWasmContractAddress(WASAM_APPCALL_PATH + "/contract.wasm")

	if err != nil {
		ctx.LogError("TestCallWasamContract GetWasmContractAddress error:%s", err)
		return false
	}

	_, err = callAdd(ctx, admin, address)
	if err != nil {
		ctx.LogError("TestCallWasamContract callAdd error:%s", err)
		return false
	}

	code := "51c56b61616701009004e629d5df8306405ce2c9074feb9a2d8d47ef03616464447b22506172616d73223a5b7b2274797065223a22696e74222c2276616c7565223a223230227d2c7b2274797065223a22696e74222c2276616c7565223a223330227d5d7d6c766b00527ac46168164e656f2e53746f726167652e476574436f6e7465787406726573756c746c766b00c3615272680f4e656f2e53746f726167652e50757461616c7566"
	codeAddress := utils.GetNeoVMContractAddress(code)
	signer, err := ctx.Wallet.GetDefaultAccount()

	if err != nil {
		ctx.LogError("TestCallWasamContract - GetDefaultAccount error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
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
		ctx.LogError("TestCallWasamContract DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallWasamContract WaitForGenerateBlock error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMSmartContract(
		0,
		0,
		signer,
		0,
		codeAddress,
		[]interface{}{},
	)

	if err != nil {
		ctx.LogError("TestCallWasamContract InvokeSmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestCallWasamContract WaitForGenerateBlock error:%s", err)
		return false
	}

	result, err := ctx.Ont.Rpc.GetStorage(codeAddress, []byte("result"))

	err = ctx.AssertToByteArray(result, []byte("50"))
	if err != nil {
		ctx.LogError(" TestCallWasamContract - AssertToByteArray error: %s", err)
		return false
	}

	return true
}

func callAdd(ctx *testframework.TestFrameworkContext, acc *account.Account, address common.Address) (common.Uint256, error) {
	method := "add"
	params := make([]interface{}, 2)
	params[0] = 20
	params[1] = 30
	txHash, err := ctx.Ont.Rpc.InvokeWasmVMSmartContract(
		0,
		0,
		acc,
		0,
		address,
		method,
		wasmvm.Json,
		params)

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func deployWasmJsonContract(ctx *testframework.TestFrameworkContext, signer *account.Account) (common.Uint256, error) {

	code, err := ioutil.ReadFile(WASAM_APPCALL_PATH + "/" + "contract.wasm")
	if err != nil {
		return common.Uint256{}, errors.New("")
	}

	codeHash := common.ToHexString(code)

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(
		0,
		0,
		signer,
		types.WASMVM,
		true,
		codeHash,
		"",
		"",
		"",
		"",
		"",
	)

	if err != nil {
		return common.Uint256{}, fmt.Errorf("DeploySmartContract error:%s", err)
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}
