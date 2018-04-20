package native

import (
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology-test/testcase/smartcontract/native/nativeOracle"
)

func TestNative() {
	testframework.TFramework.RegTestCase("TestOntTransfer", TestOntTransfer)
	nativeOracle.TestNativeOracle()
}
