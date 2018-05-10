package jsonrpc

import "github.com/ontio/ontology-test/testframework"

func TestRpc() {
	testframework.TFramework.RegTestCase("TestGetBlockByHeight", TestGetBlockByHeight)
	testframework.TFramework.RegTestCase("TestGetBlockByHash", TestGetBlockByHash)
	testframework.TFramework.RegTestCase("TestGetBalance", TestGetBalance)
	testframework.TFramework.RegTestCase("TestGetBlockCount", TestGetBlockCount)
	testframework.TFramework.RegTestCase("TestGetBlockHash", TestGetBlockHash)
	testframework.TFramework.RegTestCase("TestGetCurrentBlockHash", TestGetCurrentBlockHash)
	testframework.TFramework.RegTestCase("TestGetRawTransaction", TestGetRawTransaction)
	testframework.TFramework.RegTestCase("TestGetSmartContract", TestGetSmartContract)
	testframework.TFramework.RegTestCase("TestGetSmartContractEvent", TestGetSmartContractEvent)
	testframework.TFramework.RegTestCase("TestGetStorage", TestGetStorage)
}
