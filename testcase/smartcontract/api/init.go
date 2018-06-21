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

package api

import (
	"github.com/ontio/ontology-test/testcase/smartcontract/api/blockchain"
	"github.com/ontio/ontology-test/testcase/smartcontract/api/contract"
	"github.com/ontio/ontology-test/testcase/smartcontract/api/executionengine"
	"github.com/ontio/ontology-test/testcase/smartcontract/api/runtime"
	"github.com/ontio/ontology-test/testcase/smartcontract/api/storage"
	"github.com/ontio/ontology-test/testcase/smartcontract/api/transaction"
	"github.com/ontio/ontology-test/testframework"
)

func TestSmartContractApi() {
	testframework.TFramework.RegTestCase("TestGetBlock", blockchain.TestGetBlock)
	testframework.TFramework.RegTestCase("TestGetContract", contract.TestGetContract)
	testframework.TFramework.RegTestCase("TestContractCreate", contract.TestContractCreate)
	testframework.TFramework.RegTestCase("TestContractDestroy", contract.TestContractDestroy)
	testframework.TFramework.RegTestCase("TestCallingScriptHash", executionengine.TestCallingScriptHash)
	testframework.TFramework.RegTestCase("TestCheckWitness", runtime.TestCheckWitness)
	//testframework.TFramework.RegTestCase("TestRuntimLog", runtime.TestRuntimLog)
	testframework.TFramework.RegTestCase("TestRuntimeNotify", runtime.TestRuntimeNotify)
	testframework.TFramework.RegTestCase("TestGetTxHash", transaction.TestGetTxHash)
	testframework.TFramework.RegTestCase("TestGetTxType", transaction.TestGetTxType)
	testframework.TFramework.RegTestCase("TestStorage", storage.TestStorage)
}
