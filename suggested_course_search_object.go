package main

import (
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func suggestedCourseSearchObject(err error, client *nex.Client, callID uint32, param *nexproto.DataStoreSearchParam, extraData []string) {
	// TODO: complete this

	courseID, _ := strconv.ParseUint(extraData[0], 0, 64)

	if userNotOwnCourse(courseID, client.PID()) {
		incrementCourseAttemptCount(courseID) // We also know this is when a user attempts a course
	}

	pRankingResults := make([]*nexproto.DataStoreCustomRankingResult, 0)

	courseMetadatas := getCourseMetadatasByLimit(4) // In PCAPs param.minimalRatingFrequency is 4 but is 0 here?

	for _, courseMetadata := range courseMetadatas {
		pRankingResults = append(pRankingResults, courseMetadataToDataStoreCustomRankingResult(courseMetadata))
	}

	rmcResponseStream := nex.NewStreamOut(nexServer)

	rmcResponseStream.WriteListStructure(pRankingResults)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(nexproto.DataStoreSMMProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.DataStoreSMMMethodSuggestedCourseSearchObject, rmcResponseBody)

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
