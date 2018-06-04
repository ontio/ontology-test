package native

import (
	"bytes"
	"encoding/hex"
	"fmt"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/genesis"
	"github.com/ontio/ontology/smartcontract/service/native/global_params"
	"strconv"
	"time"
)

func TestGlobalParam(ctx *testframework.TestFrameworkContext) bool {
	originAdmin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("Wallet.GetDefaultAccount error:%s", err)
		return false
	}
	// init global params
	initParms := new(global_params.Params)
	initParms.SetParam(&global_params.Param{"init-key", "init-value"})
	setParam(originAdmin, ctx.Ont, initParms)
	// query init param, should cause error, because not create snapshot
	_, err = getParam(&global_params.ParamNameList{"init-key"}, ctx.Ont)
	if err == nil {
		ctx.LogError("Get global params error, the value should not take effect!")
		return false
	}
	// create snapshot of init params
	createSnapshot(originAdmin, ctx.Ont)
	// query init param, should not cause error
	_, err = getParam(&global_params.ParamNameList{"init-key"}, ctx.Ont)
	if err != nil {
		ctx.LogError("Get global params error: %s!", err)
		return false
	}

	// origin admin set and get global params, should not cause error
	globalParams, err := testGetAndSet(originAdmin, ctx.Ont, "originAdmin0", *initParms)
	if err != nil {
		ctx.LogError("Origin admin operate global params error:%s", err)
		return false
	}
	newAdmin, err := ctx.NewAccount()
	if err != nil {
		ctx.LogError("Wallet.NewAccount error:%s", err)
		return false
	}
	// new admin set and get global params, should cause error, because of no permission
	globalParams, err = testGetAndSet(newAdmin, ctx.Ont, "newAdmin0", globalParams)
	if err == nil {
		ctx.LogError("New admin operate global params error, should not be authorized!")
		return false
	}
	// new admin accept permission
	accpetAdmin(newAdmin, ctx.Ont)
	// new admin set and get global params, should cause error, origin admin doesn't transfer
	globalParams, err = testGetAndSet(newAdmin, ctx.Ont, "newAdmin1", globalParams)
	if err == nil {
		ctx.LogError("New admin operate global params error, should not be authorized!")
		return false
	}
	// origin admin transfer permission to new admin
	transferAdmin(originAdmin, newAdmin.Address, ctx.Ont)
	// origin admin set and get global params, should not cause error, because new admin doesn't accept permission
	globalParams, err = testGetAndSet(originAdmin, ctx.Ont, "originAdmin1", globalParams)
	if err != nil {
		ctx.LogError("Origin admin operate global params error: %s", err)
		return false
	}
	// new admin set and get global params, should cause error, because new admin doesn't accept permission
	globalParams, err = testGetAndSet(newAdmin, ctx.Ont, "newAdmin2", globalParams)
	if err == nil {
		ctx.LogError("New admin operate global params error, should not be authorized!")
		return false
	}
	// new admin accept permission
	accpetAdmin(newAdmin, ctx.Ont)
	// origin admin set and get global params, should cause error, beacause he doesn't have permission
	globalParams, err = testGetAndSet(originAdmin, ctx.Ont, "originAdmin2", globalParams)
	if err == nil {
		ctx.LogError("Origin admin operate global params error, should not be authorized!")
		return false
	}
	// new admin set and get global params, should not cause error
	globalParams, err = testGetAndSet(newAdmin, ctx.Ont, "newAdmin3", globalParams)
	if err != nil {
		ctx.LogError("New admin operate global params error: %s", err)
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

	// update new params bt init params
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

func setParam(account *account.Account, ontSdk *sdk.OntologySdk, params *global_params.Params) {
	bf := new(bytes.Buffer)
	params.Serialize(bf)
	ontSdk.Rpc.InvokeNativeContract(0, 0, account, 0, genesis.ParamContractAddress,
		"setGlobalParam", bf.Bytes())
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func createSnapshot(account *account.Account, ontSdk *sdk.OntologySdk) {
	// create snapshot
	ontSdk.Rpc.InvokeNativeContract(0, 0, account, 0, genesis.ParamContractAddress,
		"createSnapshot", []byte{})
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func transferAdmin(orignAdmin *account.Account, newAdminAddress common.Address, ontSdk *sdk.OntologySdk) {
	var destinationAdmin global_params.Admin
	copy(destinationAdmin[:], newAdminAddress[:])
	adminBuffer := new(bytes.Buffer)
	destinationAdmin.Serialize(adminBuffer)
	ontSdk.Rpc.InvokeNativeContract(0, 0, orignAdmin, 0, genesis.ParamContractAddress,
		"transferAdmin", adminBuffer.Bytes())
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func accpetAdmin(newAdmin *account.Account, ontSdk *sdk.OntologySdk) {
	var destinationAdmin global_params.Admin
	copy(destinationAdmin[:], newAdmin.Address[:])
	adminBuffer := new(bytes.Buffer)
	destinationAdmin.Serialize(adminBuffer)
	ontSdk.Rpc.InvokeNativeContract(0, 0, newAdmin, 0, genesis.ParamContractAddress,
		"acceptAdmin", adminBuffer.Bytes())
	ontSdk.Rpc.WaitForGenerateBlock(30*time.Second, 1)
}

func getParam(paramNameList *global_params.ParamNameList, ontSdk *sdk.OntologySdk) (*global_params.Params, error) {
	bf := new(bytes.Buffer)
	paramNameList.Serialize(bf)
	result, err := ontSdk.Rpc.PrepareInvokeNativeSmartContract(0, genesis.ParamContractAddress,
		"getGlobalParam", bf.Bytes())
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
