package utils

import (
	"github.com/ontio/ontology-test/testframework"
)

func TestUtils() {
	testframework.TFramework.RegTestCase("TestAsBigInteger", TestAsBigInteger)
	testframework.TFramework.RegTestCase("TestAsByteArrayBigInteger", TestAsByteArrayBigInteger)
	testframework.TFramework.RegTestCase("TestAsByteArrayString", TestAsByteArrayString)
	testframework.TFramework.RegTestCase("TestAsString", TestAsString)
	testframework.TFramework.RegTestCase("TestRange", TestRange)
	testframework.TFramework.RegTestCase("TestTake", TestTake)
}
