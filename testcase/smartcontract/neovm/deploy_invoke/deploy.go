package deploy_invoke

import (
<<<<<<< HEAD
=======
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
>>>>>>> 86d51efafdf1683a867a80a58a2561eee4dc3a44
	"time"

	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/types"
)

func TestDeploySmartContract(ctx *testframework.TestFrameworkContext) bool {
	signer, err := ctx.Wallet.GetDefaultAccount()
	if err != nil {
		ctx.LogError("TestDeploySmartContract GetDefaultAccount error:%s", err)
		return false
	}
	_, err = ctx.Ont.Rpc.DeploySmartContract(signer,
		types.NEOVM,
		true,
		contractCode,
		"TestDeploySmartContract",
		"1.0",
		"",
		"",
		"",
	)

	ctx.LogInfo("CodeAddress:%x", contractCodeAddress)
	if err != nil {
		ctx.LogError("TestDeploySmartContract DeploySmartContract error:%s", err)
		return false
	}
	//WaitForGenerateBlock
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 1)
	if err != nil {
		ctx.LogError("TestDeploySmartContract WaitForGenerateBlock error:%s", err)
		return false
	}
	return true
}
