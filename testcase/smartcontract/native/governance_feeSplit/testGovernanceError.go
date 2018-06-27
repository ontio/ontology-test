/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package governance_feeSplit

import (
	"github.com/ontio/ontology-test/testframework"
	"github.com/ontio/ontology/smartcontract/service/native/governance"
)

const (
	PEER_PUBKEY2 = "028f5cbc5f878cddda2be940eeb3643301d798f4bcac71ef81419e4c082045a3ec"
	PEER_PUBKEY3 = "0336f107dde5e8f5844bb69d4fdcb3b8a324be3e27f7350cb6172bb54340cb3309"
)

func SimulateUnConsensusVoteForPeerError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint64{1000}
	unVoteForPeer(ctx, user1, peerPubkeyList, posList)
	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint64{2000}
	voteForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user2, INIT_ONT-2000)
	if !ok {
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint64{1000}
	unVoteForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check voteInfo data
	//user2
	voteInfo, err := getVoteInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getVoteInfo error :%v", err)
	}
	if voteInfo.NewPos != 0 || voteInfo.FreezePos != 1000 || voteInfo.ConsensusPos != 0 || voteInfo.WithdrawUnfreezePos != 0 ||
		voteInfo.WithdrawPos != 0 || voteInfo.WithdrawFreezePos != 1000 {
		ctx.LogError("voteInfo data for user2 is wrong!")
		return false
	}

	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint64{2000}
	unVoteForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check voteInfo data
	//user2
	voteInfo, err = getVoteInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getVoteInfo error :%v", err)
	}
	if voteInfo.NewPos != 0 || voteInfo.FreezePos != 1000 || voteInfo.ConsensusPos != 0 || voteInfo.WithdrawUnfreezePos != 0 ||
		voteInfo.WithdrawPos != 0 || voteInfo.WithdrawFreezePos != 1000 {
		ctx.LogError("voteInfo data for user2 is wrong!")
		return false
	}
	return true
}

func SimulateConsensusVoteForPeerError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	user2, ok := getAccount2(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint64{1000}
	unVoteForPeer(ctx, user1, peerPubkeyList, posList)
	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint64{300000}
	voteForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT)
	if !ok {
		return false
	}
	ok = checkBalance(ctx, user2, INIT_ONT-300000)
	if !ok {
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint64{200000}
	unVoteForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check voteInfo data
	//user2
	voteInfo, err := getVoteInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getVoteInfo error :%v", err)
	}
	if voteInfo.NewPos != 0 || voteInfo.FreezePos != 0 || voteInfo.ConsensusPos != 100000 || voteInfo.WithdrawUnfreezePos != 0 ||
		voteInfo.WithdrawPos != 200000 || voteInfo.WithdrawFreezePos != 0 {
		ctx.LogError("voteInfo data for user2 is wrong!")
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint64{200000}
	unVoteForPeer(ctx, user2, peerPubkeyList, posList)
	waitForBlock(ctx)
	//check voteInfo data
	//user2
	voteInfo, err = getVoteInfo(ctx, PEER_PUBKEY, user2.Address)
	if err != nil {
		ctx.LogError("getVoteInfo error :%v", err)
	}
	if voteInfo.NewPos != 0 || voteInfo.FreezePos != 0 || voteInfo.ConsensusPos != 100000 || voteInfo.WithdrawUnfreezePos != 0 ||
		voteInfo.WithdrawPos != 0 || voteInfo.WithdrawFreezePos != 200000 {
		ctx.LogError("voteInfo data for user2 is wrong!")
		return false
	}
	return true
}

func SimulateWithDrawError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint64{300000}
	voteForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)
	peerPubkeyList = []string{PEER_PUBKEY}
	posList = []uint64{100000}
	unVoteForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)

	withdrawList := []uint64{80000}
	withdraw(ctx, user1, peerPubkeyList, withdrawList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-300000)
	if !ok {
		return false
	}
	commitDpos(ctx, user)
	waitForBlock(ctx)
	withdrawList = []uint64{80000}
	withdraw(ctx, user1, peerPubkeyList, withdrawList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-300000)
	if !ok {
		return false
	}

	commitDpos(ctx, user)
	waitForBlock(ctx)
	withdrawList = []uint64{80000}
	withdraw(ctx, user1, peerPubkeyList, withdrawList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-220000)
	if !ok {
		return false
	}

	withdrawList = []uint64{80000}
	withdraw(ctx, user1, peerPubkeyList, withdrawList)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-220000)
	if !ok {
		return false
	}
	return true
}

func SimulateRegisterCandidateError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	registerCandidate(ctx, user, PEER_PUBKEY, 10000)
	registerCandidate(ctx, user1, PEER_PUBKEY2, 1000)
	waitForBlock(ctx)
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.CandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY2].Status != governance.RegisterCandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	approveCandidate(ctx, user, PEER_PUBKEY2)
	waitForBlock(ctx)
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY2].Status != governance.RegisterCandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	return true
}

func SimulateApproveCandidateError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	registerCandidate(ctx, user, PEER_PUBKEY2, 10000)
	waitForBlock(ctx)
	approveCandidate(ctx, user1, PEER_PUBKEY2)
	waitForBlock(ctx)
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY2].Status != governance.RegisterCandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	approveCandidate(ctx, user, PEER_PUBKEY)
	waitForBlock(ctx)
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.CandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	approveCandidate(ctx, user, PEER_PUBKEY3)
	waitForBlock(ctx)
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY3]
	if ok {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	return true
}

func SimulateRejectCandidateError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	registerCandidate(ctx, user1, PEER_PUBKEY2, 10000)
	waitForBlock(ctx)
	rejectCandidate(ctx, user1, PEER_PUBKEY2)
	waitForBlock(ctx)
	//check pperPoolMap
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	_, ok = peerPoolMap.PeerPoolMap[PEER_PUBKEY2]
	if !ok {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY2].InitPos != 10000 ||
		peerPoolMap.PeerPoolMap[PEER_PUBKEY2].Status != governance.RegisterCandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	return true
}

func SimulateBlackNodeError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	blackNode(ctx, user1, []string{PEER_PUBKEY})
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.CandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	ok = blackNode(ctx, user, []string{PEER_PUBKEY})
	if !ok {
		return false
	}
	return true
}

func SimulateWhiteNodeError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}
	blackNode(ctx, user, []string{PEER_PUBKEY})
	waitForBlock(ctx)
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.BlackStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	whiteNode(ctx, user1, PEER_PUBKEY)
	waitForBlock(ctx)
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.BlackStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	whiteNode(ctx, user, "0253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a85")
	waitForBlock(ctx)
	peerPoolMap, err = getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap["0253ccfd439b29eca0fe90ca7c6eaa1f98572a054aa2d1d56e72ad96c466107a85"].Status != governance.ConsensusStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}
	return true
}

func SimulateQuitNodeError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	quitNode(ctx, user1, PEER_PUBKEY)
	waitForBlock(ctx)
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	if peerPoolMap.PeerPoolMap[PEER_PUBKEY].Status != governance.CandidateStatus {
		ctx.LogError("peerPoolItem data is wrong!")
		return false
	}

	quitNode(ctx, user1, PEER_PUBKEY2)
	waitForBlock(ctx)
	return true
}

func SimulateUpdateConfigError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	config := &governance.Configuration{
		N:                    1,
		C:                    2,
		K:                    3,
		L:                    4,
		BlockMsgDelay:        5,
		HashMsgDelay:         6,
		PeerHandshakeTimeout: 7,
		MaxBlockChangeView:   8,
	}
	ok = updateConfig(ctx, user, config)
	if !ok {
		return false
	}
	waitForBlock(ctx)

	//check config
	config, err := getVbftConfig(ctx)
	if err != nil {
		ctx.LogError("getVbftConfig error :%v", err)
		return false
	}
	if config.L != 112 || config.K != 7 || config.C != 2 || config.N != 7 || config.BlockMsgDelay != 10000 || config.HashMsgDelay != 10000 ||
		config.PeerHandshakeTimeout != 10 || config.MaxBlockChangeView != 1000 {
		ctx.LogError("config is error")
		return false
	}
	return true
}

func SimulateUpdateGlobalParamError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	globalParam := &governance.GlobalParam{
		CandidateFee: 500000000000,
		MinInitStake: 10000,
		CandidateNum: 8 * 8,
		PosLimit:     30,
		A:            20,
		B:            30,
		Yita:         7,
		Penalty:      20,
	}
	ok = updateGlobalParam(ctx, user, globalParam)
	if !ok {
		return false
	}
	waitForBlock(ctx)
	//check config
	globalParam, err := getGlobalParam(ctx)
	if err != nil {
		ctx.LogError("getGlobalParam error :%v", err)
		return false
	}
	if globalParam.CandidateFee != 0 || globalParam.MinInitStake != 10000 ||
		globalParam.CandidateNum != (7*7) || globalParam.PosLimit != 50 ||
		globalParam.A != 50 || globalParam.B != 50 || globalParam.Yita != 5 || globalParam.Penalty != 5 {
		ctx.LogError("globalParam is error")
		return false
	}
	return true
}

func SimulateTransferPenaltyError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	//select in consensus
	peerPubkeyList := []string{PEER_PUBKEY}
	posList := []uint64{1000}
	voteForPeer(ctx, user1, peerPubkeyList, posList)
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)
	//check total stake
	totalStake, err := getTotalStake(ctx, user.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 350000+10000 {
		ctx.LogError("total stake user is error")
		return false
	}
	totalStake, err = getTotalStake(ctx, user1.Address)
	if err != nil {
		ctx.LogError("getTotalStake error :%v", err)
		return false
	}
	if totalStake.Stake != 1000 {
		ctx.LogError("total stake user1 is error")
		return false
	}

	//blacknode
	ok = blackNode(ctx, user, []string{PEER_PUBKEY})
	if !ok {
		return false
	}
	waitForBlock(ctx)
	commitDpos(ctx, user)
	waitForBlock(ctx)

	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-1000)
	if !ok {
		return false
	}
	//check penaltyStake
	penaltyStake, err := getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPenaltyStake error :%v", err)
		return false
	}
	if penaltyStake.InitPos != 10000 || penaltyStake.VotePos != 50 {
		ctx.LogError("penalty stake is error")
		return false
	}

	ok = transferPenalty(ctx, user1, PEER_PUBKEY, user1.Address)
	if !ok {
		return false
	}
	waitForBlock(ctx)

	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-1000)
	if !ok {
		return false
	}
	//check penaltyStake
	penaltyStake, err = getPenaltyStake(ctx, PEER_PUBKEY)
	if err != nil {
		ctx.LogError("getPenaltyStake error :%v", err)
		return false
	}
	if penaltyStake.InitPos != 10000 || penaltyStake.VotePos != 50 {
		ctx.LogError("penalty stake is error")
		return false
	}
	return true
}

func SimulateUnRegisterCandidateError(ctx *testframework.TestFrameworkContext) bool {
	user, ok := getDefaultAccount(ctx)
	if !ok {
		return false
	}
	user1, ok := getAccount1(ctx)
	if !ok {
		return false
	}
	ok = setupTest(ctx, user)
	if !ok {
		return false
	}

	registerCandidate(ctx, user1, PEER_PUBKEY2, 10000)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-10000)
	if !ok {
		return false
	}

	unRegisterCandidate(ctx, user, PEER_PUBKEY2)
	waitForBlock(ctx)
	//check balance
	ok = checkBalance(ctx, user1, INIT_ONT-10000)
	if !ok {
		return false
	}
	//check peerPoolItem data
	peerPoolMap, err := getPeerPoolMap(ctx)
	if err != nil {
		ctx.LogError("getPeerPoolMap error :%v", err)
		return false
	}
	peerPoolItem, ok := peerPoolMap.PeerPoolMap[PEER_PUBKEY2]
	if !ok {
		ctx.LogError("peer should exist")
		return false
	}
	if peerPoolItem.Status != governance.RegisterCandidateStatus {
		ctx.LogError("peerPoolItem status error")
		return false
	}

	return true
}
