package deploy_invoke

import (
	"time"

	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology-test/testframework"
)

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System;
using System.ComponentModel;
using System.Numerics;

namespace Neo.SmartContract
{
    public class Domain : Framework.SmartContract
    {
        public static object Main(string operation, params object[] args)
        {
            switch (operation)
            {
                case "Query":
                    return Query((string)args[0]);
                case "Register":
                    return Register((string)args[0], (byte[])args[1]);
                case "Transfer":
                    return Transfer((string)args[0], (byte[])args[1]);
                case "Delete":
                    return Delete((string)args[0]);
                default:
                    return false;
            }
        }

        public static byte[] Query(string domain)
        {
            return Storage.Get(Storage.CurrentContext, domain);
        }

        public static bool Register(string domain, byte[] owner)
        {
            byte[] value = Storage.Get(Storage.CurrentContext, domain);
            if (value != null) return false;
            Storage.Put(Storage.CurrentContext, domain, owner);
            return true;
        }

        public static bool Transfer(string domain, byte[] to)
        {
            byte[] from = Storage.Get(Storage.CurrentContext, domain);
            if (from == null) return false;
            if (!Runtime.CheckWitness(from)) return false;
            Storage.Put(Storage.CurrentContext, domain, to);
            return true;
        }

        public static bool Delete(string domain)
        {
            byte[] owner = Storage.Get(Storage.CurrentContext, domain);
            if (owner == null) return false;
            if (!Runtime.CheckWitness(owner)) return false;
            Storage.Delete(Storage.CurrentContext, domain);
            return true;
        }
    }
}


code = 54c56b6c766b00527ac46c766b51527ac4616c766b00c36c766b52527ac46c766b52c305517565727987633a006c766b52c308526567697374657287633d006c766b52c3085472616e73666572876348006c766b52c30644656c657465876355006267006c766b51c300c3616570006c766b53527ac4625d006c766b51c300c36c766b51c351c3617c65a8006c766b53527ac46240006c766b51c300c36c766b51c351c3617c654e016c766b53527ac46223006c766b51c300c361653b026c766b53527ac4620e00006c766b53527ac46203006c766b53c3616c756652c56b6c766b00527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b51527ac46203006c766b51c3616c756655c56b6c766b00527ac46c766b51527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b52527ac46c766b52c300a06c766b53527ac46c766b53c3640e00006c766b54527ac4624c0061681953797374656d2e53746f726167652e476574436f6e746578746c766b00c36c766b51c3615272681253797374656d2e53746f726167652e50757461516c766b54527ac46203006c766b54c3616c756656c56b6c766b00527ac46c766b51527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b52527ac46c766b52c3009c6c766b53527ac46c766b53c3640e00006c766b54527ac4628b006c766b52c361681b53797374656d2e52756e74696d652e436865636b5769746e657373009c6c766b55527ac46c766b55c3640e00006c766b54527ac4624c0061681953797374656d2e53746f726167652e476574436f6e746578746c766b00c36c766b51c3615272681253797374656d2e53746f726167652e50757461516c766b54527ac46203006c766b54c3616c756655c56b6c766b00527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b51527ac46c766b51c3009c6c766b52527ac46c766b52c3640e00006c766b53527ac46288006c766b51c361681b53797374656d2e52756e74696d652e436865636b5769746e657373009c6c766b54527ac46c766b54c3640e00006c766b53527ac462490061681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681553797374656d2e53746f726167652e44656c65746561516c766b53527ac46203006c766b53c3616c7566
*/

func TestDomainSmartContract(ctx *testframework.TestFrameworkContext) bool {
	code := "54c56b6c766b00527ac46c766b51527ac4616c766b00c36c766b52527ac46c766b52c305517565727987633a006c766b52c308526567697374657287633d006c766b52c3085472616e73666572876348006c766b52c30644656c657465876355006267006c766b51c300c3616570006c766b53527ac4625d006c766b51c300c36c766b51c351c3617c65a8006c766b53527ac46240006c766b51c300c36c766b51c351c3617c654e016c766b53527ac46223006c766b51c300c361653b026c766b53527ac4620e00006c766b53527ac46203006c766b53c3616c756652c56b6c766b00527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b51527ac46203006c766b51c3616c756655c56b6c766b00527ac46c766b51527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b52527ac46c766b52c300a06c766b53527ac46c766b53c3640e00006c766b54527ac4624c0061681953797374656d2e53746f726167652e476574436f6e746578746c766b00c36c766b51c3615272681253797374656d2e53746f726167652e50757461516c766b54527ac46203006c766b54c3616c756656c56b6c766b00527ac46c766b51527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b52527ac46c766b52c3009c6c766b53527ac46c766b53c3640e00006c766b54527ac4628b006c766b52c361681b53797374656d2e52756e74696d652e436865636b5769746e657373009c6c766b55527ac46c766b55c3640e00006c766b54527ac4624c0061681953797374656d2e53746f726167652e476574436f6e746578746c766b00c36c766b51c3615272681253797374656d2e53746f726167652e50757461516c766b54527ac46203006c766b54c3616c756655c56b6c766b00527ac46161681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681253797374656d2e53746f726167652e4765746c766b51527ac46c766b51c3009c6c766b52527ac46c766b52c3640e00006c766b53527ac46288006c766b51c361681b53797374656d2e52756e74696d652e436865636b5769746e657373009c6c766b54527ac46c766b54c3640e00006c766b53527ac462490061681953797374656d2e53746f726167652e476574436f6e746578746c766b00c3617c681553797374656d2e53746f726167652e44656c65746561516c766b53527ac46203006c766b53c3616c7566"
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