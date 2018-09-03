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
package main

import (
	"flag"
	log4 "github.com/alecthomas/log4go"
	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-test/common"
	_ "github.com/ontio/ontology-test/testcase"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common/log"
	"math/rand"
	"strings"
	"time"
)

var (
	TestConfig string //Test config file
	LogConfig  string //Log config file
	TestCases  string //TestCase list in cmdline
)

func init() {
	flag.StringVar(&TestConfig, "cfg", "./config_test.json", "Config of ontology-test")
	flag.StringVar(&LogConfig, "lfg", "./log4go.xml", "Log config of ontology-test")
	flag.StringVar(&TestCases, "t", "", "Test case to run. use ',' to split test case")
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log4.LoadConfiguration(LogConfig)
	log.InitLog(1) //init log module in ontology
	defer time.Sleep(time.Second)

	err := common.DefConfig.Init(TestConfig)
	if err != nil {
		log4.Error("DefConfig.Init error:%s", err)
		return
	}

	ontSdk := sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(common.DefConfig.JsonRpcAddress)
	wallet, err := ontSdk.OpenWallet(common.DefConfig.WalletFile)
	if err != nil {
		log4.Error("OpenOrCreateWallet %s error:%s", common.DefConfig.WalletFile, err)
		return
	}
	testCases := make([]string, 0)
	if TestCases != "" {
		testCases = strings.Split(TestCases, ",")
	}
	testframework.TFramework.SetOntSdk(ontSdk)
	testframework.TFramework.SetWallet(wallet)
	//Start run test case
	testframework.TFramework.Start(testCases)
}
