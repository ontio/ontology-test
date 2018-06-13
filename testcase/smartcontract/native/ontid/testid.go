package ontid

import (
	"bytes"
	"crypto/rand"
	"github.com/ontio/ontology-crypto/keypair"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"github.com/ontio/ontology/smartcontract/states"
	"math/big"
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

	res := true

	caddr := utils.OntIDContractAddress
	inv := &states.Contract{
		Address: caddr,
		Method:  "regIDWithPublicKey",
		Args:    nil,
	}
	ok, _ := InvokeContract(ctx, inv, false)
	if ok {
		ctx.LogError("register should failed without arguments")
		res = false
	}

	buf := bytes.NewBuffer(nil)
	serialization.WriteVarBytes(buf, nil)
	serialization.WriteVarBytes(buf, pub)
	inv.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, inv, false)
	if ok {
		ctx.LogError("register should failed with invalid id")
		res = false
	}
	buf.Reset()
	serialization.WriteVarBytes(buf, []byte(test_id))
	serialization.WriteVarBytes(buf, nil)
	inv.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, inv, false)
	if ok {
		ctx.LogError("register should failed with invalid key")
		res = false
	}

	buf.Reset()
	serialization.WriteVarBytes(buf, []byte(test_id))
	serialization.WriteVarBytes(buf, pub)
	inv.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, inv, false)

	return ok && res
}

func addKey(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	_, p, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	key := keypair.SerializePublicKey(p)
	_, p1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	key1 := keypair.SerializePublicKey(p1)

	res := true

	c := &states.Contract{
		Address: utils.OntIDContractAddress,
		Method:  "addKey",
		Args:    nil,
	}
	ok, _ := InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed without arguments")
		res = false
	}

	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, nil)
	serialization.WriteVarBytes(&buf, []byte(key))
	serialization.WriteVarBytes(&buf, pub)
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with invalid id")
		res = false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte("123"))
	serialization.WriteVarBytes(&buf, key)
	serialization.WriteVarBytes(&buf, pub)
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with unregistered id")
		res = false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, key)
	serialization.WriteVarBytes(&buf, nil)
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with empty key")
		res = false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, key)
	serialization.WriteVarBytes(&buf, key1)
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with wrong key")
		res = false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, key)
	serialization.WriteVarBytes(&buf, pub)
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if !ok {
		return false
	}

	c.Method = "removeKey"
	ok, _ = InvokeContract(ctx, c, false)
	res = ok && res
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("removeKey should failed while removing the same key again")
		res = false
	}

	c.Method = "getKeyState"
	for i := 1; i <= 3; i++ {
		c.Args = keyStateArg(uint64(i))
		InvokeContract(ctx, c, false)
	}

	return res
}

func keyStateArg(i uint64) []byte {
	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarUint(&buf, i)
	return buf.Bytes()
}

func queryDDO(ctx *testframework.TestFrameworkContext) bool {
	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte(test_id))
	c := &states.Contract{
		Address: utils.OntIDContractAddress,
		Method:  "getDDO",
		Args:    buf.Bytes(),
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
	addr0.Serialize(&buf)
	//serialization.WriteVarBytes(&buf, addr0[:])
	serialization.WriteVarBytes(&buf, pub)
	args := buf.Bytes()

	c := &states.Contract{
		Address: utils.OntIDContractAddress,
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
	addr1.Serialize(&buf)
	addr0.Serialize(&buf)
	//serialization.WriteVarBytes(&buf, addr1[:])
	//serialization.WriteVarBytes(&buf, addr0[:])
	args = buf.Bytes()

	c.Method = "changeRecovery"
	c.Args = args

	tx, err := makeTx(c)
	if err != nil {
		ctx.LogError(err)
		return false
	}

	err = sdkcom.MultiSignToTransaction(tx, 1, pubs, user)
	if err != nil {
		ctx.LogError("MultiSignToTransaction error: %s", err)
		return false
	}
	ok = sendTx(ctx, tx)
	return ok
}

func regIDWithAttr(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	res := true

	c := &states.Contract{
		Address: utils.OntIDContractAddress,
		Method:  "regIDWithAttributes",
		Args:    nil,
	}

	ok, _ := InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed without arguments")
		res = false
	}

	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, nil)
	serialization.WriteVarBytes(&buf, pub)
	serialization.WriteVarUint(&buf, 1)
	serialization.WriteVarBytes(&buf, []byte("attr0"))
	serialization.WriteVarBytes(&buf, []byte{0})
	serialization.WriteVarBytes(&buf, []byte{1})
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed with invalid id")
		res = false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, nil)
	serialization.WriteVarUint(&buf, 1)
	serialization.WriteVarBytes(&buf, []byte("attr0"))
	serialization.WriteVarBytes(&buf, []byte{0})
	serialization.WriteVarBytes(&buf, []byte{1})
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed with invalid key")
		res = false
	}

	buf.Reset()
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, pub)
	serialization.WriteVarUint(&buf, 2)
	serialization.WriteVarBytes(&buf, []byte("attr0"))
	serialization.WriteVarBytes(&buf, []byte{0})
	serialization.WriteVarBytes(&buf, []byte{1})
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed with invalid attributes")
		res = false
	}

	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarBytes(&buf, pub)
	serialization.WriteVarUint(&buf, 2)
	serialization.WriteVarBytes(&buf, []byte("attr0"))
	serialization.WriteVarBytes(&buf, []byte{0})
	serialization.WriteVarBytes(&buf, []byte{1})
	serialization.WriteVarBytes(&buf, []byte("attr1"))
	serialization.WriteVarBytes(&buf, []byte{1})
	serialization.WriteVarBytes(&buf, []byte{0, 1})
	c.Args = buf.Bytes()
	ok, _ = InvokeContract(ctx, c, false)
	return ok && res
}

func testAttr(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	var buf bytes.Buffer
	serialization.WriteVarBytes(&buf, []byte(test_id))
	serialization.WriteVarUint(&buf, 2)
	serialization.WriteVarBytes(&buf, []byte("attr1"))
	serialization.WriteVarBytes(&buf, []byte{1})
	serialization.WriteVarBytes(&buf, []byte{0x01, 0x02})
	serialization.WriteVarBytes(&buf, []byte("attr2"))
	serialization.WriteVarBytes(&buf, []byte{2})
	serialization.WriteVarBytes(&buf, []byte("abcd"))
	serialization.WriteVarBytes(&buf, pub)
	args := buf.Bytes()

	c := &states.Contract{
		Address: utils.OntIDContractAddress,
		Method:  "addAttributes",
		Args:    args,
	}

	ok, _ := InvokeContract(ctx, c, false)
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
