package native

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ontio/ontology-crypto/keypair"
	sdkcomm "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/auth"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	cstates "github.com/ontio/ontology/smartcontract/states"
	vmtypes "github.com/ontio/ontology/smartcontract/types"
)

var (
	contractAddress = utils.AuthContractAddress
	adminOntID      = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
	newAdminOntID   = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
	p1              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
	p2              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
	p3              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
	p4              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
)

func TestAuthContract(ctx *testframework.TestFrameworkContext) bool {
	user, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError(err)
		return false
	}

	ret, err := TestInitContractAdmin(ctx, user)
	if err != nil {
		fmt.Println(err)
		//log.Error(err)
	}
	ret, err = TestAssignFuncsToRole(ctx, user)
	if err != nil {
		fmt.Println(err)
	}
	ret, err = TestAssignOntIDsToRole(ctx, user)
	if err != nil {
		fmt.Println(err)
	}
	ret, err = TestDelegate(ctx, user)
	if err != nil {
		fmt.Println(err)
	}
	ret, err = TestWithdraw(ctx, user)
	if err != nil {
		fmt.Println(err)
	}
	return ret
}

func sendTestTx(ctx *testframework.TestFrameworkContext, user *account.Account, param []byte, method string) (bool, error) {
	crt := &cstates.Contract{
		Address: contractAddress,
		Method:  method,
		Args:    param,
	}
	buf := bytes.NewBuffer(nil)
	if err := crt.Serialize(buf); err != nil {
		return false, fmt.Errorf("Serialize contract error:%s", err)
	}
	//prepare tx
	invokeTx := sdkcomm.NewInvokeTransaction(common.DefConfig.GasPrice, common.DefConfig.GasLimit, vmtypes.Native, buf.Bytes())
	if err := sdkcomm.SignTransaction(invokeTx, user); err != nil {
		return false, fmt.Errorf("SignTransaction error:%s", err)
	}
	if _, err := ctx.Ont.Rpc.SendRawTransaction(invokeTx); err != nil {
		return false, fmt.Errorf("SendTransaction error:%s", err)
	}
	if _, err := ctx.Ont.Rpc.WaitForGenerateBlock(15*time.Second, 1); err != nil {
		ctx.LogError("WaitForGenerateBlock error:%s", err)
		return false, err
	}
	return true, nil
}

func TestInitContractAdmin(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {
	//prepare invoke param
	param := &auth.InitContractAdminParam{
		AdminOntID: adminOntID,
	}
	buf := bytes.NewBuffer(nil)
	if err := param.Serialize(buf); err != nil {
		ctx.LogError(err)
		return false, err
	}
	//prepare contract
	method := "initContractAdmin"
	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
	if err != nil {
		ctx.LogError(err)
		return false, err
	}
	if !ret {
		return false, nil
	}

	//check result
	key, err := auth.PackKeys([]byte{}, [][]byte{utils.OntContractAddress[:], auth.Admin})
	if err != nil {
		return false, err
	}
	val, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return false, err
	}
	//ctx.LogInfo(hex.EncodeToString(val))
	if bytes.Compare(val, param.AdminOntID) != 0 {
		return false, nil
	}

	return true, nil
}

func TestAssignFuncsToRole(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {
	//prepare invoke param
	param := &auth.FuncsToRoleParam{
		ContractAddr: utils.OntContractAddress[:],
		AdminOntID:   newAdminOntID,
		Role:         []byte("role"),
		FuncNames:    []string{"init", "transfer"},
	}
	buf := bytes.NewBuffer(nil)
	if err := param.Serialize(buf); err != nil {
		//ctx.LogError()
		return false, err
	}
	method := "assignFuncsToRole"
	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
	if err != nil {
		ctx.LogError(err)
		return false, err
	}
	if !ret {
		return false, nil
	}

	roleF, err := auth.PackKeys([]byte{}, [][]byte{utils.OntContractAddress[:], auth.RoleF, []byte("role")})
	//roleF = append(utils.AuthContractAddress[:], roleF...)
	q, err := ctx.Ont.Rpc.GetStorage(contractAddress, roleF)

	if err != nil {
		return false, err
	}
	for {
		if q == nil {
			break
		}
		ctx.LogInfo(string(q))
		next, err := ctx.Ont.Rpc.GetStorage(contractAddress, append(roleF, q...))
		if err != nil {
			return false, err
		}
		fn := new(utils.LinkedlistNode)
		fn.Deserialize(next)
		q = fn.GetNext()
	}
	return true, nil
}

func TestAssignOntIDsToRole(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {

	param := &auth.OntIDsToRoleParam{
		ContractAddr: utils.OntContractAddress[:],
		AdminOntID:   adminOntID,
		Role:         []byte("role"),
		Persons:      [][]byte{p1, p2},
	}
	buf := bytes.NewBuffer(nil)
	if err := param.Serialize(buf); err != nil {
		//ctx.LogError()
		return false, err
	}
	method := "assignOntIDsToRole"
	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
	if err != nil {
		return false, err
	}
	if !ret {
		return false, errors.NewErr("send test transaction failed")
	}
	roleP, err := auth.PackKeys([]byte{}, [][]byte{utils.OntContractAddress[:], auth.RoleP, []byte("role")})
	//roleF = append(utils.AuthContractAddress[:], roleF...)
	q, err := ctx.Ont.Rpc.GetStorage(contractAddress, roleP)

	if err != nil {
		return false, err
	}
	for {
		if q == nil {
			break
		}
		ctx.LogInfo(hex.EncodeToString(q))
		next, err := ctx.Ont.Rpc.GetStorage(contractAddress, append(roleP, q...))
		if err != nil {
			return false, err
		}
		per := new(utils.LinkedlistNode)
		per.Deserialize(next)
		q = per.GetNext()
		//ctx.LogInfo(hex.EncodeToString(per.payload))

	}
	return true, nil
}

func TestDelegate(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {

	param := &auth.DelegateParam{
		ContractAddr: utils.OntContractAddress[:],
		From:         p1,
		To:           p3,
		Role:         []byte("role"),
		Period:       60 * 60 * 24 * 30,
		Level:        3,
	}
	buf := bytes.NewBuffer(nil)
	if err := param.Serialize(buf); err != nil {
		//ctx.LogError()
		return false, err
	}
	method := "delegate"
	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
	if err != nil {
		return false, err
	}
	if !ret {
		return false, errors.NewErr("send test transaction failed")
	}
	param = &auth.DelegateParam{
		ContractAddr: utils.OntContractAddress[:],
		From:         p3,
		To:           p4,
		Role:         []byte("role"),
		Period:       60 * 60 * 24 * 30,
		Level:        2,
	}
	buf = bytes.NewBuffer(nil)
	if err := param.Serialize(buf); err != nil {
		//ctx.LogError()
		return false, err
	}
	method = "delegate"
	ret, err = sendTestTx(ctx, user, buf.Bytes(), method)
	if err != nil {
		return false, err
	}
	if !ret {
		return false, errors.NewErr("send test transaction failed")
	}
	roleP, err := auth.PackKeys([]byte{}, [][]byte{utils.OntContractAddress[:], auth.RoleP, []byte("role")})
	//roleF = append(utils.AuthContractAddress[:], roleF...)
	q, err := ctx.Ont.Rpc.GetStorage(contractAddress, roleP)
	ctx.LogInfo("after delegate")
	if err != nil {
		return false, err
	}
	for {
		if q == nil {
			break
		}
		ctx.LogInfo(hex.EncodeToString(q))
		next, err := ctx.Ont.Rpc.GetStorage(contractAddress, append(roleP, q...))
		if err != nil {
			return false, err
		}
		fn := new(utils.LinkedlistNode)
		fn.Deserialize(next)
		q = fn.GetNext()
	}
	return true, nil
}

func TestWithdraw(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {
	param := &auth.DelegateParam{
		ContractAddr: utils.OntContractAddress[:],
		From:         p1,
		To:           p3,
		Role:         []byte("role"),
		Period:       60 * 60 * 24 * 30,
		Level:        2,
	}
	buf := bytes.NewBuffer(nil)
	if err := param.Serialize(buf); err != nil {
		//ctx.LogError()
		return false, err
	}
	method := "withdraw"
	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
	if err != nil {
		return false, err
	}
	if !ret {
		return false, errors.NewErr("send test transaction failed")
	}
	roleP, err := auth.PackKeys([]byte{}, [][]byte{utils.OntContractAddress[:], auth.RoleP, []byte("role")})
	//roleF = append(utils.AuthContractAddress[:], roleF...)
	q, err := ctx.Ont.Rpc.GetStorage(contractAddress, roleP)
	ctx.LogInfo("after withdraw")
	if err != nil {
		return false, err
	}
	for {
		if q == nil {
			break
		}
		ctx.LogInfo(hex.EncodeToString(q))
		next, err := ctx.Ont.Rpc.GetStorage(contractAddress, append(roleP, q...))
		if err != nil {
			return false, err
		}
		fn := new(utils.LinkedlistNode)
		fn.Deserialize(next)
		q = fn.GetNext()
	}
	return true, nil
}
