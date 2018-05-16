package ontid

import (
	"bytes"
	"crypto/rand"
	"math/big"

	"github.com/ontio/ontology-crypto/keypair"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/core/genesis"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/states"
)

var test_id = "did:ont:abcd1234"

func TestID(ctx *testframework.TestFrameworkContext) bool {
	genID()
	if !registerID(ctx) {
		ctx.LogError("register ID failed")
		return false
	}
	if !addKey(ctx) {
		ctx.LogError("add key failed")
		return false
	}
	if !testRecovery(ctx) {
		ctx.LogError("recovery failed")
		return false
	}

	queryDDO(ctx)
	return true
}

func TestAttr(ctx *testframework.TestFrameworkContext) bool {
	genID()
	if !regIDWithAttr(ctx) {
		ctx.LogError("test register with attributes failed")
		return false
	}
	if !testAttr(ctx) {
		ctx.LogError("test attribute failed")
		return false
	}
	queryDDO(ctx)
	return true
}

func genID() {
	max, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFF", 16)
	r, _ := rand.Int(rand.Reader, max)
	test_id = "did:ont:" + r.String()
}

func registerID(ctx *testframework.TestFrameworkContext) bool {
	user, err := ctx.GetDefaultAccount()
	if err != nil {
		ctx.LogError("Wallet.CreateAccount error:%s", err)
		return false
	}
	pub := keypair.SerializePublicKey(user.PublicKey)

	buf := bytes.NewBuffer(nil)

	serialization.WriteVarBytes(buf, []byte(test_id))
	serialization.WriteVarBytes(buf, pub)
	args := buf.Bytes()

	caddr := genesis.OntIDContractAddress
	inv := &states.Contract{
		Address: caddr,
		Method:  "regIDWithPublicKey",
		Args:    args,
	}

	ok, _ := InvokeContract(ctx, inv, false)
	return ok
}

func addKey(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	_, p, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	key := keypair.SerializePublicKey(p)

	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, []byte(key))
	serialization.WriteVarBytes(&buf, pub)
	args := buf.Bytes()

	c := &states.Contract{
		Address: genesis.OntIDContractAddress,
		Method:  "addKey",
		Args:    args,
	}
	ok, _ := InvokeContract(ctx, c, false)
	if !ok {
		return false
	}

	c.Method = "removeKey"
	ok, _ = InvokeContract(ctx, c, false)

	c.Method = "getKeyState"
	for i := 1; i <= 3; i++ {
		c.Args = keyStateArg(uint32(i))
		InvokeContract(ctx, c, false)
	}

	return ok
}

func keyStateArg(i uint32) []byte {
	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteUint32(&buf, i)
	return buf.Bytes()
}

func queryDDO(ctx *testframework.TestFrameworkContext) bool {
	c := &states.Contract{
		Address: genesis.OntIDContractAddress,
		Method:  "getDDO",
		Args:    []byte(test_id),
	}
	ok, _ := InvokeContract(ctx, c, false)
	if !ok {
		return false
	}
	//ctx.LogInfo(hex.EncodeToString(res))

	return true
}

func testRecovery(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	pubs := []keypair.PublicKey{user.PublicKey}
	for i := 0; i < 3; i++ {
		_, p, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
		pubs = append(pubs, p)
	}
	ctx.LogError(pubs)
	addr0, _ := types.AddressFromMultiPubKeys(pubs, 1)
	addr1, _ := types.AddressFromMultiPubKeys(pubs[:3], 1)

	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, addr0[:])
	serialization.WriteVarBytes(&buf, pub)
	args := buf.Bytes()

	c := &states.Contract{
		Address: genesis.OntIDContractAddress,
		Method:  "addRecovery",
		Args:    args,
	}
	ok, _ := InvokeContract(ctx, c, false)
	if !ok {
		ctx.LogError("add recovery failed")
		return false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, addr1[:])
	serialization.WriteVarBytes(&buf, addr0[:])
	args = buf.Bytes()

	c.Method = "changeRecovery"
	c.Args = args

	tx, err := makeTx(c)
	if err != nil {
		ctx.LogError(err)
		return false
	}
	var ac []*account.Account
	for _, v := range pubs {
		ac = append(ac, &account.Account{PublicKey: v})
	}
	ac[0].PrivateKey = user.PrivateKey

	err = sdkcom.MultiSignTransaction("SHA256withECDSA", tx, 1, ac)
	if err != nil {
		ctx.LogError("SignTransaction error: %s", err)
		return false
	}

	ok = sendTx(ctx, tx)
	return ok
}

func regIDWithAttr(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte("attr0"))
	serialization.WriteVarBytes(&buf, []byte{0})
	serialization.WriteVarBytes(&buf, []byte{1})
	serialization.WriteVarBytes(&buf, []byte("attr1"))
	serialization.WriteVarBytes(&buf, []byte{1})
	serialization.WriteVarBytes(&buf, []byte{0, 1})
	attrs := buf.Bytes()

	var buf1 bytes.Buffer
	serialization.WriteVarBytes(&buf1, []byte(test_id))
	serialization.WriteVarBytes(&buf1, pub)
	serialization.WriteVarBytes(&buf1, attrs)

	c := &states.Contract{
		Address: genesis.OntIDContractAddress,
		Method:  "regIDWithAttributes",
		Args:    buf1.Bytes(),
	}

	ok, _ := InvokeContract(ctx, c, false)

	return ok
}

func testAttr(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, []byte("attr1"))
	serialization.WriteVarBytes(&buf, []byte{1})
	serialization.WriteVarBytes(&buf, []byte{0x01, 0x02})
	serialization.WriteVarBytes(&buf, pub)
	args := buf.Bytes()

	c := &states.Contract{
		Address: genesis.OntIDContractAddress,
		Method:  "addAttribute",
		Args:    args,
	}

	ok, _ := InvokeContract(ctx, c, false)
	if !ok {
		return false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, []byte("attr2"))
	serialization.WriteVarBytes(&buf, []byte{2})
	serialization.WriteVarBytes(&buf, []byte("abcd"))
	serialization.WriteVarBytes(&buf, pub)
	args = buf.Bytes()
	c.Args = args
	ok, _ = InvokeContract(ctx, c, false)
	if !ok {
		return false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, []byte("attr2"))
	serialization.WriteVarBytes(&buf, pub)
	c.Method = "removeAttribute"
	ok, _ = InvokeContract(ctx, c, false)

	return ok
}
