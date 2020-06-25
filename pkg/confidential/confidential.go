package confidential

import (
	"errors"

	"github.com/btcsuite/btcutil/base58"
	addressPackage "github.com/vulpemventures/go-elements/address"
	"github.com/vulpemventures/go-elements/network"
)

const (
	P2Pkh = iota
	P2Sh
	ConfidentialP2Pkh
	ConfidentialP2Sh
	P2Wpkh
	P2Wsh
	ConfidentialP2Wpkh
	ConfidentialP2Wsh
)

//ToBlindingKey returns  the blinding key encoded in the address
func ToBlindingKey(address string, net network.Network) ([]byte, error) {
	addressType, err := addressPackage.DecodeType(address, net)
	if err != nil {
		return nil, err
	}

	switch addressType {
	case ConfidentialP2Pkh:
		decoded, _, err := base58.CheckDecode(address)
		if err != nil {
			return nil, err
		}
		prefixBytes := 1
		prefixPlusBlindKeySize := prefixBytes + 33
		blindingKey := decoded[prefixBytes:prefixPlusBlindKeySize]
		return blindingKey, nil
	case ConfidentialP2Sh:
		decoded, _, err := base58.CheckDecode(address)
		if err != nil {
			return nil, err
		}
		prefixBytes := 1
		prefixPlusBlindKeySize := prefixBytes + 33
		blindingKey := decoded[prefixBytes:prefixPlusBlindKeySize]
		return blindingKey, nil
	case ConfidentialP2Wpkh, ConfidentialP2Wsh:
		fromBlech32, err := addressPackage.FromBlech32(address)
		if err != nil {
			return nil, err
		}
		return fromBlech32.PublicKey, nil
	default:
		return nil, errors.New("unsupported address type")
	}
}
