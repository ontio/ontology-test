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
	"bytes"
	"encoding/hex"
	"encoding/json"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/core/genesis"
	cstates "github.com/ontio/ontology/smartcontract/states"
	vmtypes "github.com/ontio/ontology/smartcontract/types"
	"math/big"
	"time"
)

type SetOracleOutcomeParam struct {
	TxHash  string      `json:"txHash"`
	Owner   string      `json:"owner"`
	Outcome interface{} `json:"outcome"`
}

func SetOracleOutcome(ctx *testframework.TestFrameworkContext, tx string) bool {
	wallet, err := ctx.Ont.CreateWallet("./wallet_test.dat", "wangbing")
	if err != nil {
		ctx.LogError("CreateWallet ./wallet_test.dat error:%s", err)
		return false
	}

	user, err := wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("Wallet.CreateAccount error:%s", err)
		return false
	}

	params := &SetOracleOutcomeParam{
		TxHash:  tx,
		Owner:   hex.EncodeToString(user.Address[:]),
		Outcome: "helloworld",
	}

	contractAddress := genesis.OracleContractAddress

	args, err := json.Marshal(params)
	crt := &cstates.Contract{
		Address: contractAddress,
		Method:  "setOracleOutcome",
		Args:    args,
	}
	buf := bytes.NewBuffer(nil)
	err = crt.Serialize(buf)
	if err != nil {
		ctx.LogError("Serialize contract error:%s", err)
		return false
	}

	invokeTx := sdkcom.NewInvokeTransaction(new(big.Int).SetInt64(0), vmtypes.Native, buf.Bytes())
	err = sdkcom.SignTransaction("SHA256withECDSA", invokeTx, user)
	if err != nil {
		ctx.LogError("SignTransaction error:%s", err)
		return false
	}
	txHash, err := ctx.Ont.Rpc.SendRawTransaction(invokeTx)
	if err != nil {
		ctx.LogError("SendTransaction error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false
	}

	events, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("TestInvokeSmartContract GetSmartContractEvent error:%s", err)
		return false
	}

	states := events[0].States[1]

	ctx.LogInfo("setOracleOutcome result is : %+v", states.(bool))

	return true
}
