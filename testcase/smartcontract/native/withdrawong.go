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

	balanceBefore, err := ctx.Ont.Native.Ont.BalanceOf(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG GetBalance error:%s", err)
		return false
	}
	if balanceBefore == 0 {
		ctx.LogError("TestWithdrawONG ont balance = 0")
		return false
	}

	ctx.LogInfo("TestWithdrawONG Balance ONT:%d", balanceBefore)

	amount := uint64(100000)
	if balanceBefore < amount {
		amount = balanceBefore
	}
	_, err = ctx.Ont.Native.Ont.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(),  defAccount, newAccount.Address, 10000)
	if err != nil {
		ctx.LogError("TestWithdrawONG Transfer ont error:%s", err)
		return false
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)


	ongBalanceBefor, err := ctx.Ont.Native.Ong.BalanceOf(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG BalanceOf ong error:%s", err)
		return false
	}
	unBoundONG, err := ctx.Ont.Native.Ong.UnboundONG(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG UnboundONG error:%s", err)
		return false
	}
	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)

	ctx.LogInfo("TestWithdrawONG UnboundONG:%d", unBoundONG)
	if unBoundONG == 0 {
		ctx.LogError("TestWithdrawONG UnboundONG = 0")
		return false
	}

	withdrawAmount := unBoundONG - 1
	_, err = ctx.Ont.Native.Ong.WithdrawONG(ctx.GetGasPrice(), ctx.GetGasLimit(), defAccount, withdrawAmount)
	if err != nil {
		ctx.LogError("TestWithdrawONG WithdrawONG error:%s", err)
		return false
	}

	ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)

	unBoundONGAfter, err := ctx.Ont.Native.Ong.UnboundONG(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG UnboundONG error:%s", err)
		return false
	}

	if unBoundONGAfter != unBoundONG-withdrawAmount {
		ctx.LogError("TestWithdrawONG unBoundONGAfter:%d != %d", unBoundONGAfter, unBoundONG-withdrawAmount)
		return false
	}

	ongBalanceAfter, err := ctx.Ont.Native.Ong.BalanceOf(defAccount.Address)
	if err != nil {
		ctx.LogError("TestWithdrawONG GetBalance error:%s", err)
		return false
	}

	ctx.LogInfo("TestWithdrawONG Balance after ONG:%d",  ongBalanceAfter)

	if ongBalanceAfter != ongBalanceBefor+withdrawAmount {
		ctx.LogError("TestWithdrawONG ong balance %d != %d", ongBalanceAfter, ongBalanceBefor+withdrawAmount)
		return false
	}

	return true
}
