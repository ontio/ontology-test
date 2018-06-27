package governance_feeSplit

import (
	"bytes"
	"encoding/hex"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/errors"
	"github.com/ontio/ontology/smartcontract/service/native/auth"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	cstates "github.com/ontio/ontology/smartcontract/states"
)

var OntIDVersion = byte(0)

func registerCandidate(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkey string, initPos uint64) bool {
	params := &governance.RegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
		InitPos:    initPos,
		Caller:     []byte("did:ont:" + user.Address.ToBase58()),
		KeyNo:      1,
	}
	method := "registerCandidate"
	contractAddress := utils.GovernanceContractAddress
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func unRegisterCandidate(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkey string) bool {
	params := &governance.UnRegisterCandidateParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
	}
	method := "unRegisterCandidate"
	contractAddress := utils.GovernanceContractAddress
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func approveCandidate(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkey string) bool {
	params := &governance.ApproveCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "approveCandidate"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func rejectCandidate(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkey string) bool {
	params := &governance.RejectCandidateParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "rejectCandidate"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func voteForPeer(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkeyList []string, posList []uint64) bool {
	params := &governance.VoteForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "voteForPeer"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func unVoteForPeer(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkeyList []string, posList []uint64) bool {
	params := &governance.VoteForPeerParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		PosList:        posList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "unVoteForPeer"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func withdraw(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkeyList []string, withdrawList []uint64) bool {
	params := &governance.WithdrawParam{
		Address:        user.Address,
		PeerPubkeyList: peerPubkeyList,
		WithdrawList:   withdrawList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdraw"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func withdrawOng(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	params := &governance.WithdrawOngParam{
		Address:        user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "withdrawOng"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func commitDpos(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "commitDpos"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func quitNode(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkey string) bool {
	params := &governance.QuitNodeParam{
		PeerPubkey: peerPubkey,
		Address:    user.Address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "quitNode"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func blackNode(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkeyList []string) bool {
	params := &governance.BlackNodeParam{
		PeerPubkeyList: peerPubkeyList,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "blackNode"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func whiteNode(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkey string) bool {
	params := &governance.WhiteNodeParam{
		PeerPubkey: peerPubkey,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "whiteNode"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func updateConfig(ctx *testframework.TestFrameworkContext, user *account.Account, config *governance.Configuration) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateConfig"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{config})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func updateGlobalParam(ctx *testframework.TestFrameworkContext, user *account.Account, globalParam *governance.GlobalParam) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateGlobalParam"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{globalParam})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func updateSplitCurve(ctx *testframework.TestFrameworkContext, user *account.Account, splitCurve *governance.SplitCurve) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "updateSplitCurve"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{splitCurve})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func callSplit(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	contractAddress := utils.GovernanceContractAddress
	method := "callSplit"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func transferPenalty(ctx *testframework.TestFrameworkContext, user *account.Account, peerPubkey string, address common.Address) bool {
	params := &governance.TransferPenaltyParam{
		PeerPubkey: peerPubkey,
		Address:    address,
	}
	contractAddress := utils.GovernanceContractAddress
	method := "transferPenalty"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func assignFuncsToRole(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	params := &auth.FuncsToRoleParam{
		ContractAddr: utils.GovernanceContractAddress,
		AdminOntID:   []byte("did:ont:" + user.Address.ToBase58()),
		Role:         []byte("role"),
		FuncNames:    []string{"registerCandidate"},
		KeyNo:        1,
	}
	method := "assignFuncsToRole"
	contractAddress := utils.AuthContractAddress
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func assignOntIDsToRole(ctx *testframework.TestFrameworkContext, user *account.Account, users []*account.Account) bool {
	params := &auth.OntIDsToRoleParam{
		ContractAddr: utils.GovernanceContractAddress,
		AdminOntID:   []byte("did:ont:" + user.Address.ToBase58()),
		Role:         []byte("role"),
		Persons:      [][]byte{},
		KeyNo:        1,
	}
	for _, user := range users {
		params.Persons = append(params.Persons, []byte("did:ont:"+user.Address.ToBase58()))
	}
	contractAddress := utils.AuthContractAddress
	method := "assignOntIDsToRole"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

func verifyToken(ctx *testframework.TestFrameworkContext, user *account.Account, caller []byte, fn string) bool {
	params := &auth.VerifyTokenParam{
		ContractAddr: utils.GovernanceContractAddress,
		Caller:       caller,
		Fn:           fn,
		KeyNo:        1,
	}
	contractAddress := utils.AuthContractAddress
	method := "verifyToken"
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}

type RegIDWithPublicKeyParam struct {
	OntID  []byte
	Pubkey []byte
}

func regIdWithPublicKey(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	params := RegIDWithPublicKeyParam{
		OntID:  []byte("did:ont:" + user.Address.ToBase58()),
		Pubkey: keypair.SerializePublicKey(user.PublicKey),
	}
	method := "regIDWithPublicKey"
	contractAddress := utils.OntIDContractAddress
	_, err := ctx.Ont.Rpc.InvokeNativeContract(ctx.GetGasPrice(), ctx.GetGasLimit(), user, OntIDVersion,
		contractAddress, method, []interface{}{params})
	if err != nil {
		ctx.LogError("invokeNativeContract error")
	}

	return true
}

func getVbftConfig(ctx *testframework.TestFrameworkContext) (*governance.Configuration, error) {
	contractAddress := utils.GovernanceContractAddress
	config := new(governance.Configuration)
	key := []byte(governance.VBFT_CONFIG)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := config.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize config error!")
	}
	return config, nil
}

func getGlobalParam(ctx *testframework.TestFrameworkContext) (*governance.GlobalParam, error) {
	contractAddress := utils.GovernanceContractAddress
	globalParam := new(governance.GlobalParam)
	key := []byte(governance.GLOBAL_PARAM)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := globalParam.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize globalParam error!")
	}
	return globalParam, nil
}

func getSplitCurve(ctx *testframework.TestFrameworkContext) (*governance.SplitCurve, error) {
	contractAddress := utils.GovernanceContractAddress
	splitCurve := new(governance.SplitCurve)
	key := []byte(governance.SPLIT_CURVE)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := splitCurve.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize splitCurve error!")
	}
	return splitCurve, nil
}

func getGovernanceView(ctx *testframework.TestFrameworkContext) (*governance.GovernanceView, error) {
	contractAddress := utils.GovernanceContractAddress
	governanceView := new(governance.GovernanceView)
	key := []byte(governance.GOVERNANCE_VIEW)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := governanceView.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize governanceView error!")
	}
	return governanceView, nil
}

func getView(ctx *testframework.TestFrameworkContext) (uint32, error) {
	governanceView, err := getGovernanceView(ctx)
	if err != nil {
		return 0, errors.NewDetailErr(err, errors.ErrNoCode, "getGovernanceView error")
	}
	return governanceView.View, nil
}

func getPeerPoolMap(ctx *testframework.TestFrameworkContext) (*governance.PeerPoolMap, error) {
	contractAddress := utils.GovernanceContractAddress
	view, err := getView(ctx)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getView error")
	}
	peerPoolMap := &governance.PeerPoolMap{
		PeerPoolMap: make(map[string]*governance.PeerPoolItem),
	}
	viewBytes, err := governance.GetUint32Bytes(view)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "GetUint32Bytes, get viewBytes error!")
	}
	key := ConcatKey([]byte(governance.PEER_POOL), viewBytes)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := peerPoolMap.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize peerPoolMap error!")
	}
	return peerPoolMap, nil
}

func getVoteInfo(ctx *testframework.TestFrameworkContext, peerPubkey string, address common.Address) (*governance.VoteInfo, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	voteInfo := new(governance.VoteInfo)
	key := ConcatKey([]byte(governance.VOTE_INFO_POOL), peerPubkeyPrefix, address[:])
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := voteInfo.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize voteInfo error!")
	}
	return voteInfo, nil
}

func inBlackList(ctx *testframework.TestFrameworkContext, peerPubkey string) (bool, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	key := ConcatKey([]byte(governance.BLACK_LIST), peerPubkeyPrefix)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if len(value) != 0 {
		return true, nil
	}
	return false, nil
}

func getTotalStake(ctx *testframework.TestFrameworkContext, address common.Address) (*governance.TotalStake, error) {
	contractAddress := utils.GovernanceContractAddress
	totalStake := new(governance.TotalStake)
	key := ConcatKey([]byte(governance.TOTAL_STAKE), address[:])
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := totalStake.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize totalStake error!")
	}
	return totalStake, nil
}

func getPenaltyStake(ctx *testframework.TestFrameworkContext, peerPubkey string) (*governance.PenaltyStake, error) {
	contractAddress := utils.GovernanceContractAddress
	peerPubkeyPrefix, err := hex.DecodeString(peerPubkey)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "hex.DecodeString, peerPubkey format error!")
	}
	penaltyStake := new(governance.PenaltyStake)
	key := ConcatKey([]byte(governance.PENALTY_STAKE), peerPubkeyPrefix)
	value, err := ctx.Ont.Rpc.GetStorage(contractAddress, key)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "getStorage error")
	}
	if err := penaltyStake.Deserialize(bytes.NewBuffer(value)); err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNoCode, "deserialize, deserialize penaltyStake error!")
	}
	return penaltyStake, nil
}

func getDDO(ctx *testframework.TestFrameworkContext, user *account.Account) bool {
	contractAddress := utils.OntIDContractAddress

	bf := new(bytes.Buffer)
	if err := serialization.WriteVarBytes(bf, []byte("did:ont:"+user.Address.ToBase58())); err != nil {
		ctx.LogError("Serialize ontid error:%s", err)
		return false
	}
	crt := &cstates.Contract{
		Address: contractAddress,
		Method:  "getDDO",
		Args:    bf.Bytes(),
	}

	ok := invokeNativeContractWithoutWait(ctx, crt, user)
	if !ok {
		ctx.LogError("invokeNativeContract error")
	}
	return true
}
