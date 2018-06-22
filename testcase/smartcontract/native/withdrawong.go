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

package native

import (
	"github.com/ontio/ontology-test/testframework"
	"time"
)

func TestWithdrawONG(ctx *testframework.TestFrameworkContext) bool {
	defAccount, _ := ctx.GetDefaultAccount()
	newAccount := ctx.NewAccount()

	balanceBefore, err := ctx.Ont.Rpc.GetBalance(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG GetBalance error:%s", err)
		return false
	}
	if balanceBefore.Ont == 0 {
		ctx.LogError("TestWithdrawONG ont balance = 0")
		return false
	}

	ctx.LogInfo("TestWithdrawONG Balance ONT:%d ONG:%d", balanceBefore.Ont, balanceBefore.Ong)

	amount := uint64(100000)
	if balanceBefore.Ont < amount {
		amount = balanceBefore.Ont
	}
	_, err = ctx.Ont.Rpc.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), "ont", defAccount, newAccount.Address, 10000)
	if err != nil {
		ctx.LogError("TestWithdrawONG Transfer ont error:%s", err)
		return false
	}
	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	unBoungONG, err := ctx.Ont.Rpc.UnboundONG(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG UnboundONG error:%s", err)
		return false
	}
	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	ctx.LogInfo("TestWithdrawONG UnboundONG:%d", unBoungONG)
	if unBoungONG == 0 {
		ctx.LogError("TestWithdrawONG UnboundONG = 0")
		return false
	}

	withdrawAmount := unBoungONG - 1
	_, err = ctx.Ont.Rpc.WithdrawONG(ctx.GetGasPrice(), ctx.GetGasLimit(), defAccount, withdrawAmount)
	if err != nil {
		ctx.LogError("TestWithdrawONG WithdrawONG error:%s", err)
		return false
	}

	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	unBoungONGAfter, err := ctx.Ont.Rpc.UnboundONG(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG UnboundONG error:%s", err)
		return false
	}

	ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)

	if unBoungONGAfter != unBoungONG-withdrawAmount {
		ctx.LogError("TestWithdrawONG unBoungONGAfter:%d != %d", unBoungONGAfter, unBoungONG-withdrawAmount)
		return false
	}

	balanceAfter, err := ctx.Ont.Rpc.GetBalance(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG GetBalance error:%s", err)
		return false
	}

	ctx.LogInfo("TestWithdrawONG Balance after ONT:%d ONG:%d", balanceAfter.Ont, balanceAfter.Ong)

	if balanceAfter.Ong != balanceBefore.Ong+withdrawAmount {
		ctx.LogError("TestWithdrawONG ong balance %d != %d", balanceAfter.Ong, balanceBefore.Ong+withdrawAmount)
		return false
	}

	return true
}
