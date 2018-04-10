package deploy_invoke

import (
	"github.com/ontio/ontology-test/testframework"
)

func TestDeployInvoke() {
	testframework.TFramework.RegTestCase("TestDeploySmartContract", TestDeploySmartContract)
	testframework.TFramework.RegTestCase("TestInvokeSmartContract", TestInvokeSmartContract)
}
