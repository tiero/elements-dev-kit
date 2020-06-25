package confidential

import (
	"encoding/hex"
	"testing"

	"github.com/vulpemventures/go-elements/network"
)

func TestToBlindingKey(t *testing.T) {

	got, err := ToBlindingKey("AzpwTgRMptQ8CB1UTrc6ereqFt6ZDTwJSgm6iu2BHRZbrXEXyu8x2cjAkZR5BeVznVeiTCCqqsQKzcwD", network.Regtest)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if hex.EncodeToString(got) != "03cd357bea612515608889894ab600625bf64667d53f84ae3734a79ea011488f42" {
		t.Fatal("Not the right blinding key")
	}
}
