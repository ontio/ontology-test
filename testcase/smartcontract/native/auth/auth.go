package auth
//
//import (
//	"bytes"
//	"fmt"
//	"time"
//
//	"github.com/Ontology/core/genesis"
//	"github.com/ontio/ontology-crypto/keypair"
//	sdkcomm "github.com/ontio/ontology-go-sdk/common"
//	"github.com/ontio/ontology-test/testframework"
//	"github.com/ontio/ontology/account"
//	"github.com/ontio/ontology/common"
//	"github.com/ontio/ontology/smartcontract/service/native/auth"
//	"github.com/ontio/ontology/smartcontract/service/native/utils"
//	cstates "github.com/ontio/ontology/smartcontract/states"
//	vmtypes "github.com/ontio/ontology/smartcontract/types"
//)
//
//var (
//	contractAddress = utils.AuthContractAddress
//	adminOntID      = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
//	newAdminOntID   = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
//	p1              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
//	p2              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
//	p3              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
//	p4              = keypair.SerializePublicKey(account.NewAccount("SHA256withECDSA").PublicKey)
//)
//
//func TestAuthContract(ctx *testframework.TestFrameworkContext) bool {
//	user, err := ctx.GetDefaultAccount()
//	if err != nil {
//		ctx.LogError(err)
//		return false
//	}
//
//	_, err = InitContractAdmin(ctx, user)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = AssignFuncsToRole(ctx, user)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = AssignOntIDsToRole(ctx, user)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = verifyToken(ctx, user, p1, "init")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = verifyToken(ctx, user, p1, "transfer")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = Delegate(ctx, user, p1, p3, "role")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = verifyToken(ctx, user, p3, "init")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = Withdraw(ctx, user, p1, p3, "role")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = verifyToken(ctx, user, p3, "init")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = Delegate(ctx, user, p1, p3, "role")
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	_, err = verifyToken(ctx, user, p3, "init")
//	if err != nil {
//		fmt.Println(err)
//	}
//	return true
//}
//
//func InitContractAdmin(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {
//	//prepare invoke param
//	param := &auth.InitContractAdminParam{
//		AdminOntID: adminOntID,
//	}
//	buf := bytes.NewBuffer(nil)
//	if err := param.Serialize(buf); err != nil {
//		return false, err
//	}
//
//	//prepare contract
//	method := "initContractAdmin"
//	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
//	if err != nil {
//		return false, err
//	}
//	return ret, nil
//}
//
//func AssignFuncsToRole(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {
//	//prepare invoke param
//	param := &auth.FuncsToRoleParam{
//		ContractAddr: utils.OntContractAddress[:],
//		AdminOntID:   adminOntID,
//		Role:         []byte("role"),
//		FuncNames:    []string{"init", "transfer"},
//		KeyNo:        1,
//	}
//	buf := bytes.NewBuffer(nil)
//	if err := param.Serialize(buf); err != nil {
//		return false, err
//	}
//
//	method := "assignFuncsToRole"
//	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
//	if err != nil {
//		ctx.LogError(err)
//		return false, err
//	}
//	return ret, nil
//}
//
//func AssignOntIDsToRole(ctx *testframework.TestFrameworkContext, user *account.Account) (bool, error) {
//	param := &auth.OntIDsToRoleParam{
//		ContractAddr: utils.OntContractAddress[:],
//		AdminOntID:   adminOntID,
//		Role:         []byte("role"),
//		Persons:      [][]byte{p1, p2},
//		KeyNo:        1,
//	}
//	buf := bytes.NewBuffer(nil)
//	if err := param.Serialize(buf); err != nil {
//		//ctx.LogError()
//		return false, err
//	}
//
//	method := "assignOntIDsToRole"
//	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
//	if err != nil {
//		return false, err
//	}
//	return ret, nil
//}
//
//func Delegate(ctx *testframework.TestFrameworkContext, user *account.Account, from, to []byte, role string) (bool, error) {
//	param := &auth.DelegateParam{
//		ContractAddr: utils.OntContractAddress[:],
//		From:         from,
//		To:           to,
//		Role:         []byte(role),
//		Period:       2,
//		Level:        1,
//		KeyNo:        1,
//	}
//	buf := bytes.NewBuffer(nil)
//	if err := param.Serialize(buf); err != nil {
//		//ctx.LogError()
//		return false, err
//	}
//	method := "delegate"
//	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
//	if err != nil {
//		return false, err
//	}
//	return ret, nil
//
//}
//
//func Withdraw(ctx *testframework.TestFrameworkContext, user *account.Account, from, to []byte,
//	role string) (bool, error) {
//	param := &auth.WithdrawParam{
//		ContractAddr: genesis.OntContractAddress[:],
//		Initiator:    from,
//		Delegate:     to,
//		Role:         []byte(role),
//		KeyNo:        1,
//	}
//	buf := bytes.NewBuffer(nil)
//	if err := param.Serialize(buf); err != nil {
//		//ctx.LogError()
//		return false, err
//	}
//	method := "withdraw"
//	_, err := sendTestTx(ctx, user, buf.Bytes(), method)
//	if err != nil {
//		return false, err
//	}
//	return true, nil
//}
//
//func verifyToken(ctx *testframework.TestFrameworkContext, user *account.Account, caller []byte, fn string) (bool, error) {
//	param := &auth.VerifyTokenParam{
//		ContractAddr: utils.OntContractAddress[:],
//		Caller:       caller,
//		Fn:           fn,
//		KeyNo:        1,
//	}
//	buf := bytes.NewBuffer(nil)
//	if err := param.Serialize(buf); err != nil {
//		//ctx.LogError()
//		return false, err
//	}
//	method := "verifyToken"
//	ret, err := sendTestTx(ctx, user, buf.Bytes(), method)
//	if err != nil {
//		return false, err
//	}
//	return ret, nil
//}
//
//func sendTestTx(ctx *testframework.TestFrameworkContext, user *account.Account, param []byte, method string) (bool, error) {
//	var txHash common.Uint256
//
//	crt := &cstates.Contract{
//		Address: contractAddress,
//		Method:  method,
//		Args:    param,
//	}
//	buf := bytes.NewBuffer(nil)
//	if err := crt.Serialize(buf); err != nil {
//		return false, fmt.Errorf("Serialize contract error:%s", err)
//	}
//
//	//prepare tx
//	invokeTx := sdkcomm.NewInvokeTransaction(ctx.GetGasPrice(), ctx.GetGasLimit(), vmtypes.Native, buf.Bytes())
//	if err := sdkcomm.SignTransaction(invokeTx, user); err != nil {
//		return false, fmt.Errorf("SignTransaction error:%s", err)
//	}
//	txHash, err := ctx.Ont.Rpc.SendRawTransaction(invokeTx)
//	if err != nil {
//		return false, fmt.Errorf("SendTransaction error:%s", err)
//	}
//
//	//wait
//	if _, err := ctx.Ont.Rpc.WaitForGenerateBlock(15*time.Second, 1); err != nil {
//		ctx.LogError("WaitForGenerateBlock error:%s", err)
//		return false, err
//	}
//
//	//print event
//	events, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
//	if err != nil {
//		ctx.LogError("GetSmartContractEvent error: %s", err)
//		return false, err
//	}
//
//	if events.State == 0 {
//		ctx.LogError("contract invoke failed, state:0")
//		return false, nil
//	}
//	if len(events.Notify) > 0 {
//		states := events.Notify[0].States
//		ctx.LogInfo("result is : %+v", states)
//		return false, nil
//	} else {
//		return true, nil
//	}
//}
