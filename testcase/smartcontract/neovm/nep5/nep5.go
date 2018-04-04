/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
 //Test nep5 neovm contract
package nep5

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/vm/types"
	"math/big"
	"time"
)

//Test nep5 neovm contract, deploy, invoke
func TestNep5Contract(ctx *testframework.TestFrameworkContext) bool {
	nep5Wallet := "./testcase/smartcontract/neovm/nep5/nep5wallet.dat"
	nep5WalletPwd := "123"
	wallet, err := ctx.Ont.OpenWallet(nep5Wallet, nep5WalletPwd)
	if err != nil {
		ctx.LogError("OpenWallet:%s error:%s", nep5Wallet, err)
		return false
	}

	admin, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestNep5Contract nep5Wallet.GetDefaultAccount error:%s", err)
		return false
	}

	txHash, err := nep5Deploy(ctx, admin)
	if err != nil {
		ctx.LogError("TestNep5Contract deployNep5 error:%s", err)
		return false
	}

	ctx.LogInfo("TestNep5Contract deploy TxHash:%x", txHash)

	txHash, err = nep5Init(ctx, admin)
	if err != nil {
		ctx.LogError("TestNep5Contract nep5Init error:%s", err)
		return false
	}
	ctx.LogInfo("InitIx: %x\n", txHash)
	ctx.LogInfo("TestNep5Contract nep5Init success")

	notifies, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestNep5Contract init GetSmartContractEvent error:%s", err)
		return false
	}
	ctx.LogInfo("TestNep5Contract init notify %s", notifies)

	user, err := wallet.CreateAccount()
	if err != nil {
		ctx.LogError("wallet.CreateAccount error:%s", err)
		return false
	}
	txHash2, err := nep5Transfer(ctx, admin, user, 10)
	if err != nil {
		ctx.LogError("TestNep5Contract nep5Transfer error:%s", err)
		return false
	}
	ctx.LogInfo("TransferTx: %x\n", txHash2)
	ctx.LogInfo("TestNep5Contract nep5Transfer success")

	notifies, err = ctx.Ont.Rpc.GetSmartContractEvent(txHash2)
	if err != nil {
		ctx.LogError("TestNep5Contract transfer GetSmartContractEvent error:%s", err)
		return false
	}

	ctx.LogInfo("TestNep5Contract transfer notify %s", notifies)
	return true
}

func nep5Transfer(ctx *testframework.TestFrameworkContext, from, to *account.Account, amount int) (common.Uint256, error) {
	method := "transfer"
	txHash, err := ctx.Ont.Rpc.InvokeNeoVMSmartContract(from, new(big.Int), nep5Address, []interface{}{method, []interface{}{from.Address[:], to.Address[:], amount}})
	if err != nil {
		return common.Uint256{}, fmt.Errorf("InvokeNeoVM %s error:%s", method, err)
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func nep5Init(ctx *testframework.TestFrameworkContext, acc *account.Account) (common.Uint256, error) {
	method := "init"
	txHash, err := ctx.Ont.Rpc.InvokeNeoVMSmartContract(acc, new(big.Int), nep5Address, []interface{}{method, []interface{}{0}})
	if err != nil {
		return common.Uint256{}, fmt.Errorf("InvokeNeoVM %s error:%s", method, err)
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30 * time.Second)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}
	return txHash, nil
}

func nep5Deploy(ctx *testframework.TestFrameworkContext, signer *account.Account) (common.Uint256, error) {
	txHash, err := ctx.Ont.Rpc.DeploySmartContract(
		signer,
		types.NEOVM,
		true,
		nep5Contract,
		"nyc",
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
	return txHash, nil
}
