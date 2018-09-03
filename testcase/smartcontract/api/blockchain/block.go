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

package blockchain

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
	vtypes "github.com/ontio/ontology/vm/neovm/types"
)

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System;
using System.ComponentModel;
using System.Numerics;

class OnTest : SmartContract
{
    public static object[] Main(int height)
    {
        //object[] ret = new object[8];
        Block block = Blockchain.GetBlock((uint)height);
        Storage.Put(Storage.CurrentContext, "hash", block.Hash);
        Storage.Put(Storage.CurrentContext, "index", block.Index);
        Storage.Put(Storage.CurrentContext, "merkRoot", block.MerkleRoot);
        Storage.Put(Storage.CurrentContext, "nextConsensus", block.NextConsensus);
        Storage.Put(Storage.CurrentContext, "prevHash", block.PrevHash);
        Storage.Put(Storage.CurrentContext, "timeStamp", block.Timestamp);
        Storage.Put(Storage.CurrentContext, "version", block.Version);

        return null;
    }
}
*/
func TestGetBlock(ctx *testframework.TestFrameworkContext) bool {
	code := "53c56b6c766b00527ac4616c766b00c361681a53797374656d2e426c6f636b636861696e2e476574426c6f636b6c766b51527ac461681953797374656d2e53746f726167652e476574436f6e7465787404686173686c766b51c361681553797374656d2e4865616465722e47657448617368615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e7465787405696e6465786c766b51c361681653797374656d2e4865616465722e476574496e646578615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e74657874086d65726b526f6f746c766b51c361681d4f6e746f6c6f67792e4865616465722e4765744d65726b6c65526f6f74615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e746578740d6e657874436f6e73656e7375736c766b51c36168204f6e746f6c6f67792e4865616465722e4765744e657874436f6e73656e737573615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e746578740870726576486173686c766b51c361681953797374656d2e4865616465722e4765745072657648617368615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e746578740974696d655374616d706c766b51c361681a53797374656d2e4865616465722e47657454696d657374616d70615272681253797374656d2e53746f726167652e5075746161681953797374656d2e53746f726167652e476574436f6e746578740776657273696f6e6c766b51c361681a4f6e746f6c6f67792e4865616465722e47657456657273696f6e615272681253797374656d2e53746f726167652e50757461006c766b52527ac46203006c766b52c3616c7566"

	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestGetBlock - GetDefaultAccount error:%s", err)
		return false
	}

	_, err = ctx.Ont.NeoVM.DeployNeoVMSmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		true,
		code,
		"TestGetBlock",
		"1.0",
		"",
		"",
		"")

	if err != nil {
		ctx.LogError("TestGetBlock - DeploySmartContract error:%s", err)
		return false
	}

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetBlock - WaitForGenerateBlock error: %s", err)
		return false
	}

	height, err := ctx.Ont.GetCurrentBlockHeight()
	if err != nil {
		ctx.LogError("TestGetBlock - GetCurrentBlockHeight error: %s", err)
		return false
	}

	height -= 1
	block, err := ctx.Ont.GetBlockByHeight(height)
	if err != nil {
		ctx.LogError("TestGetBlock GetBlockByHeight error: %s", err)
		return false
	}

	header := block.Header
	codeHash, _ := utils.GetContractAddress(code)
	_, err = ctx.Ont.NeoVM.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeHash,
		[]interface{}{int(height)})

	if err != nil {
		ctx.LogError("TestGetBlock - InvokeNeoVMSmartContract error: %s", err)
		return false
	}

	_, err = ctx.Ont.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestGetBlock WaitForGenerateBlock error: %s", err)
		return false
	}

	if err != nil {
		ctx.LogError("TestGetBlock - InvokeSmartContract error: %s", err)
		return false
	}

	vmHash, err := ctx.Ont.GetStorage(codeHash.ToHexString(), []byte("hash"))
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage1111 error: %s", err)
		return false
	}

	hash := header.Hash()
	err = ctx.AssertToByteArray(vmHash, hash.ToArray())
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage2222 error: %s", err)
		return false
	}

	index, err := ctx.Ont.GetStorage(codeHash.ToHexString(), []byte("index"))
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage4 error: %s", err)
		return false
	}

	bValue, _ := vtypes.NewByteArray(index).GetBigInteger()
	err = ctx.AssertToInt(bValue, int(header.Height))
	if err != nil {
		ctx.LogError("TestGetBlock Height AssertToInt error: %s", err)
		return false
	}

	return true

	merker, err := ctx.Ont.GetStorage(codeHash.ToHexString(), []byte("merkRoot"))
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToByteArray(merker, header.TransactionsRoot.ToArray())
	if err != nil {
		ctx.LogError("TestGetBlock TransactionsRoot AssertToByteArray error: %s", err)
		return false
	}

	bookkeeper, err := ctx.Ont.GetStorage(codeHash.ToHexString(), []byte("nextConsensus"))
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToByteArray(bookkeeper, header.NextBookkeeper[:])
	if err != nil {
		ctx.LogError("TestGetBlock NextBookKeeper AssertToByteArray error:%s", err)
		return false
	}

	prevHash, err := ctx.Ont.GetStorage(codeHash.ToHexString(), []byte("prevHash"))
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToByteArray(prevHash, header.PrevBlockHash.ToArray())
	if err != nil {
		ctx.LogError("TestGetBlock PrevBlockHash AssertToByteArray error: %s", err)
		return false
	}

	timeStamp, err := ctx.Ont.GetStorage(codeHash.ToHexString(), []byte("timeStamp"))
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToInt(timeStamp, int(header.Timestamp))
	if err != nil {
		ctx.LogError("TestGetBlock Timestamp AssertToInt error:%s", err)
		return false
	}

	version, err := ctx.Ont.GetStorage(codeHash.ToHexString(), []byte("version"))
	if err != nil {
		ctx.LogError("TestGetBlock - GetStorage error: %s", err)
		return false
	}

	err = ctx.AssertToInt(version, int(header.Version))
	if err != nil {
		ctx.LogError("TestGetBlock Version AssertToInt error:%s", err)
		return false
	}

	return true
}
