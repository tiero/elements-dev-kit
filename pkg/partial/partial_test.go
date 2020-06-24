package partial

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/tiero/ocean/pkg/explorer/blockstream"
	"github.com/tiero/ocean/pkg/keypair"
	"github.com/vulpemventures/go-elements/network"
	"github.com/vulpemventures/go-elements/payment"
	"github.com/vulpemventures/go-elements/pset"
	"github.com/vulpemventures/go-elements/transaction"
)

const EMPTYPSET = "cHNldP8BAAsCAAAAAAAAAAAAAAA="
const aliceHex = "bfb96a215dfb07d1a193464174b9ea8e91f2a15bba79800dea838add330f6d86"
const aliceBlindHex = "dd65e215154c13b1c14f9dc0aa7cfc1f40414f214bd0c5dfe2d370880bdf8356"
const bobHex = "1804e76aa3016013bc9969103554668913cf697c03c23aecb28136d0e0ac16f0"
const bobBlindHex = "dd65e215154c13b1c14f9dc0aa7cfc1f40414f214bd0c5dfe2d370880bdf8356"
const defaultExplorer = "http://localhost:3001"

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

func TestCreatePsetWithBlindedInput(t *testing.T) {
	explorerURL, ok := os.LookupEnv("API_URL")
	if !ok {
		explorerURL = defaultExplorer
	}
	//Explorer
	e := blockstream.NewExplorer(explorerURL)
	// PSET
	p := NewPartial(&network.Regtest)
	// Alice keypair
	kp, err := keypair.FromPrivateKey(aliceHex)
	if err != nil {
		t.Fatal(err)
	}
	// Alice Blinding KeyPair
	kpBlind, err := keypair.FromPrivateKey(aliceBlindHex)
	if err != nil {
		t.Fatal(err)
	}

	//Bob KeyPair
	bobKeyPair, err := keypair.FromPrivateKey(bobHex)
	if err != nil {
		t.Fatal(err)
	}

	bobBlind, err := keypair.FromPrivateKey(bobBlindHex)
	if err != nil {
		t.Fatal(err)
	}

	alice := payment.FromPublicKey(kp.PublicKey, &network.Regtest, kpBlind.PublicKey)
	aliceConfAddr, err := alice.ConfidentialWitnessPubKeyHash()
	if err != nil {
		t.Fatal(err)
	}

	bob := payment.FromPublicKey(bobKeyPair.PublicKey, &network.Regtest, nil)

	// Fund sender address.
	_, err = faucet(aliceConfAddr)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve sender utxos.
	utxos, err := e.GetUnspents(aliceConfAddr)
	if err != nil {
		t.Fatal(err)
	}
	utxo := utxos[0]
	prevoutTxHex, err := e.GetTransactionHex(utxo.Hash())
	if err != nil {
		t.Fatal(err)
	}

	trx, err := transaction.NewTxFromHex(prevoutTxHex)
	if err != nil {
		t.Fatal(err)
	}

	err = p.AddBlindedInput(utxo.Hash(), utxo.Index(), &ConfidentialWitnessUtxo{
		AssetCommitment: utxo.AssetCommitment(),
		ValueCommitment: utxo.ValueCommitment(),
		Script:          alice.Script,
		Nonce:           trx.Outputs[utxo.Index()].Nonce,
		RangeProof:      trx.Outputs[utxo.Index()].RangeProof,
		SurjectionProof: trx.Outputs[utxo.Index()].SurjectionProof,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = p.AddOutput(network.Regtest.AssetID, 99995000, bob.Script, false)
	if err != nil {
		t.Fatal(err)
	}
	err = p.AddOutput(network.Regtest.AssetID, 5000, []byte{}, false)
	if err != nil {
		t.Fatal(err)
	}

	println(len(p.Data.Outputs))

	blindingPrivKeys := [][]byte{kpBlind.PrivateKey.Serialize()}
	blindingPubKeys := [][]byte{bobBlind.PublicKey.SerializeCompressed()}
	blinder, err := pset.NewBlinder(
		p.Data,
		blindingPrivKeys,
		blindingPubKeys,
		nil,
		nil)
	if err != nil {
		t.Fatal(err)
	}
	err = blinder.Blind()
	if err != nil {
		t.Fatal(err)
	}

	err = p.SignWithPrivateKey(0, kp)
	if err != nil {
		t.Fatal(err)
	}

	pFinalized := p.Data
	err = pset.FinalizeAll(pFinalized)
	if err != nil {
		t.Errorf("sign: %w", err)
	}

	if !pFinalized.IsComplete() {
		t.Errorf("pset not complete: %w", err)
	}

	err = pFinalized.SanityCheck()
	if err != nil {
		t.Errorf("sanity check: %w", err)
	}

	b64, err := pFinalized.ToBase64()
	if err != nil {
		t.Errorf("base64: %w", err)
	}

	fmt.Println(b64)

}

func faucet(address string) (string, error) {
	baseURL, ok := os.LookupEnv("API_URL")
	if !ok {
		baseURL = defaultExplorer
	}
	url := baseURL + "/faucet"
	payload := map[string]string{"address": address}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "appliation/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBody := map[string]string{}
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return "", err
	}

	return respBody["txId"], nil
}
