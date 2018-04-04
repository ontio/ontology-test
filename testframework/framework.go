/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

//light test framework for ontology
package testframework

import (
	"fmt"
	log4 "github.com/alecthomas/log4go"
	sdk "github.com/ontio/ontology-go-sdk"
	wa "github.com/ontio/ontology-go-sdk/wallet"
	"reflect"
	"time"
)

//Default TestFramework instance
var TFramework = NewTestFramework()

//TestCase type
type TestCase func(ctx *TestFrameworkContext) bool

//TestFramework manage test case and run test case
type TestFramework struct {
	//Test case start time
	startTime time.Time
	//hold the test case for testing
	testCases []TestCase
	//Map the test case name to test case
	testCasesMap map[string]string
	//hold the test case result for testing
	testCaseRes map[string]bool
	//OntologySdk object
	ont *sdk.OntologySdk
	//OntWallet object
	wallet *wa.OntWallet
}

//NewTestFramework return a TestFramework instance
func NewTestFramework() *TestFramework {
	return &TestFramework{
		testCases:    make([]TestCase, 0),
		testCasesMap: make(map[string]string, 0),
		testCaseRes:  make(map[string]bool, 0),
	}
}

//RegTestCase register a test case to framework
func (this *TestFramework) RegTestCase(name string, testCase TestCase) {
	this.testCases = append(this.testCases, testCase)
	this.testCasesMap[this.getTestCaseId(testCase)] = name
}

//Start run test case
func (this *TestFramework) Start() {
	this.runTestList(this.testCases)
}

func (this *TestFramework) runTestList(testCaseList []TestCase) {
	this.onTestStart()
	defer this.onTestFinish()
	failNowCh := make(chan interface{}, 0)
	for i, testCase := range this.testCases {
		select {
		case <-failNowCh:
			this.onTestFailNow()
			return
		default:
			this.runTest(i+1, failNowCh, testCase)
		}
	}
}

//Run a single test case
func (this *TestFramework) runTest(index int, failNowCh chan interface{}, testCase TestCase) {
	ctx := NewTestFrameworkContext(this.ont, this.wallet, failNowCh)
	this.onBeforeTestCaseStart(index, testCase)
	ok := testCase(ctx)
	this.onAfterTestCaseFinish(index, testCase, ok)
	this.testCaseRes[this.getTestCaseId(testCase)] = ok
}

//SetOntSdk ontology sdk instance to test framework
func (this *TestFramework) SetOntSdk(ont *sdk.OntologySdk) {
	this.ont = ont
}

//SetWallet wallet instance to test framework
func (this *TestFramework) SetWallet(wallet *wa.OntWallet) {
	this.wallet = wallet
}

//onTestStart invoke at the beginning of test
func (this *TestFramework) onTestStart() {
	version, _ := this.ont.Rpc.GetVersion()
	log4.Info("\t\t\t===============================================================")
	log4.Info("\t\t\t-------Ontology Test Start Version:%d", version)
	log4.Info("\t\t\t===============================================================")
	log4.Info("")
	this.startTime = time.Now()
}

//onTestStart invoke at the end of test
func (this *TestFramework) onTestFinish() {
	failedList := make([]string, 0)
	successList := make([]string, 0)
	for testCase, ok := range this.testCaseRes {
		if ok {
			successList = append(successList, this.getTestCaseName(testCase))
		} else {
			failedList = append(failedList, this.getTestCaseName(testCase))
		}
	}

	skipList := make([]string, 0)
	for _, testCase := range this.testCases {
		_, ok := this.testCaseRes[this.getTestCaseId(testCase)]
		if !ok {
			skipList = append(skipList, this.getTestCaseName(testCase))
		}
	}

	succCount := len(successList)
	failedCount := len(failedList)

	log4.Info("\t\t===============================================================")
	log4.Info("\t\tOntology Test Finish Total:%v Success:%v Failed:%v Skip:%v TimeCost:%.2f s.",
		len(this.testCases),
		succCount,
		failedCount,
		len(this.testCases)-succCount-failedCount,
		time.Now().Sub(this.startTime).Seconds())
	if succCount > 0 {
		log4.Info("\t\t---------------------------------------------------------------")
		log4.Info("\t\t\tSuccess list:")
		for i, succCase := range successList {
			log4.Info("\t\t\t%d.\t%s", i+1, succCase)
		}
	}
	if failedCount > 0 {
		log4.Info("\t\t---------------------------------------------------------------")
		log4.Info("\t\t\tFail list:")
		for i, failCase := range failedList {
			log4.Info("\t\t\t%d.\t%s", i+1, failCase)
		}
	}
	if len(skipList) > 0 {
		log4.Info("\t\t---------------------------------------------------------------")
		log4.Info("\t\t\tSkip list:")
		for i, failCase := range skipList {
			log4.Info("\t\t\t%d.\t%s", i+1, failCase)
		}
	}
	log4.Info("\t\t===============================================================")
}

//onTestFailNow invoke when context.FailNow() was be called
func (this *TestFramework) onTestFailNow() {
	log4.Info("Test Stop.")
}

//onBeforeTestCaseStart invoke before single test case
func (this *TestFramework) onBeforeTestCaseStart(index int, testCase TestCase) {
	log4.Info("===============================================================")
	log4.Info("%d. Start TestCase:%s", index, this.getTestCaseName(testCase))
	log4.Info("---------------------------------------------------------------")
}

//onBeforeTestCaseStart invoke after single test case
func (this *TestFramework) onAfterTestCaseFinish(index int, testCase TestCase, res bool) {
	if res {
		log4.Info("TestCase:%s success.", this.getTestCaseName(testCase))
	} else {
		log4.Info("TestCase:%s failed.", this.getTestCaseName(testCase))
	}
	log4.Info("---------------------------------------------------------------")
	log4.Info("")
}

//getTestCaseName return the name of test case
func (this *TestFramework) getTestCaseName(testCase interface{}) string {
	testCaseStr, ok := testCase.(string)
	if !ok {
		testCaseStr = this.getTestCaseId(testCase)
	}
	name, ok := this.testCasesMap[testCaseStr]
	if ok {
		return name
	}
	return ""
}

//getTestCaseId return the id of test case
func (this *TestFramework) getTestCaseId(testCase interface{}) string {
	return fmt.Sprintf("%v", reflect.ValueOf(testCase).Pointer())
}
