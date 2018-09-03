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
	//Map the test case id to test case name
	testCaseNameMap map[string]string
	//Map test case id to test case
	testCasesMap map[string]TestCase
	//Map the test case result for testing
	testCaseRes map[string]bool
	//OntologySdk object
	ont *sdk.OntologySdk
	//OntWallet object
	wallet *sdk.Wallet
	//Callback func before running test
	before func(ctx *TestFrameworkContext)
	//Callback func After running test
	after func(ctx *TestFrameworkContext)
}

//NewTestFramework return a TestFramework instance
func NewTestFramework() *TestFramework {
	return &TestFramework{
		testCases:       make([]TestCase, 0),
		testCaseNameMap: make(map[string]string, 0),
		testCasesMap:    make(map[string]TestCase, 0),
		testCaseRes:     make(map[string]bool, 0),
	}
}

//RegTestCase register a test case to framework
func (this *TestFramework) RegTestCase(name string, testCase TestCase) {
	this.testCases = append(this.testCases, testCase)
	testCaseId := this.getTestCaseId(testCase)
	this.testCaseNameMap[testCaseId] = name
	this.testCasesMap[testCaseId] = testCase
}

func (this *TestFramework) SetBeforeCallback(callback func(ctx *TestFrameworkContext)) {
	this.before = callback
}

func (this *TestFramework) SetAfterCallback(callback func(ctx *TestFrameworkContext)) {
	this.after = callback
}

//Start run test case
func (this *TestFramework) Start(testCases []string) {
	if len(testCases) > 0 {
		taseCaseList := make([]TestCase, 0, len(testCases))
		for _, t := range testCases {
			if t == "" {
				continue
			}
			testCase := this.getTestCaseByName(t)
			if testCase != nil {
				taseCaseList = append(taseCaseList, testCase)
			}
		}
		if len(taseCaseList) > 0 {
			this.runTestList(taseCaseList)
			return
		}
		log4.Info("Not test case to run")
		return
	}
	this.runTestList(this.testCases)
}

func (this *TestFramework) runTestList(testCaseList []TestCase) {
	this.onTestStart()
	defer this.onTestFinish(testCaseList)
	failNowCh := make(chan interface{}, 10)
	ctx := NewTestFrameworkContext(this.ont, this.wallet, failNowCh)
	if this.before != nil {
		this.before(ctx)
	}
	if this.after != nil {
		defer this.after(ctx)
	}
	for i, testCase := range testCaseList {
		select {
		case <-failNowCh:
			this.onTestFailNow()
			return
		default:
			this.runTest(i+1, ctx, testCase)
		}
	}
}

//Run a single test case
func (this *TestFramework) runTest(index int, ctx *TestFrameworkContext, testCase TestCase) {
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
func (this *TestFramework) SetWallet(wallet *sdk.Wallet) {
	this.wallet = wallet
}

//onTestStart invoke at the beginning of test
func (this *TestFramework) onTestStart() {
	version, _ := this.ont.GetVersion()
	log4.Info("===============================================================")
	log4.Info("-------Ontology Test Start Version:%s", version)
	log4.Info("===============================================================")
	log4.Info("")
	this.startTime = time.Now()
}

//onTestStart invoke at the end of test
func (this *TestFramework) onTestFinish(testCaseList []TestCase) {
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
	for _, testCase := range testCaseList {
		_, ok := this.testCaseRes[this.getTestCaseId(testCase)]
		if !ok {
			skipList = append(skipList, this.getTestCaseName(testCase))
		}
	}

	succCount := len(successList)
	failedCount := len(failedList)

	log4.Info("===============================================================")
	log4.Info("Ontology Test Finish Total:%v Success:%v Failed:%v Skip:%v TimeCost:%.2f s.",
		len(this.testCases),
		succCount,
		failedCount,
		len(this.testCases)-succCount-failedCount,
		time.Now().Sub(this.startTime).Seconds())
	if succCount > 0 {
		log4.Info("---------------------------------------------------------------")
		log4.Info("Success list:")
		for i, succCase := range successList {
			log4.Info("%d.\t%s", i+1, succCase)
		}
	}
	if failedCount > 0 {
		log4.Info("---------------------------------------------------------------")
		log4.Info("Fail list:")
		for i, failCase := range failedList {
			log4.Info("%d.\t%s", i+1, failCase)
		}
	}
	if len(skipList) > 0 {
		log4.Info("---------------------------------------------------------------")
		log4.Info("Skip list:")
		for i, failCase := range skipList {
			log4.Info("%d.\t%s", i+1, failCase)
		}
	}
	log4.Info("===============================================================")
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
	name, ok := this.testCaseNameMap[testCaseStr]
	if ok {
		return name
	}
	return ""
}

//getTestCaseByName return test case by test case name
func (this *TestFramework) getTestCaseByName(name string) TestCase {
	testCaseId := ""
	for id, n := range this.testCaseNameMap {
		if n == name {
			testCaseId = id
			break
		}
	}
	if testCaseId == "" {
		return nil
	}
	return this.testCasesMap[testCaseId]
}

//getTestCaseId return the id of test case
func (this *TestFramework) getTestCaseId(testCase interface{}) string {
	return fmt.Sprintf("%v", reflect.ValueOf(testCase).Pointer())
}
