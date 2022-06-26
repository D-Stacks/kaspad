package server

import (
	"context"
	"strconv"

	"github.com/kaspanet/kaspad/cmd/kaspawallet/daemon/pb"
)

func (s *server) Send(_ context.Context, request *pb.SendRequest) (*pb.SendResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	amount, err := strconv.ParseUint(request.Amount, 10, 64)
	if err != nil {
		return nil, err
	}

	unsignedTransactions, err := s.createUnsignedTransactions(request.ToAddress, amount, request.From)
	if err != nil {
		return nil, err
	}

	signedTransactions, err := s.signTransactions(unsignedTransactions, request.Password)
	if err != nil {
		return nil, err
	}

	txIDs, feePaid, err := s.broadcast(signedTransactions, false)
	if err != nil {
		return nil, err
	}

	return &pb.SendResponse{TxIDs: txIDs, Fees: feePaid}, nil
}
