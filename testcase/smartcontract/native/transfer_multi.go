package native

import (
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	"time"
)

func TestOntTransferMulti(ctx *testframework.TestFrameworkContext) bool {
	defAcc, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestOntTransferMulti GetDefaultAccount error:%s", err)
		return false
	}
	acc1 := ctx.NewAccount()
	acc2 := ctx.NewAccount()

	_, err = ctx.Ont.Rpc.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), "ont", defAcc, acc1.Address, 10)
	if err != nil {
		ctx.LogError("TestOntTransferMulti Rpc.Transfer error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), "ont", defAcc, acc2.Address, 10)
	if err != nil {
		ctx.LogError("TestOntTransferMulti Rpc.Transfer error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(time.Second*30, 1)
	if err != nil {
		ctx.LogError("TestOntTransferMulti WaitForGenerateBlock error:%s", err)
		return false
	}

	//Start multi address transfer
	multiTransfer, err := ctx.Ont.Rpc.NewMultiTransferTransfer(ctx.GetGasPrice(), ctx.GetGasLimit(), "ont", []*ont.State{
		{From: acc1.Address, To: acc2.Address, Value: 2},
		{From: acc2.Address, To: acc1.Address, Value: 8},
	})
	if err != nil {
		ctx.LogError("TestOntTransferMulti NewMultiTransferTransfer error:%s", err)
		return false
	}

	err = ctx.Ont.Rpc.SignToTransaction(multiTransfer, acc1)
	if err != nil {
		ctx.LogError("TestOntTransferMulti SignToTransaction error:%s", err)
		return false
	}
	err = ctx.Ont.Rpc.SignToTransaction(multiTransfer, acc2)
	if err != nil {
		ctx.LogError("TestOntTransferMulti SignToTransaction error:%s", err)
		return false
	}
	txHash, err := ctx.Ont.Rpc.SendRawTransaction(multiTransfer)
	if err != nil {
		ctx.LogError("TestOntTransferMulti SendRawTransaction error:%s", err)
		return false
	}
	ctx.LogInfo("TestOntTransferMulti MultiTransfer TxHash:%s", txHash.ToHexString())

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestOntTransferMulti WaitForGenerateBlock error:%s", err)
		return false
	}

	acc1_balance, err := ctx.Ont.Rpc.GetBalance(acc1.Address)
	if err != nil {
		ctx.LogError("TestOntTransferMulti Rpc.GetBalance error:%s", err)
		return false
	}
	acc2_balance, err := ctx.Ont.Rpc.GetBalance(acc2.Address)
	if err != nil {
		ctx.LogError("TestOntTransferMulti Rpc.GetBalance error:%s", err)
		return false
	}

	if acc1_balance.Ont != 16 {
		ctx.LogError("TestOntTransferMulti Account1 balance %d != %d", acc1_balance.Ont, 16)
		return false
	}

	if acc2_balance.Ont != 4 {
		ctx.LogError("TestOntTransferMulti Account1 balance %d != %d", acc2_balance.Ont, 4)
		return false
	}
	return true
}
