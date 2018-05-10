package jsonrpc

import (
	"fmt"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
)

func TestGetBlockByHeight(ctx *testframework.TestFrameworkContext) bool {
	block, err := ctx.Ont.Rpc.GetBlockByHeight(10)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBlockByHeight error:%s", err)
		return false
	}
	fmt.Println(block)
	return true
}

func TestGetBlockByHash(ctx *testframework.TestFrameworkContext) bool {
	bs, err := common.HexToBytes("72fe5d25ff63b2cf33e697d79f9b4f310d6e919deeeeaabffd79c293ab50fef9")
	txhash, err := common.Uint256ParseFromBytes(bs)
	if err != nil {
		ctx.LogError("common.Uint256ParseFromBytes error:%s", err)
		return false
	}
	block, err := ctx.Ont.Rpc.GetBlockByHash(txhash)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBlockByHash error:%s", err)
		return false
	}
	fmt.Println(block)
	return true
}

func TestGetBalance(ctx *testframework.TestFrameworkContext) bool {
	address, err := common.AddressFromBase58("TA8ZD2Z9R25Y2uibc96FdBn3owhTqZZRy6")
	if err != nil {
		ctx.LogError("common.AddressFromBase58 error:%s", err)
		return false
	}
	balance, err := ctx.Ont.Rpc.GetBalance(address)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBalance error:%s", err)
		return false
	}
	fmt.Println(balance)
	return true
}

func TestGetBlockCount(ctx *testframework.TestFrameworkContext) bool {
	num, err := ctx.Ont.Rpc.GetBlockCount()
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBlockCount error:%s", err)
		return false
	}
	fmt.Println("num:", num)
	return true
}

func TestGetBlockHash(ctx *testframework.TestFrameworkContext) bool {
	blkhash, err := ctx.Ont.Rpc.GetBlockHash(1)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetBlockHash error:%s", err)
		return false
	}
	fmt.Println("blkhash:", blkhash)
	return true
}

func TestGetCurrentBlockHash(ctx *testframework.TestFrameworkContext) bool {
	blkhash, err := ctx.Ont.Rpc.GetCurrentBlockHash()
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetCurrentBlockHash error:%s", err)
		return false
	}
	fmt.Println("blkhash:", blkhash)
	return true
}

func TestGetRawTransaction(ctx *testframework.TestFrameworkContext) bool {
	txhashbs, err := common.HexToBytes("a6733c1aa5f0885062121ce98fdd647c8e1bf4cb6e0e618ee9c9cd19b4c79997")
	if err != nil {
		ctx.LogError("common.HexToBytes error:%s", err)
		return false
	}
	txhash, err := common.Uint256ParseFromBytes(txhashbs)
	if err != nil {
		ctx.LogError("common.Uint256ParseFromBytes error:%s", err)
		return false
	}
	tx, err := ctx.Ont.Rpc.GetRawTransaction(txhash)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetRawTransaction error:%s", err)
		return false
	}
	fmt.Println("transaction:", tx)
	return true
}

func TestGetSmartContract(ctx *testframework.TestFrameworkContext) bool {
	bs, err := common.HexToBytes("80b0cc71bda8653599c5666cae084bff587e2de1")
	addr, err := common.AddressParseFromBytes(bs)
	if err != nil {
		ctx.LogError("common.AddressFromBase58 error:%s", err)
		return false
	}
	code, err := ctx.Ont.Rpc.GetSmartContract(addr)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetSmartContract error:%s", err)
		return false
	}
	fmt.Println("code:", code)
	return true
}

func TestGetSmartContractEvent(ctx *testframework.TestFrameworkContext) bool {
	txbytes, err := common.HexToBytes("a6733c1aa5f0885062121ce98fdd647c8e1bf4cb6e0e618ee9c9cd19b4c79997")
	if err != nil {
		ctx.LogError("common.HexToBytes error:%s", err)
		return false
	}
	txhash, err := common.Uint256ParseFromBytes(txbytes)
	if err != nil {
		ctx.LogError("common.Uint256ParseFromBytes error:%s", err)
		return false
	}
	smartcodeevent, err := ctx.Ont.Rpc.GetSmartContractEvent(txhash)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetSmartContractEvent error:%s", err)
		return false
	}
	fmt.Println("smartcodeevent:", smartcodeevent)
	return true
}

func TestGetStorage(ctx *testframework.TestFrameworkContext) bool {
	codeaddress, err := common.HexToBytes("80b0cc71bda8653599c5666cae084bff587e2de1")
	addr, err := common.AddressParseFromBytes(codeaddress)
	if err != nil {
		ctx.LogError("common.AddressParseFromBytes error:%s", err)
		return false
	}
	keyBytes, err := common.HexToBytes("key")
	storage, err := ctx.Ont.Rpc.GetStorage(addr, keyBytes)
	if err != nil {
		ctx.LogError("ctx.Ont.Rpc.GetSmartContractEvent error:%s", err)
		return false
	}
	fmt.Println("storage:", storage)
	return true
}
