package rpcclient

import "github.com/kaspanet/kaspad/app/appmessage"

// GetIncludingBlockHashOfTxs sends an RPC request respective to the function's name and returns the RPC server's response
func (c *RPCClient) GetIncludingBlockHashOfTxs(txID string) (*appmessage.GetIncludingBlockHashOfTxResponseMessage, error) {
	err := c.rpcRouter.outgoingRoute().Enqueue(appmessage.NewGetIncludingBlockHashOfTxRequest(txID))
	if err != nil {
		return nil, err
	}
	response, err := c.route(appmessage.CmdGetIncludingBlockHashOfTxResponseMessage).DequeueWithTimeout(c.timeout)
	if err != nil {
		return nil, err
	}
	getIncludingBlockHashOfTxResponse := response.(*appmessage.GetIncludingBlockHashOfTxResponseMessage)
	if getIncludingBlockHashOfTxResponse.Error != nil {
		return nil, c.convertRPCError(getIncludingBlockHashOfTxResponse.Error)
	}
	return getIncludingBlockHashOfTxResponse, nil
}
