package wasmvm

import "github.com/ontio/ontology-test/testframework"

func TestWasmVM() {
	testframework.TFramework.RegTestCase("TestWasmJsonContract", TestWasmJsonContract)
	testframework.TFramework.RegTestCase("TestWasmRawContract", TestWasmRawContract)

	testframework.TFramework.RegTestCase("TestCallWasmJsonContract", TestCallWasmJsonContract)

	//ICO-Test
	testframework.TFramework.RegTestCase("TestCallNativeContract", TestCallNativeContract)
	testframework.TFramework.RegTestCase("TestCallICOContract", TestICOContract)
	testframework.TFramework.RegTestCase("TestICOContractCollect", TestICOContractCollect)
	////domain-test
	testframework.TFramework.RegTestCase("TestDomainContract", TestDomainContract)
	testframework.TFramework.RegTestCase("TestDomainContract_invoke", TestDomainContract_Invoke)
	testframework.TFramework.RegTestCase("TestDomainContract_invoke2", TestDomainContract_Invoke2)
	testframework.TFramework.RegTestCase("TestDomainContract_invoke3", TestDomainContract_Invoke3)

	//call neovm test
	testframework.TFramework.RegTestCase("TestCallNeoContract", TestCallNeoContract)
}
