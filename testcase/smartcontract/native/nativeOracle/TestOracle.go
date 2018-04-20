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

package nativeOracle

import (
	"github.com/ontio/ontology-test/testframework"
)

func TestOracle(ctx *testframework.TestFrameworkContext) bool {
	txHash, err := CreateOracleRequest(ctx)
	if err != nil {
		ctx.LogError("TestOracle error:%s", err)
		return false
	}

	ok := SetOracleOutcome(ctx, txHash)
	if !ok {
		return false
	}

	ok = SetOracleOutcome(ctx, txHash)
	if !ok {
		return false
	}

	ok = SetOracleOutcome(ctx, txHash)
	if !ok {
		return false
	}

	finalOutcome := GetFinalOutcome(ctx, txHash)
	if finalOutcome.(string) != "helloworld" {
		ctx.LogError("FinalOutcome is not helloworld")
		return false
	}
	ctx.LogInfo("Test success")
	return true
}
