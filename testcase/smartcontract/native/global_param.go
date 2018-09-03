package native

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/service/neovm"
	sdk"github.com/ontio/ontology-go-sdk"
	"time"
)

func TestGlobalParam(ctx *testframework.TestFrameworkContext) bool {
	defAcc, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestGlobalParam GetDefaultAccount error:%s", err)
		return false
	}
	oldAdmin := defAcc
	oldOperator := defAcc

	newAdmin := ctx.NewAccount()
	newOperator := ctx.NewAccount()

	err = testGetGlobalParam(ctx)
	if err != nil {
		ctx.LogError("TestGlobalParam testGetGlobalParam error:%s", err)
		return false
	}
	ctx.LogInfo("TestGlobalParam testGetGlobalParam success")

	err = testSetGlobalParam(ctx, oldOperator)
	if err != nil {
		ctx.LogError("TestGlobalParam testSetGlobalParam error:%s", err)
		return false
	}
	ctx.LogInfo("TestGlobalParam testSetGlobalParam success")

	err = testSetOperator(ctx, oldAdmin, oldOperator, newOperator)
	if err != nil {
		ctx.LogError("TestGlobalParam testSetOperator error:%s", err)
		return false
	}
	ctx.LogInfo("TestGlobalParam testSetOperator success")

	err = testTransferAndAcceptAdmin(ctx, oldAdmin, newAdmin, oldOperator, newOperator)
	if err != nil {
		ctx.LogError("TestGlobalParam testTransferAndAcceptAdmin error:%s", err)
		return false
	}
	ctx.LogInfo("TestGlobalParam testTransferAndAcceptAdmin success")

	return true
}

func testGetGlobalParam(ctx *testframework.TestFrameworkContext) error {
	params := neovm.GAS_TABLE_KEYS
	values, err := ctx.Ont.Native.GlobalParams.GetGlobalParams(params)
	if err != nil || len(values) != len(params) {
		return fmt.Errorf("testGetGlobalParam GetGlobalParams error:%s", err)
	}
	return nil
}

func testSetGlobalParam(ctx *testframework.TestFrameworkContext, operator *sdk.Account) error {
	testKey := "testKey"
	params := []string{testKey}
	testValue := fmt.Sprintf("%d", time.Now().Unix())

	err := setParam(ctx, operator, map[string]string{testKey: testValue})
	if err != nil {
		return fmt.Errorf("testSetGlobalParam SetGlobalParams error:%s", err)
	}

	values, err := ctx.Ont.Native.GlobalParams.GetGlobalParams(params)
	if err != nil {
		return fmt.Errorf("testGetGlobalParam GetGlobalParams error:%s", err)
	}
	if values[testKey] == testValue {
		return fmt.Errorf("testGetGlobalParam set param should not take effect before CreateSnapshot.")
	}

	_, err = ctx.Ont.Native.GlobalParams.CreateSnapshot(0, ctx.GetGasLimit(), operator)
	if err != nil {
		return fmt.Errorf("CreateSnapshot error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)

	values, err = ctx.Ont.Native.GlobalParams.GetGlobalParams(params)
	if err != nil {
		return fmt.Errorf("testGetGlobalParam GetGlobalParams error:%s", err)
	}
	if values[testKey] != testValue {
		return fmt.Errorf("testGetGlobalParam set param failed. Param:%s Value:%s != %s", testKey, values[testKey], testValue)
	}
	return nil
}

func testSetOperator(ctx *testframework.TestFrameworkContext, admin, oldOperator, newOperator *sdk.Account) error {
	testKey := "testKey"
	testValue := "testValue"
	testParams := map[string]string{testKey: testValue}

	err := setParam(ctx, oldOperator, testParams)
	if err != nil {
		return fmt.Errorf("oldOperator set param failed before set operator. Error:%s", err)
	}
	err = setParam(ctx, newOperator, testParams)
	if err == nil {
		return fmt.Errorf("newOperator set param should failed before set operator")
	}

	_, err = ctx.Ont.Native.GlobalParams.SetOperator(0, ctx.GetGasLimit(), admin, newOperator.Address)
	if err != nil {
		return fmt.Errorf("SetOperator error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)

	err = setParam(ctx, newOperator, testParams)
	if err != nil {
		return fmt.Errorf("newOperator set param failed after set operator. Error:%s", err)
	}
	err = setParam(ctx, oldOperator, testParams)
	if err == nil {
		return fmt.Errorf("oldOperator set param should failed after operator")
	}

	//reset operator
	_, err = ctx.Ont.Native.GlobalParams.SetOperator(0, ctx.GetGasLimit(), admin, oldOperator.Address)
	if err != nil {
		return fmt.Errorf("SetOperator error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	return nil
}

func testTransferAndAcceptAdmin(ctx *testframework.TestFrameworkContext, oldAdmin, newAdmin, oldOperator, newOperator *sdk.Account) error {
	err := testSetOperator(ctx, oldAdmin, oldOperator, newOperator)
	if err != nil {
		return fmt.Errorf("oldAdmin set operator failed before tansfer admin. Error:%s", err)
	}
	err = testSetOperator(ctx, newAdmin, oldOperator, newOperator)
	if err == nil {
		return fmt.Errorf("newAdmin set operator should failed before tansfer admin.")
	}

	_, err = ctx.Ont.Native.GlobalParams.TransferAdmin(0, ctx.GetGasLimit(), oldAdmin, newAdmin.Address)
	if err != nil {
		return fmt.Errorf("TransferAdmin error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)

	err = testSetOperator(ctx, oldAdmin, oldOperator, newOperator)
	if err != nil {
		return fmt.Errorf("oldAdmin set operator failed before accept admin. Error:%s", err)
	}
	err = testSetOperator(ctx, newAdmin, oldOperator, newOperator)
	if err == nil {
		return fmt.Errorf("newAdmin set operator should failed before accept admin.")
	}

	_, err = ctx.Ont.Native.GlobalParams.AcceptAdmin(0, ctx.GetGasLimit(), newAdmin)
	if err != nil {
		return fmt.Errorf("AcceptAdmin error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)

	err = testSetOperator(ctx, oldAdmin, oldOperator, newOperator)
	if err == nil {
		return fmt.Errorf("oldAdmin set operator should fialed after accept admin.")
	}
	err = testSetOperator(ctx, newAdmin, oldOperator, newOperator)
	if err != nil {
		return fmt.Errorf("newAdmin set operator failed after accept admin. Error:%s", err)
	}

	//reset admin
	_, err = ctx.Ont.Native.GlobalParams.TransferAdmin(0, ctx.GetGasLimit(), newAdmin, oldAdmin.Address)
	if err != nil {
		return fmt.Errorf("TransferAdmin error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	_, err = ctx.Ont.Native.GlobalParams.AcceptAdmin(0, ctx.GetGasLimit(), oldAdmin)
	if err != nil {
		return fmt.Errorf("AcceptAdmin error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	return nil
}

func setParam(ctx *testframework.TestFrameworkContext, operator *sdk.Account, params map[string]string) error {
	txHash, err := ctx.Ont.Native.GlobalParams.SetGlobalParams(0, ctx.GetGasLimit(), operator, params)
	if err != nil {
		return fmt.Errorf("testSetGlobalParam SetGlobalParams error:%s", err)
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	evt, err := ctx.Ont.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		return fmt.Errorf("GetSmartContractEvent error:%s", err)
	}
	if evt.State == 0 {
		return fmt.Errorf("SetGlobalParams failed")
	}
	return nil
}
