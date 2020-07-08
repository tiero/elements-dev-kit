package coinselect

import (
	"encoding/hex"
	"errors"

	"github.com/tiero/ocean/internal/bufferutil"
	"github.com/tiero/ocean/pkg/explorer"
	"github.com/vulpemventures/go-elements/confidential"
)

// Coins defines the struct thta holds utxos and relative blinding keys.
type Coins struct {
	Utxos        []explorer.Utxo
	BlindingKey  []byte
	BlindingKeys [][]byte
}

// CoinSelect returns the utxos that satisfies the target amount and target asset.
// TODO implement branch and bound algorithm.
func (cs *Coins) CoinSelect(amount uint64, asset string) (unspents []explorer.Utxo, change uint64, err error) {
	change = 0
	unspents = []explorer.Utxo{}
	availableSats := uint64(0)

	for index, unspent := range cs.Utxos {
		u := unspent
		assetHash := u.Asset()
		amountSatoshis := u.Value()
		if len(u.AssetCommitment()) > 0 && len(u.ValueCommitment()) > 0 {
			bk := cs.BlindingKey
			if len(cs.BlindingKeys) > 0 {
				bk = cs.BlindingKeys[index]
			}
			av, err := unblindUxto(u, bk)
			if err != nil {
				return nil, 0, err
			}
			assetHash = av.asset
			amountSatoshis = av.value
		}
		if asset == assetHash {
			unspents = append(unspents, unspent)
			availableSats += amountSatoshis

			if availableSats >= amount {
				break
			}
		}
	}

	if availableSats < amount {
		return nil, 0, errors.New("You do not have enough coins")
	}

	change = availableSats - amount

	return unspents, change, nil
}

type assetAndValue struct {
	asset string
	value uint64
}

func unblindUxto(prevout explorer.Utxo, blindingKey []byte) (*assetAndValue, error) {
	assetCommitment, err := hex.DecodeString(prevout.AssetCommitment())
	if err != nil {
		return nil, err
	}
	valueCommitment, err := hex.DecodeString(prevout.ValueCommitment())
	if err != nil {
		return nil, err
	}
	nonce, err := confidential.NonceHash(
		prevout.Nonce(),
		blindingKey,
	)
	if err != nil {
		return nil, err
	}
	unblindOutputArg := confidential.UnblindOutputArg{
		Nonce:           nonce,
		Rangeproof:      prevout.RangeProof(),
		ValueCommitment: valueCommitment,
		AssetCommitment: assetCommitment,
		ScriptPubkey:    prevout.Script(),
	}

	output, err := confidential.UnblindOutput(unblindOutputArg)
	if err != nil {
		return nil, err
	}
	assetHash := hex.EncodeToString(bufferutil.ReverseBytes(output.Asset[:]))
	return &assetAndValue{assetHash, output.Value}, nil
}
