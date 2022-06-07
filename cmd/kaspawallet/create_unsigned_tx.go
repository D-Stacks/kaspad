package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kaspanet/kaspad/cmd/kaspawallet/daemon/client"
	"github.com/kaspanet/kaspad/cmd/kaspawallet/daemon/pb"
	"github.com/kaspanet/kaspad/domain/consensus/utils/constants"
)

func createUnsignedTransaction(conf *createUnsignedTransactionConfig) error {
	daemonClient, tearDown, uuid, err := client.Connect(conf.DaemonAddress)
	if err != nil {
		return err
	}
	defer tearDown()

	ctx, cancel := context.WithTimeout(context.Background(), daemonTimeout)
	defer cancel()

	sendAmountSompi := uint64(conf.SendAmount * constants.SompiPerKaspa)
	response, err := daemonClient.CreateUnsignedTransactions(ctx, &pb.CreateUnsignedTransactionsRequest{
		Id:      uuid.String(),
		From:    conf.FromAddresses,
		Address: conf.ToAddress,
		Amount:  sendAmountSompi,
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Created unsigned transaction")
	fmt.Println(encodeTransactionsToHex(response.UnsignedTransactions))

	return nil
}
