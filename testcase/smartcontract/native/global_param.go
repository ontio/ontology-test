package native

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/native/global_params"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

func TestGlobalParam(ctx *testframework.TestFrameworkContext) bool {
	originAdmin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("Wallet.GetDefaultAccount error:%s", err)
		return false
	}
	initKey := "gasPrice"
	// query init parmas
	initParms, err := getParam(&global_params.ParamNameList{initKey}, ctx.Ont)
	if err != nil {
		ctx.LogError("Get global params error, the value should initialize in genesis!")
		return false
	}
	// add a global params
	initParms.SetParam(&global_params.Param{"init-key", "init-value"})
	// origin admin is origin operator
	originOperator := originAdmin
	setParam(originOperator, ctx.Ont, initParms)
	// query init param, should cause error, because not create snapshot
	_, err = getParam(&global_params.ParamNameList{"init-key"}, ctx.Ont)
	if err == nil {
		ctx.LogError("Get global params error, the value should not take effect!")
		return false
	}
	// create snapshot of init params
	createSnapshot(originOperator, ctx.Ont)
	// query init param, should not cause error
	initParms, err = getParam(&global_params.ParamNameList{initKey, "init-key"}, ctx.Ont)
	if err != nil {
		ctx.LogError("Get global params error: %s!", err)
		return false
	}
	// origin admin set and get global params, should not cause error
	globalParams, err := testGetAndSet(originOperator, ctx.Ont, "originOperator0", *initParms)
	if err != nil {
		ctx.LogError("Origin operator operate global params error:%s", err)
		return false
	}
	newOperator, err := ctx.NewAccount()
	if err != nil {
		ctx.LogError("Wallet.NewAccount error:%s", err)
		return false
	}
	// new operator set and get global params, should cause error, because  he is not operator
	globalParams, err = testGetAndSet(newOperator, ctx.Ont, "newOperator0", globalParams)
	if err == nil {
		ctx.LogError("New operator operate global params error, should not be authorized!")
		return false
	}
	// set newOperator as operator
	setOperator(originAdmin, newOperator.Address, ctx.Ont)
	// newOperator set and get global params, should not cause error, because he is operator
	globalParams, err = testGetAndSet(newOperator, ctx.Ont, "newOperator0", globalParams)
	if err != nil {
		ctx.LogError("New operator operate global params error: ", err)
		return false
	}
	// originOperator set and get global param, should cause error, because he is not operator
	globalParams, err = testGetAndSet(originOperator, ctx.Ont, "originOperator1", globalParams)
	if err == nil {
		ctx.LogError("Origin operator operate global params error:, should not be authorized!")
		return false
	}
	// transfer admin to newOperator
	newAdmin := newOperator
	transferAdmin(originAdmin, newAdmin.Address, ctx.Ont)
	// new admin set originOperator as current operator, should not success, because new admin does not accept
	setOperator(newAdmin, originOperator.Address, ctx.Ont)
	// originOperator operate global params, should cause error, because newAdmin is not admin
	// and he's setOperator is not effective
	globalParams, err = testGetAndSet(originOperator, ctx.Ont, "originOperator1", globalParams)
	if err == nil {
		ctx.LogError("Origin operator operate global params error, should not be authorized!")
		return false
	}
	// new admin accept permission
	acceptAdmin(newAdmin, ctx.Ont)
	// new admin set originOperator as current operator
	setOperator(newAdmin, originOperator.Address, ctx.Ont)
	// originOperator operate global params, should not cause error
	globalParams, err = testGetAndSet(originOperator, ctx.Ont, "originOperator1", globalParams)
	if err != nil {
		ctx.LogError("Origin operator operate global params error: ", err)
		return false
	}
	// origin admin setOperator, should not success
	setOperator(originAdmin, newOperator.Address, ctx.Ont)
	// newOperator set and get global params, should cause error
	globalParams, err = testGetAndSet(newOperator, ctx.Ont, "newOperator1", globalParams)
	if err == nil {
		ctx.LogError("New operator operate global params error, should not be authorized!")
		return false
	}
	return true
}

func testGetAndSet(account *account.Account, ontSdk *sdk.OntologySdk, keyword string, initParams global_params.Params,
) (global_params.Params, error) {
	// read original param
	paramNameList := new(global_params.ParamNameList)
	for _, param := range initParams {
		(*paramNameList) = append(*paramNameList, param.Key)
	}
	queriedParams, err := getParam(paramNameList, ontSdk)
	if err != nil {
		return initParams, fmt.Errorf("query param failed: %s", err)
	}
	// compare to original param
	if len(*queriedParams) != len(initParams) {
		return initParams, fmt.Errorf("query param failed: init param not equals genesis init param!")
	}
	for index, param := range initParams {
		if (*queriedParams)[index].Key != param.Key || (*queriedParams)[index].Value != param.Value {
			return initParams, fmt.Errorf("query param failed: init param not equals genesis init param!")
		}
	}

	// set global param
	newParams := new(global_params.Params)
	for i := 0; i < 3; i++ {
		key := "test-key-" + keyword + "-" + strconv.Itoa(i)
		value := "test-value-" + keyword + "-" + strconv.Itoa(i)
		(*paramNameList) = append(*paramNameList, key)
		newParams.SetParam(&global_params.Param{key, value})
	}
	setParam(account, ontSdk, newParams)

	// create Snapshot
	createSnapshot(account, ontSdk)

	// read param
	queriedParams, err = getParam(paramNameList, ontSdk)
	if err != nil {
		// snapshot has been created, so return the new params
		return initParams, fmt.Errorf("query param failed: %s", err)
	}

	// update new params by init params
	for _, v := range initParams {
		newParams.SetParam(v)
	}
	if len(*newParams) != len(*queriedParams) {
		return *newParams, fmt.Errorf("set param failed, the new param isn't effective!")
	}
	for _, param := range *newParams {
		if index, _ := queriedParams.GetParam(param.Key); index < 0 {
			return *newParams, fmt.Errorf("set param failed: the new param isn't effective!")
		}
	}
	return *queriedParams, err
}

func transferAdmin(orignAdmin *account.Account, newAdminAddress common.Address, ontSdk *sdk.OntologySdk) {
	ontSdk.Rpc.InvokeNativeContract(0, 0, orignAdmin, 0, utils.ParamContractAddress,
		"transferAdmin", []interface{}{newAdminAddress})
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func acceptAdmin(newAdmin *account.Account, ontSdk *sdk.OntologySdk) {
	ontSdk.Rpc.InvokeNativeContract(0, 0, newAdmin, 0, utils.ParamContractAddress,
		"acceptAdmin", []interface{}{newAdmin.Address})
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func setOperator(admin *account.Account, newOperator common.Address, ontSdk *sdk.OntologySdk) {
	ontSdk.Rpc.InvokeNativeContract(0, 0, admin, 0, utils.ParamContractAddress,
		"setOperator", []interface{}{newOperator})
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func setParam(account *account.Account, ontSdk *sdk.OntologySdk, params *global_params.Params) {
	ontSdk.Rpc.InvokeNativeContract(0, 0, account, 0, utils.ParamContractAddress,
		"setGlobalParam", []interface{}{*params})
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func createSnapshot(account *account.Account, ontSdk *sdk.OntologySdk) {
	// create snapshot
	ontSdk.Rpc.InvokeNativeContract(0, 0, account, 0, utils.ParamContractAddress,
		"createSnapshot", nil)
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func getParam(paramNameList *global_params.ParamNameList, ontSdk *sdk.OntologySdk) (*global_params.Params, error) {
	tx, err := ontSdk.Rpc.NewNativeInvokeTransaction(0, 0, global_params.VERSION_CONTRACT_GLOBAL_PARAMS,
		utils.ParamContractAddress, "getGlobalParam", []interface{}{*paramNameList})
	if err != nil{
		return nil, err
	}
	result, err := ontSdk.Rpc.PrepareInvokeContract(tx)
	if err != nil {
		return nil, err
	}
	queriedParams := new(global_params.Params)
	data, err := hex.DecodeString(result.Result.(string))
	if err != nil {
		err = fmt.Errorf("get param error: decode result error!")
	}
	err = queriedParams.Deserialize(bytes.NewBuffer([]byte(data)))
	if err != nil {
		err = fmt.Errorf("get param error: deserialize result error!")
	}
	return queriedParams, err
}
