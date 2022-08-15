package rpchandlers

import (
	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/app/rpc/rpccontext"
	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/kaspanet/kaspad/domain/consensus/utils/transactionhelper"
	"github.com/kaspanet/kaspad/domain/consensus/utils/txscript"
	"github.com/kaspanet/kaspad/infrastructure/network/netadapter/router"
	"github.com/kaspanet/kaspad/util"
	"github.com/kaspanet/kaspad/version"
)

// HandleGetBlockTemplates handles the respectively named RPC command
func HandleGetBlockTemplates(context *rpccontext.Context, _ *router.Router, request appmessage.Message) (appmessage.Message, error) {
	getBlockTemplatesRequest := request.(*appmessage.GetBlockTemplatesRequestMessage)

	templeateResponses := make([]*appmessage.TemplatesResponse, 0)
	var isNearlySynced bool
	var templateBlock *externalapi.DomainBlock
	for _, templateRequest := range getBlockTemplatesRequest.TemplatesRequest {
		payAddress, err := util.DecodeAddress(templateRequest.PayAddress, context.Config.ActiveNetParams.Prefix)
		if err != nil {
			continue // could have better handling, but do not want to exclude valids in call
		}
		scriptPublicKey, err := txscript.PayToAddrScript(payAddress)
		if err != nil {
			return nil, err
		}

		coinbaseData := &externalapi.DomainCoinbaseData{ScriptPublicKey: scriptPublicKey, ExtraData: []byte(version.Version() + "/" + templateRequest.ExtraData)}

		templateBlock, isNearlySynced, err = context.Domain.MiningManager().GetBlockTemplate(coinbaseData)
		if err != nil {
			return nil, err
		}

		if uint64(len(templateBlock.Transactions[transactionhelper.CoinbaseTransactionIndex].Payload)) > context.Config.NetParams().MaxCoinbasePayloadLength {
			errorMessage := &appmessage.GetBlockTemplateResponseMessage{}
			errorMessage.Error = appmessage.RPCErrorf("Coinbase payload is above max length (%d). Try to shorten the extra data.", context.Config.NetParams().MaxCoinbasePayloadLength)
			return errorMessage, nil
		}

		rpcBlock := appmessage.DomainBlockToRPCBlock(templateBlock)
		templeateResponses = append(
			templeateResponses,
			&appmessage.TemplatesResponse{
				Address: templateRequest.PayAddress,
				Block:   rpcBlock,
			},
		)
	}

	return appmessage.NewGetBlockTemplatesResponseMessage(templeateResponses, context.ProtocolManager.Context().HasPeers() && isNearlySynced), nil
}
