package native

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
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
	// origin admin set and get global params, should not cause error
	globalParams, err := testGetAndSet(originAdmin, ctx.Ont, "originAdmin0", global_params.GLOBAL_PARAM)
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

func testGetAndSet(account *account.Account, ontSdk *sdk.OntologySdk, keyword string, initParam map[string]string,
) (map[string]string, error) {
	// read original param
	paramNameList := new(global_params.ParamNameList)
	for k, _ := range initParam {
		(*paramNameList) = append(*paramNameList, k)
	}
	queriedParams, err := getParam(paramNameList, ontSdk)
	if err != nil {
		return initParam, fmt.Errorf("query param failed: %s", err)
	}
	// compare to original param
	if len(*queriedParams) != len(initParam) {
		return initParam, fmt.Errorf("query param failed: init param not equals genesis init param!")
	}
	for k, v := range initParam {
		if (*queriedParams)[k] != v {
			return initParam, fmt.Errorf("query param failed: init param not equals genesis init param!")
		}
	}

	// set param
	bf := new(bytes.Buffer)
	newParams := new(global_params.Params)
	*newParams = make(map[string]string)
	for i := 0; i < 3; i++ {
		key := "test-key-" + keyword + "-" + strconv.Itoa(i)
		value := "test-value-" + keyword + "-" + strconv.Itoa(i)
		(*paramNameList) = append(*paramNameList, key)
		(*newParams)[key] = value
	}
	newParams.Serialize(bf)
	_, err = ontSdk.Rpc.InvokeNativeContract(0, 0, account, 0, genesis.ParamContractAddress,
		"setGlobalParam", bf.Bytes())
	if err != nil {
		return initParam, fmt.Errorf("set param failed: %s", err)
	}
	ontSdk.Rpc.WaitForGenerateBlock(30 * time.Second, 1)

	// create snapshot
	_, err = ontSdk.Rpc.InvokeNativeContract(0, 0, account, 0, genesis.ParamContractAddress,
		"createSnapshot", []byte{})
	if err != nil {
		return initParam, fmt.Errorf("create snapshot failed: %s", err)
	}
	ontSdk.Rpc.WaitForGenerateBlock(30 * time.Second, 1)

	// read param
	queriedParams, err = getParam(paramNameList, ontSdk)
	if err != nil {
		// snapshot has been created, so return the new params
		return initParam, fmt.Errorf("query param failed: %s", err)
	}

	// append init params to new params
	for k, v := range initParam {
		(*newParams)[k] = v
	}
	for k, _ := range *newParams {
		if _, ok := (*queriedParams)[k]; !ok {
			return *newParams, fmt.Errorf("set param failed: the new param take effect immediately!")
		}
	}
	return *queriedParams, err
}

func transferAdmin(orignAdmin *account.Account, newAdminAddress common.Address, ontSdk *sdk.OntologySdk)  {
	var destinationAdmin global_params.Admin
	copy(destinationAdmin[:], newAdminAddress[:])
	adminBuffer := new(bytes.Buffer)
	destinationAdmin.Serialize(adminBuffer)
	ontSdk.Rpc.InvokeNativeContract(0, 0, orignAdmin, 0, genesis.ParamContractAddress,
		"transferAdmin", adminBuffer.Bytes())
	ontSdk.Rpc.WaitForGenerateBlock(30 * time.Second, 1)
}

func accpetAdmin(newAdmin *account.Account, ontSdk *sdk.OntologySdk) error {
	var destinationAdmin global_params.Admin
	copy(destinationAdmin[:], newAdmin.Address[:])
	adminBuffer := new(bytes.Buffer)
	destinationAdmin.Serialize(adminBuffer)
	_, err := ontSdk.Rpc.InvokeNativeContract(0, 0, newAdmin, 0, genesis.ParamContractAddress,
		"acceptAdmin", adminBuffer.Bytes())
	ontSdk.Rpc.WaitForGenerateBlock(30 * time.Second, 1)
	return err
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
	err = json.Unmarshal(data, queriedParams)
	if err != nil {
		err = fmt.Errorf("get param error: unmarshal result error!")
	}
	return queriedParams, nil
}
