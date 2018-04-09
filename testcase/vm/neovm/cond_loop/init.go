package cond_loop

import (
	 "github.com/ontio/ontology-test/testframework"
)

func TestCondLoop(){
	testframework.TFramework.RegTestCase("TestIfElse", TestIfElse)
	testframework.TFramework.RegTestCase("TestSwitch", TestSwitch)
	testframework.TFramework.RegTestCase("TestFor", TestFor)
}