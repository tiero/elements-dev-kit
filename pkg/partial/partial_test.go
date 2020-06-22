package partial

import (
	"reflect"
	"testing"

	"github.com/vulpemventures/go-elements/network"
	"github.com/vulpemventures/go-elements/pset"
	"github.com/vulpemventures/go-elements/transaction"
)

const EMPTYPSET = "cHNldP8BAAsCAAAAAAAAAAAAAAA="

func TestNewPartial(t *testing.T) {
	emptyPset, _ := pset.New([]*transaction.TxInput{}, []*transaction.TxOutput{}, 2, 0)
	want := &Partial{Data: emptyPset, Network: &network.Liquid}

	if got := NewPartial(nil); !reflect.DeepEqual(got, want) {
		t.Errorf("NewPartial() = %v, want %v", got, want)
	}
	if gotB64, _ := want.Data.ToBase64(); gotB64 != EMPTYPSET {
		t.Errorf("NewPartial() = %v, want %v", gotB64, EMPTYPSET)
	}

}
