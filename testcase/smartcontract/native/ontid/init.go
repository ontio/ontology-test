package ontid

import "github.com/ontio/ontology-test/testframework"

func TestNativeOntID() {
	testframework.TFramework.RegTestCase("ontid", TestID)
	testframework.TFramework.RegTestCase("ontid-attr", TestAttr)
}
