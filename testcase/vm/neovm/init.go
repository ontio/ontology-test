package neovm

import (
	"github.com/ontio/ontology-test/testcase/vm/neovm/call"
	"github.com/ontio/ontology-test/testcase/vm/neovm/cond_loop"
	"github.com/ontio/ontology-test/testcase/vm/neovm/datatype"
	"github.com/ontio/ontology-test/testcase/vm/neovm/operator"
)

func TestNeoVM() {
	operator.TestNeoVMOperator()
	datatype.TestDataType()
	cond_loop.TestCondLoop()
	call.TestCall()
}
