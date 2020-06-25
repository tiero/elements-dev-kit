package coinselect

import (
	"testing"

	"github.com/tiero/ocean/pkg/explorer"
)

type utxo struct {
	hash            string
	index           int
	value           int
	asset           string
	valuecommitment string
	assetcommitment string
}

func (u utxo) Hash() string {
	return u.hash
}

func (u utxo) Index() uint32 {
	return uint32(u.index)
}

func (u utxo) Value() uint64 {
	return uint64(u.value)
}

func (u utxo) Asset() string {
	return u.asset
}

func (u utxo) ValueCommitment() string {
	return u.valuecommitment
}

func (u utxo) AssetCommitment() string {
	return u.assetcommitment
}

func TestCoinSelect(t *testing.T) {

	testUtxo1 := &utxo{"foo", 0, 1000, "dollar", "", ""}
	testUtxo2 := &utxo{"bar", 0, 500, "dollar", "", ""}

	gotUnspents, gotChange, err := CoinSelect([]explorer.Utxo{testUtxo1, testUtxo2}, 800, "dollar")
	if err != nil {
		t.Errorf("Should not have throwed any error")
	}
	if len(gotUnspents) != 1 {
		t.Errorf("CoinSelect() gotUnspents")
	}
	if gotChange != 200 {
		t.Errorf("CoinSelect() gotChange")
	}
}
