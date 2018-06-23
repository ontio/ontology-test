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
)

func TestGovernanceContract() {
	testframework.TFramework.RegTestCase("SimulateVoteForPeerAndWithdraw", SimulateVoteForPeerAndWithdraw)
	testframework.TFramework.RegTestCase("SimulateRejectCandidate", SimulateRejectCandidate)
	testframework.TFramework.RegTestCase("SimulateUnConsensusToConsensus", SimulateUnConsensusToConsensus)
	testframework.TFramework.RegTestCase("SimulateUnConsensusToUnConsensus", SimulateUnConsensusToUnConsensus)
	testframework.TFramework.RegTestCase("SimulateConsensusToUnConsensus", SimulateConsensusToUnConsensus)
	testframework.TFramework.RegTestCase("SimulateConsensusToConsensus", SimulateConsensusToConsensus)
	testframework.TFramework.RegTestCase("SimulateQuitUnConsensus", SimulateQuitUnConsensus)
	testframework.TFramework.RegTestCase("SimulateQuitConsensus", SimulateQuitConsensus)
	testframework.TFramework.RegTestCase("SimulateBlackConsensusAndWhite", SimulateBlackConsensusAndWhite)
	testframework.TFramework.RegTestCase("SimulateBlackUnConsensusAndWhite", SimulateBlackUnConsensusAndWhite)
	testframework.TFramework.RegTestCase("SimulateUpdateConfig", SimulateUpdateConfig)
	testframework.TFramework.RegTestCase("SimulateUpdateGlobalParam", SimulateUpdateGlobalParam)
	testframework.TFramework.RegTestCase("SimulateUpdateSplitCurve", SimulateUpdateSplitCurve)
	testframework.TFramework.RegTestCase("SimulateCommitDPosAuth", SimulateCommitDPosAuth)
	testframework.TFramework.RegTestCase("SimulateTransferPenalty", SimulateTransferPenalty)
	testframework.TFramework.RegTestCase("SimulateOntIDAndAuth", SimulateOntIDAndAuth)
	testframework.TFramework.RegTestCase("SimulateUnRegisterCandidate", SimulateUnRegisterCandidate)
}

func TestGovernanceContractError() {
	testframework.TFramework.RegTestCase("SimulateUnConsensusVoteForPeerError", SimulateUnConsensusVoteForPeerError)
	testframework.TFramework.RegTestCase("SimulateConsensusVoteForPeerError", SimulateConsensusVoteForPeerError)
	testframework.TFramework.RegTestCase("SimulateWithDrawError", SimulateWithDrawError)
	testframework.TFramework.RegTestCase("SimulateRegisterCandidateError", SimulateRegisterCandidateError)
	testframework.TFramework.RegTestCase("SimulateRejectCandidateError", SimulateRejectCandidateError)
	testframework.TFramework.RegTestCase("SimulateApproveCandidateError", SimulateApproveCandidateError)
	testframework.TFramework.RegTestCase("SimulateBlackNodeError", SimulateBlackNodeError)
	testframework.TFramework.RegTestCase("SimulateWhiteNodeError", SimulateWhiteNodeError)
	testframework.TFramework.RegTestCase("SimulateQuitNodeError", SimulateQuitNodeError)
	testframework.TFramework.RegTestCase("SimulateUpdateConfigError", SimulateUpdateConfigError)
	testframework.TFramework.RegTestCase("SimulateUpdateGlobalParamError", SimulateUpdateGlobalParamError)
	testframework.TFramework.RegTestCase("SimulateTransferPenaltyError", SimulateTransferPenaltyError)
	testframework.TFramework.RegTestCase("SimulateUnRegisterCandidateError", SimulateUnRegisterCandidateError)
}

func TestGovernanceMethods() {
	testframework.TFramework.RegTestCase("RegIdWithPublicKey", RegIdWithPublicKey)
	testframework.TFramework.RegTestCase("AssignFuncsToRole", AssignFuncsToRole)
	testframework.TFramework.RegTestCase("AssignOntIDsToRole", AssignOntIDsToRole)
	testframework.TFramework.RegTestCase("RegisterCandidate", RegisterCandidate)
	testframework.TFramework.RegTestCase("ApproveCandidate", ApproveCandidate)
	testframework.TFramework.RegTestCase("RejectCandidate", RejectCandidate)
	testframework.TFramework.RegTestCase("VoteForPeer", VoteForPeer)
	testframework.TFramework.RegTestCase("UnVoteForPeer", UnVoteForPeer)
	testframework.TFramework.RegTestCase("Withdraw", Withdraw)
	testframework.TFramework.RegTestCase("QuitNode", QuitNode)
	testframework.TFramework.RegTestCase("BlackNode", BlackNode)
	testframework.TFramework.RegTestCase("WhiteNode", WhiteNode)
	testframework.TFramework.RegTestCase("CommitDpos", CommitDpos)
	testframework.TFramework.RegTestCase("CallSplit", CallSplit)
	testframework.TFramework.RegTestCase("UpdateConfig", UpdateConfig)
	testframework.TFramework.RegTestCase("UpdateGlobalParam", UpdateGlobalParam)
	testframework.TFramework.RegTestCase("UpdateSplitCurve", UpdateSplitCurve)
	testframework.TFramework.RegTestCase("TransferPenalty", TransferPenalty)
	testframework.TFramework.RegTestCase("GetVbftConfig", GetVbftConfig)
	testframework.TFramework.RegTestCase("GetGlobalParam", GetGlobalParam)
	testframework.TFramework.RegTestCase("GetSplitCurve", GetSplitCurve)
	testframework.TFramework.RegTestCase("GetGovernanceView", GetGovernanceView)
	testframework.TFramework.RegTestCase("GetPeerPoolItem", GetPeerPoolItem)
	testframework.TFramework.RegTestCase("GetVoteInfo", GetVoteInfo)
	testframework.TFramework.RegTestCase("GetTotalStake", GetTotalStake)
	testframework.TFramework.RegTestCase("GetPenaltyStake", GetPenaltyStake)
	testframework.TFramework.RegTestCase("InBlackList", InBlackList)
	testframework.TFramework.RegTestCase("WithdrawOng", WithdrawOng)
}
