package nativeOracle

import "github.com/ontio/ontology-test/testframework"

func TestNativeOracle() {
	testframework.TFramework.RegTestCase("TestOracle", TestOracle)
	testframework.TFramework.RegTestCase("TestCronOracle", TestCronOracle)
}
