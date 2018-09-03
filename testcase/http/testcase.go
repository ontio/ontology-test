package http

import "github.com/ontio/ontology-test/testframework"

func TestHttp() {
	testframework.TFramework.RegTestCase("TestGetBlockByHeight", TestGetBlockByHeight)
	testframework.TFramework.RegTestCase("TestGetBlockByHash", TestGetBlockByHash)
	testframework.TFramework.RegTestCase("TestGetCurrentBlockHeight", TestGetCurrentBlockHeight)
	testframework.TFramework.RegTestCase("TestGetBlockHash", TestGetBlockHash)
	testframework.TFramework.RegTestCase("TestGetCurrentBlockHash", TestGetCurrentBlockHash)
	testframework.TFramework.RegTestCase("TestGetRawTransaction", TestGetRawTransaction)
	testframework.TFramework.RegTestCase("TestGetSmartContract", TestGetSmartContract)
	testframework.TFramework.RegTestCase("TestGetSmartContractEvent", TestGetSmartContractEvent)
	testframework.TFramework.RegTestCase("TestGetStorage", TestGetStorage)
	//testframework.TFramework.RegTestCase("TestGetVbftInfo", TestGetVbftInfo)
}
