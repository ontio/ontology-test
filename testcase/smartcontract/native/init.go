package native

import "github.com/ontio/ontology-test/testframework"

func TestNative() {
	testframework.TFramework.RegTestCase("TestOntTransfer", TestOntTransfer)
}
