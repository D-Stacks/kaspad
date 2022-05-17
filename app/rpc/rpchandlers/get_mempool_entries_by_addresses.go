package rpchandlers

import (
	"errors"

	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/app/rpc/rpccontext"
	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/kaspanet/kaspad/domain/consensus/utils/txscript"
	"github.com/kaspanet/kaspad/infrastructure/network/netadapter/router"
	"github.com/kaspanet/kaspad/util"
)

// HandleGetMempoolEntriesByAddresses handles the respectively named RPC command
func HandleGetMempoolEntriesByAddresses(context *rpccontext.Context, _ *router.Router, request appmessage.Message) (appmessage.Message, error) {

	getMempoolEntriesByAddressesRequest := request.(*appmessage.GetMempoolEntriesByAddressesRequestMessage)
	var mempoolEntriesByAddresses []*appmessage.MempoolEntryByAddress

	if getMempoolEntriesByAddressesRequest.IncludeTransactionPool && getMempoolEntriesByAddressesRequest.IncludeOrphanPool { //both true

		transactionPoolTransactions := context.Domain.MiningManager().AllTransactions()
		transactionPoolEntriesByAddresses, err := extractMempoolEntriesByAddressesFromTransactions(
			context,
			getMempoolEntriesByAddressesRequest.Addresses,
			transactionPoolTransactions,
			false,
		)
		if err != nil {
			rpcError := &appmessage.RPCError{}
			if !errors.As(err, &rpcError) {
				return nil, err
			}
			errorMessage := &appmessage.GetUTXOsByAddressesResponseMessage{}
			errorMessage.Error = rpcError
			return errorMessage, nil
		}

		orphanPoolTransactions := context.Domain.MiningManager().AllOrphanTransactions()
		orphanPoolEntriesByAddresse, err := extractMempoolEntriesByAddressesFromTransactions(
			context,
			getMempoolEntriesByAddressesRequest.Addresses,
			orphanPoolTransactions,
			true,
		)
		if err != nil {
			rpcError := &appmessage.RPCError{}
			if !errors.As(err, &rpcError) {
				return nil, err
			}
			errorMessage := &appmessage.GetUTXOsByAddressesResponseMessage{}
			errorMessage.Error = rpcError
			return errorMessage, nil
		}

		mempoolEntriesByAddresses = append(transactionPoolEntriesByAddresses, orphanPoolEntriesByAddresse...)

	} else if getMempoolEntriesByAddressesRequest.IncludeTransactionPool && !(getMempoolEntriesByAddressesRequest.IncludeOrphanPool) { //only transactions
		transactionPoolTransactions := context.Domain.MiningManager().AllTransactions()
		transactionPoolEntriesByAddresses, err := extractMempoolEntriesByAddressesFromTransactions(
			context,
			getMempoolEntriesByAddressesRequest.Addresses,
			transactionPoolTransactions,
			true,
		)
		if err != nil {
			return nil, err
		}

		mempoolEntriesByAddresses = transactionPoolEntriesByAddresses

	} else if !(getMempoolEntriesByAddressesRequest.IncludeTransactionPool) && getMempoolEntriesByAddressesRequest.IncludeOrphanPool { //only orphans
		orphanPoolTransactions := context.Domain.MiningManager().AllOrphanTransactions()
		orphanPoolEntriesByAddresse, err := extractMempoolEntriesByAddressesFromTransactions(
			context,
			getMempoolEntriesByAddressesRequest.Addresses,
			orphanPoolTransactions,
			true,
		)
		if err != nil {
			rpcError := &appmessage.RPCError{}
			if !errors.As(err, &rpcError) {
				return nil, err
			}
			errorMessage := &appmessage.GetUTXOsByAddressesResponseMessage{}
			errorMessage.Error = rpcError
			return errorMessage, nil
		}

		mempoolEntriesByAddresses = orphanPoolEntriesByAddresse

	} else if !(getMempoolEntriesByAddressesRequest.IncludeTransactionPool || getMempoolEntriesByAddressesRequest.IncludeOrphanPool) {
		errorMessage := &appmessage.GetMempoolEntryResponseMessage{}
		errorMessage.Error = appmessage.RPCErrorf("Request is not querying any mempool pools")
		return errorMessage, nil
	}

	return appmessage.NewGetMempoolEntriesByAddressesResponseMessage(mempoolEntriesByAddresses), nil
}

//TO DO: optimize extractMempoolEntriesByAddressesFromTransactions
func extractMempoolEntriesByAddressesFromTransactions(context *rpccontext.Context, addresses []string, transactions []*externalapi.DomainTransaction, areOrphans bool) ([]*appmessage.MempoolEntryByAddress, error) {
	mempoolEntriesByAddresses := make([]*appmessage.MempoolEntryByAddress, 0)
	for _, addressString := range addresses {
		_, err := util.DecodeAddress(addressString, context.Config.ActiveNetParams.Prefix)
		if err != nil {
			return nil, appmessage.RPCErrorf("Could not decode address '%s': %s", addressString, err)
		}

		sending := make([]*appmessage.MempoolEntry, 0)
		receiving := make([]*appmessage.MempoolEntry, 0)

		for _, transaction := range transactions {

			for _, input := range transaction.Inputs {
				_, transactionSendingAddress, err := txscript.ExtractScriptPubKeyAddress(
					input.UTXOEntry.ScriptPublicKey(),
					context.Config.ActiveNetParams)
				if err != nil {
					return nil, err
				}
				if addressString == transactionSendingAddress.String() {
					rpcTransaction := appmessage.DomainTransactionToRPCTransaction(transaction)
					sending = append(
						sending,
						&appmessage.MempoolEntry{
							Fee:         transaction.Fee,
							Transaction: rpcTransaction,
							IsOrphan:    areOrphans,
						},
					)
					break //one input is enough
				}
			}

			for _, output := range transaction.Outputs {
				_, transactionReceivingAddress, err := txscript.ExtractScriptPubKeyAddress(
					output.ScriptPublicKey,
					context.Config.ActiveNetParams,
				)
				if err != nil {
					return nil, err
				}
				if addressString == transactionReceivingAddress.String() {
					rpcTransaction := appmessage.DomainTransactionToRPCTransaction(transaction)
					receiving = append(
						receiving,
						&appmessage.MempoolEntry{
							Fee:         transaction.Fee,
							Transaction: rpcTransaction,
							IsOrphan:    areOrphans,
						},
					)
					break //one output is enough
				}
			}

			//Only append mempoolEntriesByAddress, if at least 1 mempoolEntry for the address is found.
			//This mimics the behaviour of GetUtxosByAddresses RPC call.
			if len(sending) > 0 || len(receiving) > 0 {
				mempoolEntriesByAddresses = append(
					mempoolEntriesByAddresses,
					&appmessage.MempoolEntryByAddress{
						Address:   addressString,
						Sending:   sending,
						Receiving: receiving,
					},
				)
			}
		}

	}
	return mempoolEntriesByAddresses, nil
}
