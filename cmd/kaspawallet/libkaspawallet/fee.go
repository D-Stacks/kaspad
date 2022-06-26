package libkaspawallet

import "github.com/kaspanet/kaspad/domain/consensus/model/externalapi"


// TODO: Implement a better fee estimation mechanism

//FeePerIntput is the current constant per input to pay for transactions.
const FeePerInput uint64 = 10000 

func CalculateFees(transactions []*externalapi.DomainTransaction) uint64 {
	var totalFee uint64 
	for _, tx := range transactions {
		totalFee += CalculateFee(tx)
	}

	return totalFee
}

func CalculateFee(transaction *externalapi.DomainTransaction) uint64 {
	return uint64(len(transaction.Inputs)) * FeePerInput
}