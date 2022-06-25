package nodestats

type P2PStats struct {
	totalP2PMessagesSent 		uint64
	totalRPCMessagesRecived 	uint64
	p2pMessagesSent      		[appmessage.MessageCommand]uint64
	p2pMessagesRecived   		[appmessage.MessageCommand]uint64
	p2pBytesSent	     		uint64
	p2pBytesRecived	     		uint64
}

type RPCstats struct {
	totalRPCMessagesSent 		uint64
	totalRPCMessagesRecived 	uint64
	rpcMessagesSent      		[appmessage.MessageCommand]uint64
	rpcMessagesRecived   		[appmessage.MessageCommand]uint64
	rpcBytesSent	     		uint64
	rpcBytesRecived			uint64
}

