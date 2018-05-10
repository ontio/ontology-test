package call

import (
	"github.com/ontio/ontology-test/testframework"
)

func TestCall() {
	testframework.TFramework.RegTestCase("TestCallContractStatic", TestCallContractStatic)
}
