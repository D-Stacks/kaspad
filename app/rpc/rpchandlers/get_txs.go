package rpchandlers

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/app/rpc/rpccontext"
	"github.com/kaspanet/kaspad/infrastructure/network/netadapter/router"
)

// HandleGetTxs handles the respectively named RPC command
func HandleGetTxs(context *rpccontext.Context, _ *router.Router, request appmessage.Message) (appmessage.Message, error) {
	return nil, nil
}
