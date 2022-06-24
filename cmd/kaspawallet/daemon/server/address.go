package server

import (
	"context"
	"fmt"

	"github.com/kaspanet/kaspad/cmd/kaspawallet/daemon/pb"
	"github.com/kaspanet/kaspad/cmd/kaspawallet/libkaspawallet"
	"github.com/kaspanet/kaspad/util"
	"github.com/pkg/errors"
)

func (s *server) changeAddress() (util.Address, *walletAddress, error) {
	err := s.keysFile.SetLastUsedInternalIndex(s.keysFile.LastUsedInternalIndex() + 1)
	if err != nil {
		return nil, nil, err
	}

	err = s.keysFile.Save()
	if err != nil {
		return nil, nil, err
	}

	walletAddr := &walletAddress{
		index:         s.keysFile.LastUsedInternalIndex(),
		cosignerIndex: s.keysFile.CosignerIndex,
		keyChain:      libkaspawallet.InternalKeychain,
	}
	path := s.walletAddressPath(walletAddr)
	address, err := libkaspawallet.Address(s.params, s.keysFile.ExtendedPublicKeys, s.keysFile.MinimumSignatures, path, s.keysFile.ECDSA)
	if err != nil {
		return nil, nil, err
	}
	return address, walletAddr, nil
}

func (s *server) ShowAddresses(_ context.Context, request *pb.ShowAddressesRequest) (*pb.ShowAddressesResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.isSynced() {
		return nil, errors.Errorf("wallet daemon is not synced yet, %s", s.formatSyncStateReport())
	}

	addresses := make([]string, s.keysFile.LastUsedExternalIndex())
	for i := uint32(1); i <= s.keysFile.LastUsedExternalIndex(); i++ {
		walletAddr := &walletAddress{
			index:         i,
			cosignerIndex: s.keysFile.CosignerIndex,
			keyChain:      libkaspawallet.ExternalKeychain,
		}
		path := s.walletAddressPath(walletAddr)
		address, err := libkaspawallet.Address(s.params, s.keysFile.ExtendedPublicKeys, s.keysFile.MinimumSignatures, path, s.keysFile.ECDSA)
		if err != nil {
			return nil, err
		}
		addresses[i-1] = address.String()
	}

	return &pb.ShowAddressesResponse{Address: addresses}, nil
}

func (s *server) NewAddress(_ context.Context, request *pb.NewAddressRequest) (*pb.NewAddressResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.isSynced() {
		return nil, errors.Errorf("wallet daemon is not synced yet, %s", s.formatSyncStateReport())
	}

	err := s.keysFile.SetLastUsedExternalIndex(s.keysFile.LastUsedExternalIndex() + 1)
	if err != nil {
		return nil, err
	}

	err = s.keysFile.Save()
	if err != nil {
		return nil, err
	}

	walletAddr := &walletAddress{
		index:         s.keysFile.LastUsedExternalIndex(),
		cosignerIndex: s.keysFile.CosignerIndex,
		keyChain:      libkaspawallet.ExternalKeychain,
	}
	path := s.walletAddressPath(walletAddr)
	address, err := libkaspawallet.Address(s.params, s.keysFile.ExtendedPublicKeys, s.keysFile.MinimumSignatures, path, s.keysFile.ECDSA)
	if err != nil {
		return nil, err
	}

	return &pb.NewAddressResponse{Address: address.String()}, nil
}

func (s *server) CheckIfAddressesAreValid(_ context.Context, request *pb.CheckIfAddressesAreValidRequest) (*pb.CheckIfAddressesAreValidResponse, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	areValidAddresses := make([]*pb.IsValidAddress, len(request.Addresses))
	for i, address := range request.Addresses {
		_, err := util.DecodeAddress(address, s.params.Prefix)
		if err != nil {
			areValidAddresses[i] = &pb.IsValidAddress{
				Address: address,
				IsValid: false,
				Error:   err.Error(),
			}
		} else {
			areValidAddresses[i] = &pb.IsValidAddress{
				Address: address,
				IsValid: true,
			}
		}
	}

	return &pb.CheckIfAddressesAreValidResponse{
		AreValidAddresses: areValidAddresses,
	}, nil
}

func (s *server) walletAddressString(wAddr *walletAddress) (string, error) {
	path := s.walletAddressPath(wAddr)
	addr, err := libkaspawallet.Address(s.params, s.keysFile.ExtendedPublicKeys, s.keysFile.MinimumSignatures, path, s.keysFile.ECDSA)
	if err != nil {
		return "", err
	}

	return addr.String(), nil
}

func (s *server) walletAddressPath(wAddr *walletAddress) string {
	if s.isMultisig() {
		return fmt.Sprintf("m/%d/%d/%d", wAddr.cosignerIndex, wAddr.keyChain, wAddr.index)
	}
	return fmt.Sprintf("m/%d/%d", wAddr.keyChain, wAddr.index)
}

func (s *server) isMultisig() bool {
	return len(s.keysFile.ExtendedPublicKeys) > 1
}
