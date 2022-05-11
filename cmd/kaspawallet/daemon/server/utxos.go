package server

import (
	"sort"

	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
)

type walletUTXOSet map[externalapi.DomainOutpoint]*walletUTXO

func (s *server) utxosSortedByAmount() []*walletUTXO {
	utxos := make([]*walletUTXO, 0)
	for _, utxo := range s.utxoSet {
		utxos = append(utxos, utxo)
	}

	sort.Slice(utxos, func(i, j int) bool { return utxos[i].UTXOEntry.Amount() > utxos[j].UTXOEntry.Amount() })
	return utxos
}
