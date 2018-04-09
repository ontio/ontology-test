package operator

import "github.com/ontio/ontology-test/testframework"

func TestNeoVMOperator() {
	testframework.TFramework.RegTestCase("TestOperationAdd", TestOperationAdd)
	testframework.TFramework.RegTestCase("TestOperationSub", TestOperationSub)
	testframework.TFramework.RegTestCase("TestOperationMulti", TestOperationMulti)
	testframework.TFramework.RegTestCase("TestOperationDivide", TestOperationDivide)
	testframework.TFramework.RegTestCase("TestOperationSelfAdd", TestOperationSelfAdd)
	testframework.TFramework.RegTestCase("TestOperationSelfSub", TestOperationSelfSub)
	testframework.TFramework.RegTestCase("TestOperationLarger", TestOperationLarger)
	testframework.TFramework.RegTestCase("TestOperationLargerEqual", TestOperationLargerEqual)
	testframework.TFramework.RegTestCase("TestOperationSmaller", TestOperationSmaller)
	testframework.TFramework.RegTestCase("TestOperationSmallerEqual", TestOperationSmallerEqual)
	testframework.TFramework.RegTestCase("TestOperationEqual", TestOperationEqual)
	testframework.TFramework.RegTestCase("TestOperationNotEqual", TestOperationNotEqual)
	testframework.TFramework.RegTestCase("TestOperationNegative", TestOperationNegative)
	testframework.TFramework.RegTestCase("TestOperationOr", TestOperationOr)
	testframework.TFramework.RegTestCase("TestOperationAnd", TestOperationAnd)
	testframework.TFramework.RegTestCase("TestOperationLeftShift", TestOperationLeftShift)
	testframework.TFramework.RegTestCase("TestOperationRightShift", TestOperationRightShift)
}
