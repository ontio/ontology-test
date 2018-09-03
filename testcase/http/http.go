package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common/constants"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/consensus/vbft"
	"github.com/ontio/ontology/consensus/vbft/config"
	"github.com/ontio/ontology/core/payload"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
	nvutils "github.com/ontio/ontology/smartcontract/service/native/utils"
)

func TestGetBlockByHeight(ctx *testframework.TestFrameworkContext) bool {
	curBlockHeight, err := ctx.Ont.GetCurrentBlockHeight()
	if err != nil {
		ctx.LogError("TestGetBlockByHeight GetBlockCount error:%s", err)
		return false
	}
	_, err = ctx.Ont.GetBlockByHeight(curBlockHeight )
	if err != nil {
		ctx.LogError("ctx.Ont.GetBlockByHeight error:%s", err)
		return false
	}
	return true
}

func TestGetBlockByHash(ctx *testframework.TestFrameworkContext) bool {
	curBlockHeight, err := ctx.Ont.GetCurrentBlockHeight()
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockCount error:%s", err)
		return false
	}

	blockHash, err := ctx.Ont.GetBlockHash(curBlockHeight)
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockHash error:%s", err)
		return false
	}
	block, err := ctx.Ont.GetBlockByHash(blockHash.ToHexString())
	if err != nil {
		ctx.LogError("ctx.Ont.GetBlockByHash error:%s", err)
		return false
	}
	bHash := block.Hash()
	if bHash != blockHash {
		ctx.LogError("TestGetBlockByHash block hash %s != %s", blockHash.ToHexString(), bHash.ToHexString())
		return false
	}
	return true
}

func TestGetCurrentBlockHeight(ctx *testframework.TestFrameworkContext) bool {
	num, err := ctx.Ont.GetCurrentBlockHeight()
	if err != nil {
		ctx.LogError("ctx.Ont.GetCurrentBlockHeight error:%s", err)
		return false
	}
	ctx.LogInfo("CurrentBlockHeight:", num)
	return true
}

func TestGetBlockHash(ctx *testframework.TestFrameworkContext) bool {
	blockCount, err := ctx.Ont.GetCurrentBlockHeight()
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockCount error:%s", err)
		return false
	}

	blockHash, err := ctx.Ont.GetBlockHash(blockCount)
	if err != nil {
		ctx.LogError("TestGetBlockByHash GetBlockHash error:%s", err)
		return false
	}
	ctx.LogInfo("blkhash:%s", blockHash.ToHexString())
	return true
}

func TestGetCurrentBlockHash(ctx *testframework.TestFrameworkContext) bool {
	blkhash, err := ctx.Ont.GetCurrentBlockHash()
	if err != nil {
		ctx.LogError("ctx.Ont.GetCurrentBlockHash error:%s", err)
		return false
	}
	ctx.LogInfo("TestGetCurrentBlockHash blkhash:%s", blkhash.ToHexString())
	return true
}

func TestGetRawTransaction(ctx *testframework.TestFrameworkContext) bool {
	block, err := ctx.Ont.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("TestGetRawTransaction GetBlockByHeight error:%s", err)
		return false
	}
	txBaseHash := block.Transactions[0].Hash()
	tx, err := ctx.Ont.GetTransaction(txBaseHash.ToHexString())
	if err != nil {
		ctx.LogError("ctx.Ont.GetRawTransaction error:%s", err)
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
	block, err := ctx.Ont.GetBlockByHeight(0)
	if err != nil {
		ctx.LogError("GetBlockByHeight error:%s", err)
		return false
	}
	//The first transaction is ont deploy transaction
	ont := block.Transactions[0]
	payload := ont.Payload.(*payload.DeployCode)

	contractAddress := types.AddressFromVmCode(payload.Code)
	contract, err := ctx.Ont.GetSmartContract(contractAddress.ToHexString())
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
	events, err := ctx.Ont.GetSmartContractEventByBlock(0)
	if err != nil {
		ctx.LogError("TestGetSmartContractEvent GetSmartContractEventByBlock error:%s", err)
		return false
	}

	scEvt, err := ctx.Ont.GetSmartContractEvent(events[0].TxHash)
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
	value, err := ctx.Ont.GetStorage(nvutils.OntContractAddress.ToHexString(), []byte(ont.TOTALSUPPLY_NAME))
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

func TestGetVbftInfo(ctx *testframework.TestFrameworkContext) bool {
	blkNum, err := ctx.Ont.GetCurrentBlockHeight()
	if err != nil {
		ctx.LogError("TestGetVbftInfo GetBlockCount error:%s", err)
		return false
	}
	blk, err := ctx.Ont.GetBlockByHeight(blkNum)
	if err != nil {
		ctx.LogError("TestGetVbftInfo GetBlockByHeight error:%s", err)
		return false
	}
	block, err := initVbftBlock(blk)
	if err != nil {
		ctx.LogError("TestGetVbftInfo initVbftBlock error:%s", err)
		return false
	}

	var cfg vconfig.ChainConfig
	if block.Info.NewChainConfig != nil {
		cfg = *block.Info.NewChainConfig
	} else {
		var cfgBlock *types.Block
		if block.Info.LastConfigBlockNum != math.MaxUint32 {
			cfgBlock, err = ctx.Ont.GetBlockByHeight(block.Info.LastConfigBlockNum)
			if err != nil {
				ctx.LogError("TestGetVbftInfo chainconfig GetBlockByHeight error:%s", err)
				return false
			}
		}
		blk, err := initVbftBlock(cfgBlock)
		if err != nil {
			ctx.LogError("TestGetVbftInfo initVbftBlock error:%s", err)
			return false
		}
		if blk.Info.NewChainConfig == nil {
			ctx.LogError("TestGetVbftInfo newchainconfig error:%s", err)
			return false
		}
		cfg = *blk.Info.NewChainConfig
	}
	fmt.Printf("block vbft chainConfig, View:%d, N:%d, C:%d, BlockMsgDelay:%v, HashMsgDelay:%v, PeerHandshakeTimeout:%v, MaxBlockChangeView:%d, PosTable:%v\n",
		cfg.View, cfg.N, cfg.C, cfg.BlockMsgDelay, cfg.HashMsgDelay, cfg.PeerHandshakeTimeout, cfg.MaxBlockChangeView, cfg.PosTable)
	for _, p := range cfg.Peers {
		fmt.Printf("peerInfo Index: %d, ID:%v\n", p.Index, p.ID)
	}
	return true
}

func initVbftBlock(block *types.Block) (*vbft.Block, error) {
	if block == nil {
		return nil, fmt.Errorf("nil block in initVbftBlock")
	}

	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(block.Header.ConsensusPayload, blkInfo); err != nil {
		return nil, fmt.Errorf("unmarshal blockInfo: %s", err)
	}
	return &vbft.Block{
		Block: block,
		Info:  blkInfo,
	}, nil
}
