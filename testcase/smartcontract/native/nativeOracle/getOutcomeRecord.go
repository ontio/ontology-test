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
	"github.com/ontio/ontology/core/genesis"
	"encoding/json"
	"encoding/hex"
)

type OutcomeRecord struct {
	OutcomeRecord map[string]interface{} `json:"outcomeRecord"`
}

func GetOutcomeRecord(ctx *testframework.TestFrameworkContext) *OutcomeRecord {

	contractAddress := genesis.OracleContractAddress
	txHashHex := "4d40a58a68f6d76a36854735c0111de01e7a6986a0deebf70e3d19371a84436f"
	txHash, _:= hex.DecodeString(txHashHex)
	key := append([]byte("OutcomeRecord"), txHash...)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, []byte(key))
	if err != nil {
		ctx.LogError("GetOutcomeRecord GetStorageItem key:%s error:%s", key, err)
		return nil
	}
	if len(value) == 0 {
		ctx.LogError("GetOutcomeRecord OutcomeRecord is not set!")
		return nil
	}

	result := new(OutcomeRecord)
	err = json.Unmarshal(value, &result)
	if err != nil {
		ctx.LogError("GetOutcomeRecord Unmarshal result error:%s", err)
		return nil
	}

	return result
}