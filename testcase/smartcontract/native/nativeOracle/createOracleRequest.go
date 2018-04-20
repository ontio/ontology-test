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
	"encoding/json"
	"fmt"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/core/genesis"
	cstates "github.com/ontio/ontology/smartcontract/states"
	vmtypes "github.com/ontio/ontology/smartcontract/types"
	"math/big"
	"time"
	"encoding/hex"
)

type CreateOracleRequestParam struct {
	Request   string         `json:"request"`
	OracleNum int            `json:"oracleNum"`
	Address   string `json:"address"`
}

func CreateOracleRequest(ctx *testframework.TestFrameworkContext) (string, error) {
	var request = `{
	"scheduler":{
		"type": "cron",
		"params": "0/10 * * * * ?"
	},
	"tasks":[
	  {
		"type": "httpGet",
		"params": {
		  "url": "https://bitstamp.net/api/ticker/"
		}
	  },
	  {
		"type": "jsonParse",
		"params": {
		  "path": ["last"]
		}
	  }
	]
	}`

	wallet, err := ctx.Ont.OpenWallet("./wallet.dat", "passwordtest")
	if err != nil {
		return "", fmt.Errorf("OpenWallet ./wallet.dat error:%s", err)
	}

	user, err := wallet.GetDefaultAccount()
	if err != nil {
		return "", fmt.Errorf("Wallet.CreateAccount error:%s", err)
	}

	params := &CreateOracleRequestParam{
		Request:   request,
		OracleNum: 3,
		Address:   hex.EncodeToString(user.Address[:]),
	}

	contractAddress := genesis.OracleContractAddress

	args, err := json.Marshal(params)
	crt := &cstates.Contract{
		Address: contractAddress,
		Method:  "createOracleRequest",
		Args:    args,
	}
	buf := bytes.NewBuffer(nil)
	err = crt.Serialize(buf)
	if err != nil {
		return "", fmt.Errorf("Serialize contract error:%s", err)
	}

	invokeTx := sdkcom.NewInvokeTransaction(new(big.Int).SetInt64(0), vmtypes.Native, buf.Bytes())
	err = sdkcom.SignTransaction("SHA256withECDSA", invokeTx, user)
	if err != nil {
		return "", fmt.Errorf("SignTransaction error:%s", err)
	}
	txHash, err := ctx.Ont.Rpc.SendRawTransaction(invokeTx)
	if err != nil {
		return "", fmt.Errorf("SendTransaction error:%s", err)
	}

	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		return "", fmt.Errorf("WaitForGenerateBlock error:%s", err)
	}

	return hex.EncodeToString(txHash[:]), nil
}
