package wasmvm

import (
	"github.com/ontio/ontology-test/testframework"
)

func TestWasmVM() {
	testframework.TFramework.RegTestCase("TestWasmJsonContract", TestWasmJsonContract)
	testframework.TFramework.RegTestCase("TestWasmRawContract", TestWasmRawContract)
}
