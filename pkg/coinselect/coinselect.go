package coinselect

import (
	"errors"

	"github.com/tiero/ocean/pkg/explorer"
)

// CoinSelect returns the utxos that satisfies the target amount.
func CoinSelect(utxos []explorer.Utxo, amount uint64, asset string) (unspents []explorer.Utxo, change uint64, err error) {
	change = 0
	unspents = []explorer.Utxo{}
	availableSats := uint64(0)

	for _, unspent := range utxos {
		if asset == unspent.Asset() {
			unspents = append(unspents, unspent)
			availableSats += unspent.Value()

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
