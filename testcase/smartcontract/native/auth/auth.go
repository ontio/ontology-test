package auth

import (
	"fmt"
	"io"
	"math/big"
	"time"

	//"time"
	"crypto/rand"
	"crypto/sha256"

	base58 "github.com/itchyny/base58-go"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	//"github.com/ontio/ontology/common"
	sdkutils "github.com/ontio/ontology-go-sdk/utils"
	//"github.com/ontio/ontology/core/genesis"
	Auth "github.com/ontio/ontology/smartcontract/service/native/auth"

	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"golang.org/x/crypto/ripemd160"
)

var (
	adminOntID = []byte("did:ont:AeoWeBLt6wnRhow2qyrCuZKcGFU3qhwHnm")
	p1         = genOntID(rand.Reader)
	p2         = genOntID(rand.Reader)
	p3         = genOntID(rand.Reader)
	p4         = genOntID(rand.Reader)
	admin      = account.NewAccount("SHA256withECDSA")
	a1         = account.NewAccount("SHA256withECDSA")
	a2         = account.NewAccount("SHA256withECDSA")
	a3         = account.NewAccount("SHA256withECDSA")
	a4         = account.NewAccount("SHA256withECDSA")
)

var contractAddr common.Address

func TestAuthContract(ctx *testframework.TestFrameworkContext) bool {
	ret := setup(ctx)
	if !ret {
		ctx.LogError("failed to set up test prerequisites")
	}

	/*
		ret = testBasic(ctx)
		if !ret {
			ctx.LogError("basic test failed")
		}

		ret = testDelegate(ctx)
		if !ret {
			ctx.LogError("delegate test failed")
		}
	*/
	ret = testWithdraw(ctx)

	return ret
}

func testWithdraw(ctx *testframework.TestFrameworkContext) bool {
	ret := initContractAdmin(ctx, admin)
	if !ret {
		ctx.LogError("failed to init app contract's admin")
	}

	ret = assignFuncsToRole(ctx, admin, []string{"B", "C"}, "roleB")
	if !ret {
		ctx.LogError("failed to assign funcs to role")
	}

	ret = assignOntIDsToRole(ctx, admin, [][]byte{p2}, "roleB")
	if !ret {
		ctx.LogError("failed to assign OntID to role")
	}

	ret = delegate(ctx, a2, p2, p1, "roleB", 10000)
	if !ret {
		ctx.LogError("failed to delegate ")
	}

	callFunc(ctx, a1, p1, "B")

	ret = withdraw(ctx, a2, p2, p1, "roleB")
	if !ret {
		ctx.LogError("failed to withdraw")
	}

	callFunc(ctx, a1, p1, "B")
	return true
}
func testDelegate(ctx *testframework.TestFrameworkContext) bool {
	ret := initContractAdmin(ctx, admin)
	if !ret {
		ctx.LogError("failed to init app contract's admin")
	}

	ret = assignFuncsToRole(ctx, admin, []string{"B", "C"}, "roleB")
	if !ret {
		ctx.LogError("failed to assign funcs to role")
	}

	ret = assignOntIDsToRole(ctx, admin, [][]byte{p2}, "roleB")
	if !ret {
		ctx.LogError("failed to assign OntID to role")
	}

	ret = delegate(ctx, a2, p2, p1, "roleB", 10000)
	if !ret {
		ctx.LogError("failed to delegate ")
	}

	ret = callFunc(ctx, a1, p1, "B") //should push true event

	return ret
}
func testBasic(ctx *testframework.TestFrameworkContext) bool {
	ret := initContractAdmin(ctx, admin)
	if !ret {
		ctx.LogError("failed to init app contract's admin")
	}

	ret = assignFuncsToRole(ctx, admin, []string{"A", "C"}, "roleA")
	if !ret {
		ctx.LogError("failed to assign funcs to role")
	}

	ret = assignOntIDsToRole(ctx, admin, [][]byte{p1}, "roleA")
	if !ret {
		ctx.LogError("failed to assign OntID to role")
	}

	ret = callFunc(ctx, a1, p1, "B") //should push false event

	ret = callFunc(ctx, a1, p1, "A") //should push true event

	ret = callFunc(ctx, a1, p1, "C") //should push true event
	return true
}
func setup(ctx *testframework.TestFrameworkContext) bool {
	user, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError(err)
		return false
	}

	//get default user's balance
	userBalance, err := ctx.Ont.Rpc.GetBalance(user.Address)
	if err != nil {
		ctx.LogError("cannot get user's balance: ", err)
		return false
	}
	ctx.LogInfo("user's ong balance is %d", userBalance.Ong)
	ctx.LogInfo("user's ont balance is %d", userBalance.Ont)

	ret := deployAppContract(ctx, user)
	if !ret {
		ctx.LogError("deploy AppContract failed")
		return false
	}
	ret = regOntID(ctx, adminOntID, admin)
	if !ret {
		ctx.LogError("reg admin's ONT ID failed: ", string(adminOntID))
	}
	ret = regOntID(ctx, p1, a1)
	if !ret {
		ctx.LogError("reg p1's ONT ID failed: ", string(p1))
	}
	ret = regOntID(ctx, p2, a2)
	if !ret {
		ctx.LogError("reg p2's ONT ID failed: ", string(p2))
	}

	return ret
}

type Param struct {
	ID  []byte
	Key []byte
}

func toBase58(f []byte) string {
	data := append([]byte{23}, f[:]...)
	temp := sha256.Sum256(data)
	temps := sha256.Sum256(temp[:])
	data = append(data, temps[0:4]...)

	bi := new(big.Int).SetBytes(data).String()
	encoded, _ := base58.BitcoinEncoding.Encode([]byte(bi))
	return string(encoded)
}

func genOntID(rand io.Reader) []byte {
	id := make([]byte, 20)

	nonce := make([]byte, 32)
	n, err := rand.Read(nonce)
	if err != nil || n != 32 {
		fmt.Println(fmt.Errorf("genOntID failed"))
	}

	temp := sha256.Sum256(nonce)
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(id[:0])

	idString := toBase58(id)
	idString = "did:ont:" + idString
	return []byte(idString)
}

func deployAppContract(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	/*
		code, err := ioutil.ReadFile("./contract.nvm")
		if err != nil {
			ctx.LogError("read contract code failed: %v", err)
			return false
		}

		contractCode := hex.EncodeToString(code)
	*/
	contractCode := "5ac56b6c766b00527ac46c766b51527ac4616c766b00c304696e6974876c766b52527ac46c766b52c364120061616590016c766b53527ac46205016c766b00c30141876c766b54527ac46c766b54c3643d00616c766b00c36c766b51c3617c653902009c6c766b55527ac46c766b55c3640e00006c766b53527ac462c2006165c7006c766b53527ac462b4006c766b00c30142876c766b56527ac46c766b56c3643d00616c766b00c36c766b51c3617c65e801009c6c766b57527ac46c766b57c3640e00006c766b53527ac462710061659e006c766b53527ac46263006c766b00c30143876c766b58527ac46c766b58c3643d00616c766b00c36c766b51c3617c659701009c6c766b59527ac46c766b59c3640e00006c766b53527ac4622000616575006c766b53527ac46212006c766b00c36c766b53527ac46203006c766b53c3616c756651c56b6110494e564f4b45204120535543434553536c766b00527ac46203006c766b00c3616c756651c56b6110494e564f4b45204220535543434553536c766b00527ac46203006c766b00c3616c756651c56b6110494e564f4b45204320535543434553536c766b00527ac46203006c766b00c3616c756655c56b611400000000000000000000000000000000000000066c766b00527ac4536151c66c766b527a527ac46c766b53c3612a6469643a6f6e743a41656f5765424c7436776e52686f773271797243755a4b6347465533716877486e6d007cc46c766b53c36c766b51527ac4006c766b00c311696e6974436f6e747261637441646d696e6c766b51c361537951795572755172755279527954727552727568164f6e746f6c6f67792e4e61746976652e496e766f6b656c766b52527ac46c766b52c300517f519c6c766b54527ac46203006c766b54c3616c756657c56b6c766b00527ac46c766b51527ac4611400000000000000000000000000000000000000066c766b52527ac4556154c66c766b527a527ac46c766b55c36c766b53527ac46c766b53c361682d53797374656d2e457865637574696f6e456e67696e652e476574457865637574696e6753637269707448617368007cc46c766b53c36c766b00c3527cc46c766b53c36c766b51c300c3517cc46c766b53c36c766b51c351c3537cc4006c766b52c30b766572696679546f6b656e6c766b53c361537951795572755172755279527954727552727568164f6e746f6c6f67792e4e61746976652e496e766f6b656c766b54527ac46c766b54c300517f519c6c766b56527ac46203006c766b56c3616c7566"
	contractAddr, _ = sdkutils.GetContractAddress(contractCode)

	signer, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("depolyAppContract GetDefaultAccount error:%s", err)
		return false
	}
	txHash, err := ctx.Ont.Rpc.DeploySmartContract(ctx.GetGasPrice(), ctx.GetGasLimit(),
		signer,
		true,
		contractCode,
		"AppContract",
		"1.0",
		"",
		"",
		"",
	)
	if err != nil {
		ctx.LogError("depolyAppContract DeploySmartContract error:%s", err)
		return false
	}
	ctx.LogInfo("CodeAddress:%s", contractAddr.ToHexString())
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

	ctx.LogInfo("txHash: %s", txHash.ToHexString())
	return true
}

type Contract struct {
	Address common.Address
	Method  string
	Args    []interface{}
}

func InvokeContract(ctx *testframework.TestFrameworkContext, contract *Contract, user *account.Account, native bool) bool {
	var txHash common.Uint256
	var err error
	if native {
		txHash, err = ctx.Ont.Rpc.InvokeNativeContract(
			ctx.GetGasPrice(),
			ctx.GetGasLimit(),
			user,
			0,
			contract.Address,
			contract.Method,
			contract.Args,
		)
	} else {
		txHash, err = ctx.Ont.Rpc.InvokeNeoVMContract(
			ctx.GetGasPrice(),
			ctx.GetGasLimit(),
			user,
			contract.Address,
			contract.Args,
		)
	}
	ctx.LogInfo("txHash: ", txHash.ToHexString())
	if err != nil {
		ctx.LogError(err)
		return false
	}
	_, err = ctx.Ont.Rpc.WaitForGenerateBlock(30*time.Second, 2)
	if err != nil {
		ctx.LogError("WaitForGenerateBlock error: %s", err)
		return false
	}

	events, err := ctx.Ont.Rpc.GetSmartContractEvent(txHash)
	if err != nil {
		ctx.LogError("GetSmartContractEvent error: %s", err)
		return false
	}

	if events.State == 0 {
		ctx.LogError("ontio contract invoke failed, state:0")
		return false
	}
	if len(events.Notify) > 0 {
		states := events.Notify[0].States
		ctx.LogInfo("result is : %+v", states)
		return true
	} else {
		return false
	}
}

func regOntID(ctx *testframework.TestFrameworkContext, ontID []byte, user *account.Account) bool {
	pubKey := keypair.SerializePublicKey(user.PublicKey)
	param := &Param{
		ontID, pubKey,
	}
	contract := &Contract{
		Address: utils.OntIDContractAddress,
		Method:  "regIDWithPublicKey",
		Args:    []interface{}{param},
	}
	return InvokeContract(ctx, contract, user, true)
}

func initContractAdmin(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	//prepare invoke param
	contract := &Contract{
		Address: (contractAddr),
		Method:  "init",
		Args:    []interface{}{"init", []interface{}{}},
	}
	return InvokeContract(ctx, contract, user, false)
}

func callFunc(ctx *testframework.TestFrameworkContext, user *account.Account, caller []byte, fn string) bool {
	contract := &Contract{
		Address: contractAddr,
		Method:  "",
		Args:    []interface{}{fn, []interface{}{caller, 1}},
	}
	return InvokeContract(ctx, contract, user, false)
}
func assignFuncsToRole(ctx *testframework.TestFrameworkContext, user *account.Account, funcs []string, role string) bool {
	//prepare invoke param
	param := &Auth.FuncsToRoleParam{
		ContractAddr: contractAddr,
		AdminOntID:   adminOntID,
		FuncNames:    funcs,
		Role:         []byte(role),
		KeyNo:        1,
	}
	contract := &Contract{
		Address: utils.AuthContractAddress,
		Method:  "assignFuncsToRole",
		Args:    []interface{}{param},
	}
	return InvokeContract(ctx, contract, user, true)
}

func assignOntIDsToRole(ctx *testframework.TestFrameworkContext, user *account.Account, persons [][]byte, role string) bool {
	param := &Auth.OntIDsToRoleParam{
		ContractAddr: contractAddr,
		AdminOntID:   adminOntID,
		Persons:      persons,
		Role:         []byte(role),
		KeyNo:        1,
	}
	contract := &Contract{
		Address: utils.AuthContractAddress,
		Method:  "assignOntIDsToRole",
		Args:    []interface{}{param},
	}
	return InvokeContract(ctx, contract, user, true)
}

func verifyToken(ctx *testframework.TestFrameworkContext, user *account.Account, caller []byte, fn string) bool {
	param := &Auth.VerifyTokenParam{
		ContractAddr: contractAddr,
		Caller:       caller,
		Fn:           fn,
		KeyNo:        1,
	}
	contract := &Contract{
		Address: utils.AuthContractAddress,
		Method:  "verifyToken",
		Args:    []interface{}{param},
	}
	return InvokeContract(ctx, contract, user, true)
}

func delegate(ctx *testframework.TestFrameworkContext, user *account.Account, from, to []byte, role string, period uint64) bool {
	param := &Auth.DelegateParam{
		ContractAddr: contractAddr,
		From:         from,
		To:           to,
		Role:         []byte(role),
		Period:       period,
		Level:        1,
		KeyNo:        1,
	}
	contract := &Contract{
		Address: utils.AuthContractAddress,
		Method:  "delegate",
		Args:    []interface{}{param},
	}
	return InvokeContract(ctx, contract, user, true)
}

func withdraw(ctx *testframework.TestFrameworkContext, user *account.Account, from, to []byte,
	role string) bool {
	param := &Auth.WithdrawParam{
		ContractAddr: contractAddr,
		Initiator:    from,
		Delegate:     to,
		Role:         []byte(role),
		KeyNo:        1,
	}
	contract := &Contract{
		Address: utils.AuthContractAddress,
		Method:  "withdraw",
		Args:    []interface{}{param},
	}
	return InvokeContract(ctx, contract, user, true)
}
