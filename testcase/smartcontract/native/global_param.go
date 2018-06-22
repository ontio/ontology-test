package native

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

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
	initParms, err := getParam(ctx, global_params.ParamNameList{initKey})
	if err != nil {
		ctx.LogError("Get global params error, the value should initialize in genesis!")
		return false
	}
	testKey := "init-key"
	testKeyValue := "init-value"

	ps, err := getParam(ctx, global_params.ParamNameList{testKey})
	if err == nil {
		_, p := ps.GetParam(testKey)
		testKeyValue = p.Value + "1"
	}

	// add a global params
	initParms.SetParam(&global_params.Param{testKey, testKeyValue})
	// origin admin is origin operator
	originOperator := originAdmin
	setParam(ctx, originOperator, initParms)
	// query init param, should cause error, because not create snapshot
	ps, err = getParam(ctx, global_params.ParamNameList{testKey})
	if err == nil {
		_, p := ps.GetParam(testKey)
		if p.Value == testKeyValue {
			ctx.LogError("Get global params:%s error, the value should not take effect!", testKey)
			return false
		}
	}
	// create snapshot of init params
	createSnapshot(ctx, originOperator)
	// query init param, should not cause error
	initParms, err = getParam(ctx, global_params.ParamNameList{testKey})
	if err != nil {
		ctx.LogError("Get global params:%s error: %s!", testKey, err)
		return false
	}
	// origin admin set and get global params, should not cause error
	globalParams, err := testGetAndSet(ctx, originOperator, "originOperator0", initParms)
	if err != nil {
		ctx.LogError("Origin operator operate global params error:%s", err)
		return false
	}

	ctx.LogInfo("TestGlobalParam Step 1 success")

	newOperator := ctx.NewAccount()

	// new operator set and get global params, should cause error, because  he is not operator
	globalParams, err = testGetAndSet(ctx, newOperator, "newOperator0", globalParams)
	if err == nil {
		ctx.LogError("New operator operate global params error, should not be authorized!")
		return false
	}
	// set newOperator as operator
	setOperator(ctx, originAdmin, newOperator.Address)
	// newOperator set and get global params, should not cause error, because he is operator
	globalParams, err = testGetAndSet(ctx, newOperator, "newOperator0", globalParams)
	if err != nil {
		ctx.LogError("New operator operate global params error: ", err)
		return false
	}
	// originOperator set and get global param, should cause error, because he is not operator
	globalParams, err = testGetAndSet(ctx, originOperator, "originOperator1", globalParams)
	if err == nil {
		ctx.LogError("Origin operator operate global params error:, should not be authorized!")
		return false
	}

	ctx.LogInfo("TestGlobalParam Step 2 success")

	// transfer admin to newOperator
	newAdmin := newOperator
	transferAdmin(ctx, originAdmin, newAdmin.Address)
	// new admin set originOperator as current operator, should not success, because new admin does not accept
	setOperator(ctx, newAdmin, originOperator.Address)
	// originOperator operate global params, should cause error, because newAdmin is not admin
	// and he's setOperator is not effective
	globalParams, err = testGetAndSet(ctx, originOperator, "originOperator1", globalParams)
	if err == nil {
		ctx.LogError("Origin operator operate global params error, should not be authorized!")
		return false
	}

	ctx.LogInfo("TestGlobalParam Step 3 success")

	// new admin accept permission
	acceptAdmin(ctx, newAdmin)
	// new admin set originOperator as current operator
	setOperator(ctx, newAdmin, originOperator.Address)
	// originOperator operate global params, should not cause error
	globalParams, err = testGetAndSet(ctx, originOperator, "originOperator1", globalParams)
	if err != nil {
		ctx.LogError("Origin operator operate global params error: ", err)
		return false
	}
	// origin admin setOperator, should not success
	setOperator(ctx, originAdmin, newOperator.Address)
	// newOperator set and get global params, should cause error
	globalParams, err = testGetAndSet(ctx, newOperator, "newOperator1", globalParams)
	if err == nil {
		ctx.LogError("New operator operate global params error, should not be authorized!")
		return false
	}

	ctx.LogInfo("TestGlobalParam Step 4 success")
	//reset admin and operator
	transferAdmin(ctx, newAdmin, originAdmin.Address)
	acceptAdmin(ctx, originAdmin)
	setOperator(ctx, originOperator, originAdmin.Address)
	return true
}

func testGetAndSet(ctx *testframework.TestFrameworkContext, account *account.Account, keyword string, initParams global_params.Params,
) (global_params.Params, error) {
	// read original param
	paramNameList := global_params.ParamNameList(make([]string, 0))
	for _, param := range initParams {
		paramNameList = append(paramNameList, param.Key)
	}
	queriedParams, err := getParam(ctx, paramNameList)
	if err != nil {
		return initParams, fmt.Errorf("query param failed: %s", err)
	}
	// compare to original param
	if len(queriedParams) != len(initParams) {
		return initParams, fmt.Errorf("query param failed: init param not equals genesis init param!")
	}
	for index, param := range initParams {
		if queriedParams[index].Key != param.Key || queriedParams[index].Value != param.Value {
			return initParams, fmt.Errorf("query param failed: init param not equals genesis init param!")
		}
	}

	// set global param
	newParams := global_params.Params(make([]*global_params.Param, 0))
	timeStemp := time.Now().String()
	for i := 0; i < len(initParams); i++ {
		key := "test-key-" + keyword + "-" + strconv.Itoa(i) + timeStemp
		value := "test-value-" + keyword + "-" + strconv.Itoa(i)
		paramNameList = append(paramNameList, key)
		newParams.SetParam(&global_params.Param{key, value})
	}
	setParam(ctx, account, newParams)

	// create Snapshot
	createSnapshot(ctx, account)

	// read param
	queriedParams, err = getParam(ctx, paramNameList)
	if err != nil {
		// snapshot has been created, so return the new params
		return initParams, fmt.Errorf("query param failed: %s", err)
	}

	// update new params by init params
	for _, v := range initParams {
		newParams.SetParam(v)
	}
	if len(newParams) != len(queriedParams) {
		return newParams, fmt.Errorf("set param failed, the new param isn't effective!")
	}
	for _, param := range newParams {
		if index, _ := queriedParams.GetParam(param.Key); index < 0 {
			return newParams, fmt.Errorf("set param failed: the new param isn't effective!")
		}
	}
	return queriedParams, err
}

func transferAdmin(ctx *testframework.TestFrameworkContext, orignAdmin *account.Account, newAdminAddress common.Address) {
	ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), orignAdmin, global_params.VERSION_CONTRACT_GLOBAL_PARAMS,
		utils.ParamContractAddress, "transferAdmin", []interface{}{newAdminAddress})
	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func acceptAdmin(ctx *testframework.TestFrameworkContext, newAdmin *account.Account) {
	ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), newAdmin, global_params.VERSION_CONTRACT_GLOBAL_PARAMS,
		utils.ParamContractAddress, "acceptAdmin", []interface{}{newAdmin.Address})
	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func setOperator(ctx *testframework.TestFrameworkContext, admin *account.Account, newOperator common.Address) {
	ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), admin, global_params.VERSION_CONTRACT_GLOBAL_PARAMS,
		utils.ParamContractAddress, "setOperator", []interface{}{newOperator})
	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func setParam(ctx *testframework.TestFrameworkContext, account *account.Account, params global_params.Params) {
	ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), account, global_params.VERSION_CONTRACT_GLOBAL_PARAMS,
		utils.ParamContractAddress, "setGlobalParam", []interface{}{params})
	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func createSnapshot(ctx *testframework.TestFrameworkContext, account *account.Account) {
	// create snapshot
	ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), account, global_params.VERSION_CONTRACT_GLOBAL_PARAMS,
		utils.ParamContractAddress, "createSnapshot", []interface{}{})
	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func getParam(ctx *testframework.TestFrameworkContext, paramNameList global_params.ParamNameList) (global_params.Params, error) {
	tx, err := ctx.Ont.Rpc.NewNativeInvokeTransaction(ctx.GetGasPrice(), ctx.GetGasLimit(), global_params.VERSION_CONTRACT_GLOBAL_PARAMS,
		utils.ParamContractAddress, "getGlobalParam", []interface{}{paramNameList})
	if err != nil {
		return nil, err
	}
	result, err := ctx.Ont.Rpc.PrepareInvokeContract(tx)
	if err != nil {
		return nil, err
	}
	queriedParams := global_params.Params(make([]*global_params.Param, 0))
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
