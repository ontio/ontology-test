package jsonrpc

import (
	"bytes"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common/constants"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	nvutils "github.com/ontio/ontology/smartcontract/service/native/utils"
)

func TestGetBlockByHeight(ctx *testframework.TestFrameworkContext) bool {
	blockCount, err := ctx.Ont.Rpc.GetBlockCount()
	if err != nil {
		ctx.LogError("TestGetBlockByHeight GetBlockCount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.GetBlockByHeight(blockCount - 1)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBlockByHeight error:%s", err)
		return false
	}
	return true
}

func TestGetBlockByHash(ctx *testframework.TestFrameworkContext) bool {
	blockCount, err := ctx.Ont.Rpc.GetBlockCount()
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockCount error:%s", err)
		return false
	}

	blockHash, err := ctx.Ont.Rpc.GetBlockHash(blockCount - 1)
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockHash error:%s", err)
		return false
	}
	block, err := ctx.Ont.Rpc.GetBlockByHash(blockHash)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBlockByHash error:%s", err)
		return false
	}
	bHash := block.Hash()
	if bHash != blockHash {
		ctx.LogError("TestGetBlockByHash block hash %s != %s", blockHash.ToHexString(), bHash.ToHexString())
		return false
	}
	return true
}

func TestGetBalance(ctx *testframework.TestFrameworkContext) bool {
	defAcc, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestGetBalance GetDefaultAccount error:%s", err)
		return false
	}
	balance, err := ctx.Ont.Rpc.GetBalance(defAcc.Address)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBalance error:%s", err)
		return false
	}
	ctx.LogInfo("%v", balance)
	return true
}

func TestGetBlockCount(ctx *testframework.TestFrameworkContext) bool {
	num, err := ctx.Ont.Rpc.GetBlockCount()
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBlockCount error:%s", err)
		return false
	}
	ctx.LogInfo("GetBlockCount:", num)
	return true
}

func TestGetBlockHash(ctx *testframework.TestFrameworkContext) bool {
	blockCount, err := ctx.Ont.Rpc.GetBlockCount()
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockCount error:%s", err)
		return false
	}

	blockHash, err := ctx.Ont.Rpc.GetBlockHash(blockCount - 1)
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockHash error:%s", err)
		return false
	}
	ctx.LogInfo("blkhash:%s", blockHash)
	return true
}

func TestGetCurrentBlockHash(ctx *testframework.TestFrameworkContext) bool {
	blkhash, err := ctx.Ont.Rpc.GetCurrentBlockHash()
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetCurrentBlockHash error:%s", err)
		return false
	}
	ctx.LogInfo("TestGetCurrentBlockHash blkhash:%s", blkhash.ToHexString())
	return true
}

func TestGetRawTransaction(ctx *testframework.TestFrameworkContext) bool {
	block, err := ctx.Ont.Rpc.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("TestGetRawTransaction GetBlockByHeight error:%s", err)
		return false
	}
	txBaseHash := block.Transactions[0].Hash()
	tx, err := ctx.Ont.Rpc.GetRawTransaction(txBaseHash)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetRawTransaction error:%s", err)
		return false
	}
	txHash := tx.Hash()
	if txHash != txBaseHash {
		ctx.LogError("TestGetRawTransaction %x != %x", txHash, txBaseHash)
		return false
	}
	return true
}

func TestGetSmartContract(ctx *testframework.TestFrameworkContext) bool {
	block, err := ctx.Ont.Rpc.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("GetBlockByHeight error:%s", err)
		return false
	}
	//The first transaction is ont deploy transaction
	ont := block.Transactions[0]
	payload := ont.Payload.(*payload.DeployCode)

	contractAddress := types.AddressFromVmCode(payload.Code)
	contract, err := ctx.Ont.Rpc.GetSmartContract(contractAddress)
	if err != nil {
		ctx.LogError("GetSmartContract error:%s", err)
		return false
	}
	ctx.LogInfo("TestGetSmartContract:")
	ctx.LogInfo("Code:%x", contract.Code)
	ctx.LogInfo("Author:%s", contract.Author)
	ctx.LogInfo("Version:%s", contract.Version)
	ctx.LogInfo("NeedStorage:%v", contract.NeedStorage)
	ctx.LogInfo("Email:%s", contract.Email)
	ctx.LogInfo("Description:%s", contract.Description)
	return true
}

func TestGetSmartContractEvent(ctx *testframework.TestFrameworkContext) bool {
	events, err := ctx.Ont.Rpc.GetSmartContractEventByBlock(0)
	if err != nil {
		ctx.LogError("TestGetSmartContractEvent GetSmartContractEventByBlock error:%s", err)
		return false
	}

	scEvt, err := ctx.Ont.Rpc.GetSmartContractEventWithHexString(events[0].TxHash)
	if err != nil {
		ctx.LogError("TestGetSmartContractEvent GetSmartContractEvent error:%s", err)
		return false
	}

	ctx.LogInfo(" TxHash:%s", scEvt.TxHash)
	ctx.LogInfo(" State:%d", scEvt.State)
	ctx.LogInfo(" GasConsumed:%d", scEvt.GasConsumed)
	for _, notify := range scEvt.Notify {
		ctx.LogInfo(" SmartContractAddress:%s", notify.ContractAddress)
		states := notify.States.([]interface{})
		name := states[0].(string)
		from := states[1].(string)
		to := states[2].(string)
		value := states[3].(float64)
		ctx.LogInfo(" State Name:%s from:%s to:%s value:%d", name, from, to, int(value))
	}
	return true
}

func TestGetStorage(ctx *testframework.TestFrameworkContext) bool {
	value, err := ctx.Ont.Rpc.GetStorage(nvutils.OntContractAddress, []byte(ont.TOTALSUPPLY_NAME))
	if err != nil {
		ctx.LogError("TestGetStorage error:%s", err)
		return false
	}
	if value == nil {
		ctx.LogError("TestGetStorage value is nil")
		return false
	}
	totalSupply, err := serialization.ReadUint64(bytes.NewReader(value))
	if err != nil {
		ctx.LogError("TestGetStorage serialization.ReadUint64 error:%s", err)
		return false
	}
	if totalSupply != constants.ONT_TOTAL_SUPPLY {
		ctx.LogError("TestGetStorage totalSupply %d != %d", totalSupply, constants.ONT_TOTAL_SUPPLY)
		return false
	}
	ctx.LogInfo("TestGetStorage %d\n", totalSupply)
	return true
}
