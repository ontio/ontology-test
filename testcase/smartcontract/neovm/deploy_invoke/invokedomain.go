package deploy_invoke

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

/*
python smart contract

from boa.blockchain.vm.Neo.Runtime import Log, Notify
from boa.blockchain.vm.System.ExecutionEngine import GetScriptContainer, GetExecutingScriptHash
from boa.blockchain.vm.Neo.Transaction import *
from boa.blockchain.vm.Neo.Blockchain import GetHeight, GetHeader
from boa.blockchain.vm.Neo.Action import RegisterAction
from boa.blockchain.vm.Neo.Runtime import GetTrigger, CheckWitness
from boa.blockchain.vm.Neo.TriggerType import Application, Verification
from boa.blockchain.vm.Neo.Output import GetScriptHash, GetValue, GetAssetId
from boa.blockchain.vm.Neo.Storage import GetContext, Get, Put, Delete
from boa.blockchain.vm.Neo.Header import GetTimestamp, GetNextConsensus

Push = RegisterAction('transfer')

def Main(operation, args):
    if operation == 'Query':
        domain = args[0]
        return Query(domain)

    if operation == 'Register':
        domain = args[0]
        owner = args[1]
        return Register(domain, owner)

    if operation == 'Transfer':
        domain = args[0]
        to = args[1]
        return Transfer(domain, to)

    if operation == 'Delete':
        domain = args[0]
        return Delete(domain)

    return False


def Query(domain):
    context = GetContext()
    owner = Get(context, domain);

    if owner != None:
        return False

    return owner

def Register(domain, owner):
    context = GetContext()
    occupy = Get(context, domain);
    if occupy != None:
        return False;
    Put(context, domain, owner)
    Push('hello')
    return True

def  Transfer(domain, to):
    if to == None:
        return False

    context = GetContext()
    owner = Get(context, domain)
    if owner == None:
        return False
    if owner == to:
        return True

    is_owner = CheckWitness(owner)

    if not is_owner:
        return False

    Put(context, domain, to)

    return True

def  Delete(domain):
    context = GetContext()
    occupy = Get(context, domain);
    if occupy != None:
        return False;
    # Put(context, domain, owner)
    Push('hello')

    return True


code = 0111c56b6c766b00527ac46c766b51527ac46c766b00c30551756572799c6421006c766b51c300c36c766b52527ac46203006c766b52c361650503616c75666c766b00c30852656769737465729c6435006c766b51c300c36c766b52527ac46c766b51c351c36c766b53527ac46203006c766b52c36c766b53c37c61650602616c75666c766b00c3085472616e736665729c6435006c766b51c300c36c766b52527ac46c766b51c351c36c766b54527ac46203006c766b52c36c766b54c37c6165cd00616c75666c766b00c30644656c6574659c6421006c766b51c300c36c766b52527ac46203006c766b52c361650f00616c756662030000616c756657c56b6c766b00527ac46168164e656f2e53746f726167652e476574436f6e74657874616c766b51527ac46c766b51c36c766b00c37c61680f4e656f2e53746f726167652e476574616c766b52527ac46c766b52c3009e640b0062030000616c75660568656c6c6f61087472616e7366657251c168124e656f2e52756e74696d652e4e6f7469667962030051616c75665fc56b6c766b00527ac46c766b51527ac46c766b51c3009c640c0062030000616c75666168164e656f2e53746f726167652e476574436f6e74657874616c766b52527ac46c766b52c36c766b00c37c61680f4e656f2e53746f726167652e476574616c766b53527ac46c766b53c3009c640b0062030000616c75666c766b53c36c766b51c39c640b0062030051616c75666c766b53c36168184e656f2e52756e74696d652e436865636b5769746e657373616c766b54527ac46c766b54c3630b0062030000616c75666c766b52c36c766b00c36c766b51c3527261680f4e656f2e53746f726167652e5075746162030051616c756659c56b6c766b00527ac46c766b51527ac46168164e656f2e53746f726167652e476574436f6e74657874616c766b52527ac46c766b52c36c766b00c37c61680f4e656f2e53746f726167652e476574616c766b53527ac40568656c6c6f61087472616e7366657251c168124e656f2e52756e74696d652e4e6f746966796c766b53c3009e640b0062030000616c75666c766b52c36c766b00c36c766b51c3527261680f4e656f2e53746f726167652e5075746162030051616c756656c56b6c766b00527ac46168164e656f2e53746f726167652e476574436f6e74657874616c766b51527ac46c766b51c36c766b00c37c61680f4e656f2e53746f726167652e476574616c766b52527ac46c766b52c3009e640b0062030000616c75666203006c766b52c3616c756652c56b6c766b00527ac46203006c766b00c361681b4e656f2e4865616465722e4765744e657874436f6e73656e73757361616c756652c56b6c766b00527ac46203006c766b00c361650002616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756653c56b6c766b00527ac46c766b51527ac462030000616c756654c56b6c766b00527ac46c766b51527ac46c766b52527ac462030000616c756653c56b6c766b00527ac46c766b51527ac462030000616c756651c56b62030000616c756652c56b6c766b00527ac46203006c766b00c36168184e656f2e4f75747075742e4765745363726970744861736861616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756651c56b6203000110616c756651c56b6203000100616c756652c56b6c766b00527ac462030000616c756651c56b62030000616c756653c56b6c766b00527ac46c766b51527ac462030000616c756652c56b6c766b00527ac462030000616c756651c56b62030000616c756652c56b6c766b00527ac46203006c766b00c361681f4e656f2e5472616e73616374696f6e2e476574556e7370656e74436f696e7361616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c756651c56b62030000616c756651c56b62030000616c756652c56b6c766b00527ac462030000616c756652c56b6c766b00527ac462030000616c7566
*/

func TestDomainSmartContract(ctx *testframework.TestFrameworkContext) bool {
	code := "0113c56b6a00527ac46a51527ac46a00c30551756572799c6416006a51c300c36a52527ac46a52c3659c026c7566616a00c30852656769737465729c6424006a51c300c36a52527ac46a51c351c36a53527ac46a52c36a53c37c65c6016c7566616a00c3085472616e736665729c6424006a51c300c36a52527ac46a51c351c36a54527ac46a52c36a54c37c65a9006c7566616a00c30644656c6574659c6416006a51c300c36a52527ac46a52c3650b006c756661006c756659c56b6a00527ac468164e656f2e53746f726167652e476574436f6e74657874616a51527ac46a51c36a00c37c680f4e656f2e53746f726167652e476574616a52527ac46a52c3009e640700006c7566610644656c6574656a00c37c056576656e7453c168124e656f2e52756e74696d652e4e6f74696679516c75660112c56b6a00527ac46a51527ac46a51c3009c640700006c75666168164e656f2e53746f726167652e476574436f6e74657874616a52527ac46a52c36a00c37c680f4e656f2e53746f726167652e476574616a53527ac46a53c3009c640700006c7566616a53c36a51c39c640700516c7566616a53c368184e656f2e52756e74696d652e436865636b5769746e657373616a54527ac46a54c3630700006c7566616a52c36a00c36a51c35272680f4e656f2e53746f726167652e50757461085472616e736665726a00c37c056576656e7453c168124e656f2e52756e74696d652e4e6f74696679516c75665bc56b6a00527ac46a51527ac468164e656f2e53746f726167652e476574436f6e74657874616a52527ac46a52c36a00c37c680f4e656f2e53746f726167652e476574616a53527ac46a53c3009e640700006c7566616a52c36a00c36a51c35272680f4e656f2e53746f726167652e507574610852656769737465726a00c36a51c35272056576656e7453c168124e656f2e52756e74696d652e4e6f74696679516c756659c56b6a00527ac468164e656f2e53746f726167652e476574436f6e74657874616a51527ac46a51c36a00c37c680f4e656f2e53746f726167652e476574616a52527ac40571756572796a00c37c056576656e7453c168124e656f2e52756e74696d652e4e6f746966796a52c3009e640700006c7566616a52c36c7566"
	codeAddress, _ := utils.GetContractAddress(code)

	ctx.LogInfo("=====CodeAddress===%s", codeAddress.ToHexString())
	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestDomainSmartContract GetDefaultAccount error:%s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		true,
		code,
		"TestDomainSmartContract",
		"1.0",
		"",
		"",
		"",
	)

	if err != nil {
		ctx.LogError("TestDomainSmartContract DeploySmartContract error: %s", err)
	}

	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestDomainSmartContract WaitForGenerateBlock error: %s", err)
		return false
	}

	_, err = ctx.Ont.Rpc.InvokeNeoVMContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		codeAddress,
		[]interface{}{"Register", []interface{}{[]byte("ont.io"), []byte("onchain")}})
	if err != nil {
		ctx.LogError("TestDomainSmartContract InvokeNeoVMSmartContract error: %s", err)
	}

	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestDomainSmartContract WaitForGenerateBlock error: %s", err)
		return false
	}

	svalue, err := ctx.Ont.Rpc.GetStorage(codeAddress, []byte("ont.io"))
	if err != nil {
		ctx.LogError("TestDomainSmartContract GetStorageItem key:hello error: %s", err)
		return false
	}

	ctx.LogInfo("==svalue = %v", string(svalue))

	err = ctx.AssertToString(string(svalue), "onchain")
	if err != nil {
		ctx.LogError("TestDomainSmartContract AssertToString error: %s", err)
		return false
	}

	return true
}
