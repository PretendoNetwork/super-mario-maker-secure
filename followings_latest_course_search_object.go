package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func followingsLatestCourseSearchObject(err error, client *nex.Client, callID uint32, dataStoreSearchParam *nexproto.DataStoreSearchParam, extraData []string) {
	pRankingResults := make([]*nexproto.DataStoreCustomRankingResult, 0)

	for _, pid := range dataStoreSearchParam.OwnerIds {
		courseMetadatas := getCourseMetadatasByPID(pid)

		for _, courseMetadata := range courseMetadatas {
			pRankingResults = append(pRankingResults, courseMetadataToDataStoreCustomRankingResult(courseMetadata))
		}
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodFollowingsLatestCourseSearchObject, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
