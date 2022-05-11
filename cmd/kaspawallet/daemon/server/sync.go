package server

import (
	"fmt"
	"time"

	"github.com/kaspanet/kaspad/cmd/kaspawallet/libkaspawallet"

	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/pkg/errors"
)

var keyChains = []uint8{libkaspawallet.ExternalKeychain, libkaspawallet.InternalKeychain}

type walletAddressSet map[string]*walletAddress

func (was walletAddressSet) strings() []string {
	addresses := make([]string, 0, len(was))
	for addr := range was {
		addresses = append(addresses, addr)
	}
	return addresses
}

func (s *server) initialize() error {
	err := s.collectRecentAddresses()
	if err != nil {
		return err
	}

	err = s.collectFarAddresses()
	if err != nil {
		return err
	}

	err = s.getExistingUTXOsWithLock()
	if err != nil {
		return err
	}
	return nil
}

func (s *server) onUTXOsChanged(notification *appmessage.UTXOsChangedNotificationMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()
	fmt.Println("removing", len(notification.Removed))
	for _, removedUTXOByAddressesEntry := range notification.Removed {
		removedDomainOutpoint, err := appmessage.RPCOutpointToDomainOutpoint(removedUTXOByAddressesEntry.Outpoint)
		if err != nil {
			log.Warn(err)
		}
		delete(s.utxoSet, *removedDomainOutpoint)
	}
	fmt.Println("adding", len(notification.Added))
	for _, addedUTXOByAddressesEntry := range notification.Added {
		newDomainOutpoint, err := appmessage.RPCOutpointToDomainOutpoint(addedUTXOByAddressesEntry.Outpoint)
		if err != nil {
			log.Warn(err)
		}
		address, ok := s.addressSet[addedUTXOByAddressesEntry.Address]
		if !ok {
			log.Warn(errors.Errorf("Got result from address %s even though it wasn't requested", addedUTXOByAddressesEntry.Address))
		}
		utxo, err := appmessage.RPCUTXOEntryToUTXOEntry(addedUTXOByAddressesEntry.UTXOEntry)
		if err != nil {
			log.Warn(err)
		}
		s.utxoSet[*newDomainOutpoint] = &walletUTXO{
			Outpoint:  newDomainOutpoint,
			UTXOEntry: utxo,
			address:   address,
		}
	}
}

func (s *server) sync() error {
	s.initialize()

	err := s.rpcClient.RegisterForUTXOsChangedNotifications(s.addressSet.strings(), s.onUTXOsChanged)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := s.collectRecentAddresses()
		if err != nil {
			return err
		}

		err = s.collectFarAddresses()
		if err != nil {
			return err
		}
	}

	return nil
}

const numIndexesToQuery = 100

// addressesToQuery scans the addresses in the given range. Because
// each cosigner in a multisig has its own unique path for generating
// addresses it goes over all the cosigners and add their addresses
// for each key chain.
func (s *server) addressesToQuery(start, end uint32) (walletAddressSet, error) {
	addresses := make(walletAddressSet)
	for index := start; index < end; index++ {
		for cosignerIndex := uint32(0); cosignerIndex < uint32(len(s.keysFile.ExtendedPublicKeys)); cosignerIndex++ {
			for _, keychain := range keyChains {
				address := &walletAddress{
					index:         index,
					cosignerIndex: cosignerIndex,
					keyChain:      keychain,
				}
				addressString, err := s.walletAddressString(address)
				if err != nil {
					return nil, err
				}
				addresses[addressString] = address
			}
		}
	}

	return addresses, nil
}

// collectFarAddresses collects numIndexesToQuery addresses
// from the last point it stopped in the previous call.
func (s *server) collectFarAddresses() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	err := s.collectAddresses(s.nextSyncStartIndex, s.nextSyncStartIndex+numIndexesToQuery)
	if err != nil {
		return err
	}

	s.nextSyncStartIndex += numIndexesToQuery
	return nil
}

func (s *server) maxUsedIndex() uint32 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	maxUsedIndex := s.keysFile.LastUsedExternalIndex()
	if s.keysFile.LastUsedInternalIndex() > maxUsedIndex {
		maxUsedIndex = s.keysFile.LastUsedInternalIndex()
	}

	return maxUsedIndex
}

// collectRecentAddresses collects addresses from used addresses until
// the address with the index of the last used address + 1000.
// collectRecentAddresses scans addresses in batches of numIndexesToQuery,
// and releases the lock between scans.
func (s *server) collectRecentAddresses() error {
	maxUsedIndex := s.maxUsedIndex()
	for i := uint32(0); i < maxUsedIndex+1000; i += numIndexesToQuery {
		err := s.collectAddressesWithLock(i, i+numIndexesToQuery)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *server) collectAddressesWithLock(start, end uint32) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.collectAddresses(start, end)
}

func (s *server) collectAddresses(start, end uint32) error {
	addressSet, err := s.addressesToQuery(start, end)
	if err != nil {
		return err
	}

	getBalancesByAddressesResponse, err := s.rpcClient.GetBalancesByAddresses(addressSet.strings())
	if err != nil {
		return err
	}

	err = s.updateAddressesAndLastUsedIndexes(addressSet, getBalancesByAddressesResponse)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) updateAddressesAndLastUsedIndexes(requestedAddressSet walletAddressSet,
	getBalancesByAddressesResponse *appmessage.GetBalancesByAddressesResponseMessage) error {

	lastUsedExternalIndex := s.keysFile.LastUsedExternalIndex()
	lastUsedInternalIndex := s.keysFile.LastUsedInternalIndex()

	for _, entry := range getBalancesByAddressesResponse.Entries {
		walletAddress, ok := requestedAddressSet[entry.Address]
		if !ok {
			return errors.Errorf("Got result from address %s even though it wasn't requested", entry.Address)
		}

		if entry.Balance == 0 {
			continue
		}

		s.addressSet[entry.Address] = walletAddress

		if walletAddress.keyChain == libkaspawallet.ExternalKeychain {
			if walletAddress.index > lastUsedExternalIndex {
				lastUsedExternalIndex = walletAddress.index
			}
			continue
		}

		if walletAddress.index > lastUsedInternalIndex {
			lastUsedInternalIndex = walletAddress.index
		}
	}

	err := s.keysFile.SetLastUsedExternalIndex(lastUsedExternalIndex)
	if err != nil {
		return err
	}

	return s.keysFile.SetLastUsedInternalIndex(lastUsedInternalIndex)
}

func (s *server) getExistingUTXOsWithLock() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.getUTXOs()
}

// getUTXOSet clears the current UTXO set, and re-fills it with the given entries
func (s *server) getUTXOSet(entries []*appmessage.UTXOsByAddressesEntry) error {
	utxos := make(walletUTXOSet)

	for _, entry := range entries {
		outpoint, err := appmessage.RPCOutpointToDomainOutpoint(entry.Outpoint)
		if err != nil {
			return err
		}

		utxoEntry, err := appmessage.RPCUTXOEntryToUTXOEntry(entry.UTXOEntry)
		if err != nil {
			return err
		}

		address, ok := s.addressSet[entry.Address]
		if !ok {
			return errors.Errorf("Got result from address %s even though it wasn't requested", entry.Address)
		}

		utxos[*outpoint] = &walletUTXO{
			Outpoint:  outpoint,
			UTXOEntry: utxoEntry,
			address:   address,
		}
	}

	s.utxoSet = utxos

	return nil
}

func (s *server) getUTXOs() error {
	getUTXOsByAddressesResponse, err := s.rpcClient.GetUTXOsByAddresses(s.addressSet.strings())
	if err != nil {
		return err
	}

	return s.getUTXOSet(getUTXOsByAddressesResponse.Entries)
}

func (s *server) isSynced() bool {
	return s.nextSyncStartIndex > s.keysFile.LastUsedInternalIndex() && s.nextSyncStartIndex > s.keysFile.LastUsedExternalIndex()
}
