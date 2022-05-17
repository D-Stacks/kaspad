package rpcclient

import "github.com/kaspanet/kaspad/app/appmessage"

// GetMempoolEntry sends an RPC request respective to the function's name and returns the RPC server's response
func (c *RPCClient) GetMempoolEntry(txID string, inludeOrphanPool bool, includeTransactionPool bool) (*appmessage.GetMempoolEntryResponseMessage, error) {
	err := c.rpcRouter.outgoingRoute().Enqueue(appmessage.NewGetMempoolEntryRequestMessage(txID, inludeOrphanPool, includeTransactionPool))
	if err != nil {
		return nil, err
	}
	response, err := c.route(appmessage.CmdGetMempoolEntryResponseMessage).DequeueWithTimeout(c.timeout)
	if err != nil {
		return nil, err
	}
	getMempoolEntryResponse := response.(*appmessage.GetMempoolEntryResponseMessage)
	if getMempoolEntryResponse.Error != nil {
		return nil, c.convertRPCError(getMempoolEntryResponse.Error)
	}
	return getMempoolEntryResponse, nil
}
