package nodestats

import (
	"sync"
	"time"

	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/util"
)

type nodestats struct {
	startTime time.Time

	totalP2PMessagesSent 		uint64
	totalRPCMessagesRecived 	uint64
	p2pMessagesSent      		[appmessage.MessageCommand]uint64
	p2pMessagesRecived   		[appmessage.MessageCommand]uint64
	p2pBytesSent	     		uint64
	p2pBytesRecived	     		uint64

	totalRPCMessagesSent 		uint64
	totalRPCMessagesRecived 	uint64
	rpcMessagesSent      		[appmessage.MessageCommand]uint64
	rpcMessagesRecived   		[appmessage.MessageCommand]uint64
	rpcBytesSent	     		uint64
	rpcBytesRecived			uint64

	blocksSubmitted          	uint64
	blcoksSubmittedByAddress 	[util.Address]uint64
	blocksAccepted		 	uint64
	blocksAccptedByAddress	 	[util.Address]uint64

	blocksProcessed			uint64
	transactionsProcessed		uint64

	blocksAddedToMempool		uint64
	transactionsAddedToMempool	uint64
	
	totalDatabaseSizeInBytes        uint64
}

type 

New()