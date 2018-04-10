package datatype

import (
	"github.com/ontio/ontology-test/testframework"
)

func TestDataType() {
	testframework.TFramework.RegTestCase("TestBoolean", TestBoolean)
	testframework.TFramework.RegTestCase("TestInteger", TestInteger)
	testframework.TFramework.RegTestCase("TestString", TestString)
	testframework.TFramework.RegTestCase("TestArray", TestArray)
	testframework.TFramework.RegTestCase("TestReturnType", TestReturnType)
	testframework.TFramework.RegTestCase("TestByteArray", TestByteArray)
}
