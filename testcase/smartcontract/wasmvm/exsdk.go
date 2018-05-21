package wasmvm

import (
	"errors"
	"fmt"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/types"
	"io/ioutil"
	"time"
)

/* move to go sdk
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
	//isPre := false
	//if len(isPreExec) > 0 && isPreExec[0] {
	//	isPre = true
	//}
	return ctx.Ont.Rpc.SendRawTransaction(tx)

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
		return bf.Bytes(),nil
	default:
		return nil,errors.New("unsupported type")
	}
}
*/

func DeployWasmJsonContract(ctx *testframework.TestFrameworkContext, signer *account.Account, wasmfile string, contractName string, version string) (common.Uint256, error) {

	code, err := ioutil.ReadFile(wasmfile)
	if err != nil {
		return common.Uint256{}, errors.New("")
	}

	codeHash := common.ToHexString(code)

	txHash, err := ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		types.WASMVM,
		true,
		codeHash,
		contractName,
		version,
		"test",
		"",
		"",
	)

	if err != nil {
		return common.Uint256{}, fmt.Errorf(" DeploySmartContract error:%s", err)
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}

	//jsonContractAddres = utils.GetContractAddress(codeHash,types.WASMVM)
	return txHash, nil
}

func GetWasmContractAddress(path string) (common.Address, error) {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return common.Address{}, errors.New("")
	}

	codeHash := common.ToHexString(code)
	return utils.GetContractAddress(codeHash, types.WASMVM), nil
}
