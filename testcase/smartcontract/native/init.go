package native

import (
	"github.com/ontio/ontology-test/testcase/smartcontract/native/ontid"
	"github.com/ontio/ontology-test/testframework"
)

func TestNative() {
	testframework.TFramework.RegTestCase("TestOntTransfer", TestOntTransfer)
	ontid.TestNativeOntID()
}
