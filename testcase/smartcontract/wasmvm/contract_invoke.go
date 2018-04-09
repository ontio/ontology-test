package wasmvm

import (
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/vm/types"
	"io/ioutil"
	"errors"
	"fmt"
	"time"
	"github.com/ontio/ontology-go-sdk/utils"
	"math/big"
	"github.com/ontio/ontology/smartcontract/service/native/states"
	"github.com/ontio/ontology/smartcontract/service/wasm"
	"github.com/ontio/ontology/vm/wasmvm/exec"
	"strconv"
	"bytes"
	"encoding/json"
	"github.com/ontio/ontology/common/serialization"
	"encoding/binary"
)

const (
	 filePath = "/home/zhoupw/work/go/src/github.com/ontio/ontology/vm/wasmvm/exec/test_data2"
)

var jsonContractAddres common.Address

func TestWasmJsonContract(ctx *testframework.TestFrameworkContext) bool{
	nep5Wallet := "/home/zhoupw/work/go/src/github.com/ontio/ontology/wallet.dat"
	nep5WalletPwd := "123456"
	wallet, err := ctx.Ont.OpenWallet(nep5Wallet, nep5WalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", nep5Wallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestWasmJsonContract wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := deployWasmJsonContract(ctx, admin)
	if err != nil {
		ctx.LogError("TestWasmJsonContract deploy error:%s", err)
		return false
	}

	ctx.LogInfo("TestWasmJsonContract deploy TxHash:%x", txHash)

	txHash,err = invokeContract(ctx,admin)
	if err != nil {
		ctx.LogError("TestWasmJsonContract invokeContract error:%s", err)
		return false
	}
	ctx.LogInfo("invokeContract: %x\n", txHash)
	ctx.LogInfo("TestWasmJsonContract invokeContract success")
	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestWasmJsonContract init GetSmartContractEvent error:%s", err)
		return false
	}
	ctx.LogInfo("TestNep5Contract init notify %s", notifies)
	fmt.Println("============result is===============")
	bs ,_:= common.HexToBytes(notifies[0].States[0].(string))

	fmt.Printf("+==========%s\n",string(bs))
	return true
}


func invokeContract(ctx *testframework.TestFrameworkContext, acc *account.Account) (common.Uint256, error) {
	method := "add"
	params := make([]interface{},2)
	params[0] = 20
	params[1] = 30
	txHash,err := InvokeWasmVMContract(ctx,acc,new(big.Int),jsonContractAddres,method,wasm.Json,params,1,false)
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func deployWasmJsonContract(ctx *testframework.TestFrameworkContext, signer *account.Account) (common.Uint256, error){

	code, err := ioutil.ReadFile(filePath + "/" + "contract.wasm")
	if err != nil {
		return common.Uint256{}, errors.New("")
	}

	codeHash := common.ToHexString(code)

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(
		signer,
		types.WASMVM,
		true,
		codeHash,
		"wjc",
		"1.0",
		"test",
		"",
		"",
	)

	if err != nil {
		return common.Uint256{}, fmt.Errorf("TestNep5Contract DeploySmartContract error:%s", err)
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}

	jsonContractAddres = utils.GetContractAddress(codeHash,types.WASMVM)
	return txHash, nil
}

func InvokeWasmVMContract(ctx *testframework.TestFrameworkContext,
			siger*account.Account,
			gasLimit *big.Int,
			smartcodeAddress common.Address,
			methodName string,
			paramType wasm.ParamType,
			params []interface{},
			ver byte,
			isPreExec ...bool)(common.Uint256, error) {

	contract := &states.Contract{}
	contract.Address = smartcodeAddress
	contract.Method = methodName
	contract.Version = ver

	argbytes,err := buildWasmContractParam(params,paramType)

	if err != nil {
		return common.UINT256_EMPTY,err
	}

	contract.Args = argbytes

	bf := bytes.NewBuffer(nil)
	contract.Serialize(bf)

	tx :=  ctx.Ont.Rpc.NewInvokeTransaction( new(big.Int),types.WASMVM,bf.Bytes())
	err = ctx.Ont.Rpc.SignTransaction(tx, siger)
	if err != nil {
		return common.Uint256{}, nil
	}
	isPre := false
	if len(isPreExec) > 0 && isPreExec[0] {
		isPre = true
	}
	return ctx.Ont.Rpc.SendRawTransaction(tx, isPre)

}


func buildWasmContractParam(params []interface{},paramType wasm.ParamType)([]byte,error){
	switch paramType {
	case wasm.Json:
		args:= make([]exec.Param,len(params))

		for i, param := range params {
			switch param.(type){
			case string:
				arg := exec.Param{Ptype:"string",Pval:param.(string)}
				args[i] = arg
			case int:
				arg := exec.Param{Ptype:"int",Pval:strconv.Itoa(param.(int))}
				args[i] = arg
			case int64:
				arg := exec.Param{Ptype:"int64",Pval:strconv.FormatInt(param.(int64),10)  }
				args[i] = arg
			case []int:
				bf := bytes.NewBuffer(nil)
				array := param.([]int)
				for i,tmp := range array {
					bf.WriteString(strconv.Itoa(tmp))
					if i != len(array) - 1{
						bf.WriteString(",")
					}
				}
				arg := exec.Param{Ptype:"int_array",Pval:bf.String() }
				args[i] = arg
			case []int64:
				bf := bytes.NewBuffer(nil)
				array := param.([]int64)
				for i,tmp := range array {
					bf.WriteString(strconv.FormatInt(tmp,10) )
					if i != len(array) - 1{
						bf.WriteString(",")
					}
				}
				arg := exec.Param{Ptype:"int_array",Pval:bf.String() }
				args[i] = arg
			default:
				return nil,errors.New(fmt.Sprintf("not a supported type :%v\n",param))
			}
		}

		bs,err := json.Marshal(exec.Args{args})
		if err != nil {
			return nil,err
		}
		return bs,nil
	case wasm.Raw:
		bf := bytes.NewBuffer(nil)
		for _, param := range params {
			switch param.(type){
			case string :
				tmp := bytes.NewBuffer(nil)
				serialization.WriteVarString(tmp, param.(string))
				bf.Write(tmp.Bytes())

			case int :
				tmpBytes := make([]byte, 4)
				binary.LittleEndian.PutUint32(tmpBytes, uint32(param.(int)))
				bf.Write(tmpBytes)

			case int64:
				tmpBytes := make([]byte, 8)
				binary.LittleEndian.PutUint64(tmpBytes, uint64(param.(int64)))
				bf.Write(tmpBytes)

			default:
				return nil,errors.New(fmt.Sprintf("not a supported type :%v\n",param))
			}
		}
	default:
		return nil,errors.New("unsupported type")
	}
	return nil,errors.New("unsupported type")
}


