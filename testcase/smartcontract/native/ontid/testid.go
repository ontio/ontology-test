package ontid

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	base58 "github.com/itchyny/base58-go"
	"github.com/ontio/ontology-crypto/keypair"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

var test_id = "did:ont:abcd1234"

func TestID(ctx *testframework.TestFrameworkContext) bool {
	genID()
	ctx.LogInfo("ID: %s", test_id)
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
	buf := make([]byte, 32)
	rand.Read(buf)
	h := sha256.Sum256(buf)
	data := append([]byte{0x41}, h[:20]...)
	checksum := sha256.Sum256(data)
	checksum = sha256.Sum256(checksum[:])
	data = append(data, checksum[0:4]...)
	b, err := base58.BitcoinEncoding.Encode([]byte(new(big.Int).SetBytes(data).String()))
	if err != nil {
		fmt.Println(err)
	}
	test_id = "did:ont:" + string(b)
}

type Param struct {
	ID  string
	Key []byte
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

	inv := &Contract{
		Address: caddr,
		Method:  "regIDWithPublicKey",
		Args:    nil,
	}
	ok, _ := InvokeContract(ctx, inv, false)
	if ok {
		ctx.LogError("register should failed without arguments")
		res = false
	}

	inv.Args = []interface{}{Param{"", pub}}
	ok, _ = InvokeContract(ctx, inv, false)
	if ok {
		ctx.LogError("register should failed with invalid id")
		res = false
	}
	inv.Args = []interface{}{Param{test_id, nil}}
	ok, _ = InvokeContract(ctx, inv, false)
	if ok {
		ctx.LogError("register should failed with invalid key")
		res = false
	}
	inv.Args = []interface{}{Param{test_id, pub}}
	ok, _ = InvokeContract(ctx, inv, false)

	return ok && res
}

type AddKeyParam struct {
	ID     string
	NewKey []byte
	Key    []byte
}

func addKey(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	_, p, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	key := keypair.SerializePublicKey(p)
	_, p1, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	key1 := keypair.SerializePublicKey(p1)

	res := true

	c := &Contract{
		Address: utils.OntIDContractAddress,
		Method:  "addKey",
		Args:    nil,
	}
	ok, _ := InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed without arguments")
		res = false
	}

	c.Args = []interface{}{AddKeyParam{"", key, pub}}
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with invalid id")
		res = false
	}

	c.Args = []interface{}{AddKeyParam{"123", key, pub}}
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with unregistered id")
		res = false
	}

	c.Args = []interface{}{AddKeyParam{test_id, key, nil}}
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with empty key")
		res = false
	}

	c.Args = []interface{}{AddKeyParam{test_id, key, key1}}
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("addKey should failed with wrong key")
		res = false
	}

	c.Args = []interface{}{AddKeyParam{test_id, key, pub}}
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
		c.Args = []interface{}{StateParam{test_id, i}}
		InvokeContract(ctx, c, false)
	}

	return res
}

type StateParam struct {
	ID    string
	Index int
}

func queryDDO(ctx *testframework.TestFrameworkContext) bool {
	c := &Contract{
		Address: utils.OntIDContractAddress,
		Method:  "getDDO",
		Args:    []interface{}{[]byte(test_id)},
	}
	ok, buf := InvokeContract(ctx, c, true)
	if !ok {
		return false
	}
	ctx.LogInfo(hex.EncodeToString(buf))

	return true
}

type AddRecParam struct {
	ID   string
	Addr common.Address
	Key  []byte
}

type ChangeRecParam struct {
	ID      string
	NewAddr common.Address
	Addr    common.Address
}

func testRecovery(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	pubs := []keypair.PublicKey{user.PublicKey}
	for i := 0; i < 3; i++ {
		_, p, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
		pubs = append(pubs, p)
	}

	addr0, _ := types.AddressFromMultiPubKeys(pubs, 1)
	addr1, _ := types.AddressFromMultiPubKeys(pubs[:3], 1)

	args := []interface{}{AddRecParam{test_id, addr0, pub}}
	c := &Contract{
		Address: utils.OntIDContractAddress,
		Method:  "addRecovery",
		Args:    args,
	}
	ok, _ := InvokeContract(ctx, c, false)
	if !ok {
		ctx.LogError("add recovery failed")
		return false
	}

	c.Method = "changeRecovery"
	c.Args = []interface{}{ChangeRecParam{test_id, addr1, addr0}}

	tx, err := ctx.Ont.Rpc.NewNativeInvokeTransaction(
		ctx.GetGasPrice(),
		ctx.GetGasLimit(),
		0,
		c.Address,
		c.Method,
		c.Args,
	)
	if err != nil {
		ctx.LogError(err)
		return false
	}

	err = sdkcom.MultiSignToTransaction(tx, 1, pubs, user)
	if err != nil {
		ctx.LogError("MultiSignToTransaction error: %s", err)
		return false
	}
	txHash, err := ctx.Ont.Rpc.SendRawTransaction(tx)
	if err != nil {
		ctx.LogError("SendRawTransaction error: %s", err)
	}
	ok = getEvent(ctx, txHash)

	return ok
}

type Attribute struct {
	Key     []byte
	Value   []byte
	Valtype []byte
}

type RegIDAttrParam struct {
	ID   string
	Key  []byte
	Attr []Attribute
}

func regIDWithAttr(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	res := true

	c := &Contract{
		Address: utils.OntIDContractAddress,
		Method:  "regIDWithAttributes",
		Args:    nil,
	}

	ok, _ := InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed without arguments")
		res = false
	}

	c.Args = []interface{}{RegIDAttrParam{
		"",
		pub,
		[]Attribute{
			Attribute{
				[]byte("attr0"),
				[]byte{0},
				[]byte{1},
			},
		},
	}}
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed with invalid id")
		res = false
	}

	c.Args = []interface{}{RegIDAttrParam{
		test_id,
		nil,
		[]Attribute{
			Attribute{
				[]byte("attr0"),
				[]byte{0},
				[]byte{1},
			},
		},
	}}
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed with invalid key")
		res = false
	}

	c.Args = []interface{}{RegIDAttrParam{
		test_id,
		pub,
		[]Attribute{
			Attribute{
				nil,
				[]byte{0},
				[]byte{1},
			},
		},
	}}
	ok, _ = InvokeContract(ctx, c, false)
	if ok {
		ctx.LogError("register should failed with invalid attributes")
		res = false
	}

	c.Args = []interface{}{RegIDAttrParam{
		test_id,
		pub,
		[]Attribute{
			Attribute{
				[]byte("attr0"),
				[]byte{0},
				[]byte{1},
			},
			Attribute{
				[]byte("attr1"),
				[]byte{1},
				[]byte{0, 1},
			},
		},
	}}
	ok, _ = InvokeContract(ctx, c, false)
	return ok && res
}

type AddAttrParam struct {
	ID   string
	Attr []Attribute
	Key  []byte
}

type RmAttrParam struct {
	ID      string
	AttrKey []byte
	Key     []byte
}

func testAttr(ctx *testframework.TestFrameworkContext) bool {
	user, _ := ctx.GetDefaultAccount()
	pub := keypair.SerializePublicKey(user.PublicKey)

	args := []interface{}{AddAttrParam{
		test_id,
		[]Attribute{
			Attribute{
				[]byte("attr1"),
				[]byte{1},
				[]byte{0x01, 0x02},
			},
			Attribute{
				[]byte("attr2"),
				[]byte{2},
				[]byte("abcd"),
			},
		},
		pub,
	}}

	c := &Contract{
		Address: utils.OntIDContractAddress,
		Method:  "addAttributes",
		Args:    args,
	}

	ok, _ := InvokeContract(ctx, c, false)
	if !ok {
		return false
	}

	c.Method = "removeAttribute"
	c.Args = []interface{}{RmAttrParam{
		test_id,
		[]byte("attr2"),
		pub,
	}}
	ok, _ = InvokeContract(ctx, c, false)

	return ok
}
