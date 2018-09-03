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

//TestOntTransfer test native transfer case
func TestOntTransfer(ctx *testframework.TestFrameworkContext) bool {
	admin, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("Wallet.GetDefaultAccount error:%s", err)
		return false
	}
	user := ctx.NewAccount()

	adminBalanceBefore, err := ctx.Ont.Native.Ont.BalanceOf(admin.Address)
	if err != nil {
		ctx.LogError("Rpc.GetBalance error:%s", err)
		return false
	}

	if adminBalanceBefore == 0 {
		ctx.LogWarn("TestOntTransfer failed. Balance of admin is 0")
		return false
	}
	ctx.LogInfo("adminBalanceBefore %d", adminBalanceBefore)

	userBalanceBefore, err := ctx.Ont.Native.Ont.BalanceOf(user.Address)
	if err != nil {
		ctx.LogError("Rpc.GetBalance error:%s", err)
		return false
	}
	ctx.LogInfo("userBalanceBefore %d", userBalanceBefore)

	amount := uint64(100)
	_, err = ctx.Ont.Native.Ont.Transfer(ctx.GetGasPrice(), ctx.GetGasLimit(), admin, user.Address, amount)
	if err != nil {
		ctx.LogError("Rpc.Transfer error:%s", err)
		return false
	}

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false
	}

	adminBalanceAfter, err := ctx.Ont.Native.Ont.BalanceOf(admin.Address)
	if err != nil {
		if err != nil {
			ctx.LogError("Rpc.GetBalance error:%s", err)
			return false
		}
	}
	ctx.LogInfo("adminBalanceAfter :%d", adminBalanceAfter)

	userBalanceAfter, err := ctx.Ont.Native.Ont.BalanceOf(user.Address)
	if err != nil {
		ctx.LogError("Rpc.GetBalance error:%s", err)
		return false
	}
	ctx.LogInfo("userBalanceAfter :%d", userBalanceAfter)

	//Assert admin balance
	adminRes := adminBalanceBefore - amount
	if adminRes != adminBalanceAfter {
		ctx.LogError("TestOntTransfer failed. Admin balance after transfer %d != %d", adminBalanceAfter, adminRes)
		return false
	}

	//Assert user balance
	userRes := userBalanceBefore + amount
	if userRes != userBalanceAfter {
		ctx.LogError("TestOntTransfer failed. User balance after transfer %d != %d", userBalanceAfter, userRes)
		return false
	}
	return true
}
