package partial

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/tiero/ocean/pkg/coinselect"
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
	wif, err := btcutil.DecodeWIF("cPCJXXKTkHrur9UcYxdh51MQ8NYArfnjdV4xSG9eqEbhzjRQh6h8")
	if err != nil {
		t.Fatal(err)
	}
	kp, err := keypair.FromPrivateKey(hex.EncodeToString(wif.PrivKey.Serialize()))
	if err != nil {
		t.Fatal(err)
	}
	// Alice Blinding KeyPair
	kpBlind, err := keypair.FromPrivateKey("fd9123214784758c69351f45aebf3c719533a05c5fa017a466b4f31328487552")
	if err != nil {
		t.Fatal(err)
	}

	//Bob KeyPair
	bobKeyPair, err := keypair.FromPrivateKey(bobHex)
	if err != nil {
		t.Fatal(err)
	}
	// Bob Blinding KeyPair
	bobBlind, err := keypair.FromPrivateKey(bobBlindHex)
	if err != nil {
		t.Fatal(err)
	}

	alice := payment.FromPublicKey(kp.PublicKey, &network.Regtest, kpBlind.PublicKey)
	wrappedAlice, err := payment.FromPayment(alice)
	if err != nil {
		t.Fatal(err)
	}
	aliceConfAddress := wrappedAlice.ConfidentialScriptHash()
	//aliceConfAddr := "AzpwTgRMptQ8CB1UTrc6ereqFt6ZDTwJSgm6iu2BHRZbrXEXyu8x2cjAkZR5BeVznVeiTCCqqsQKzcwD"
	println(aliceConfAddress)
	nativeSegwitAliceConfAddr, err := alice.ConfidentialWitnessPubKeyHash()
	if err != nil {
		t.Fatal(err)
	}
	println(nativeSegwitAliceConfAddr)

	bob := payment.FromPublicKey(bobKeyPair.PublicKey, &network.Regtest, nil)

	// Fund sender address.
	_, err = faucet(aliceConfAddress)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	// Retrieve sender utxos.
	utxos, err := e.GetUnspents(aliceConfAddress)
	if err != nil {
		t.Fatal(err)
	}

	coins := &coinselect.Coins{Utxos: utxos, BlindingKey: kpBlind.PrivateKey.Serialize()}
	selectedUtxos, change, err := coins.CoinSelect(50000000, network.Regtest.AssetID)
	if err != nil {
		t.Fatal(err)
	}

	for _, utxo := range selectedUtxos {
		err := p.AddBlindedInput(utxo.Hash(), utxo.Index(), &ConfidentialWitnessUtxo{
			AssetCommitment: utxo.AssetCommitment(),
			ValueCommitment: utxo.ValueCommitment(),
			Script:          utxo.Script(),
			Nonce:           utxo.Nonce(),
			RangeProof:      utxo.RangeProof(),
			SurjectionProof: utxo.SurjectionProof(),
		}, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	var fee uint64 = 500
	p.AddOutput(network.Regtest.AssetID, 50000000, bob.Script, false)
	p.AddOutput(network.Regtest.AssetID, change-fee, alice.Script, false)
	p.AddOutput(network.Regtest.AssetID, fee, []byte{}, false)

	blindingPrivKeysOfInputs := [][]byte{kpBlind.PrivateKey.Serialize()}
	blindingPubKeysOfOutputs := [][]byte{bobBlind.PublicKey.SerializeCompressed(), bobBlind.PublicKey.SerializeCompressed()}
	err = p.BlindWithKeys(blindingPrivKeysOfInputs, blindingPubKeysOfOutputs)
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

	/* 	b64, err := pFinalized.ToBase64()
	   	if err != nil {
	   		t.Errorf("base64: %w", err)
	   	}

	   	fmt.Println(b64) */

	// Extract the final signed transaction from the Pset wrapper.
	finalTx, err := pset.Extract(pFinalized)
	if err != nil {
		t.Fatal(err)
	}

	// Serialize the transaction and try to broadcast.
	txHex, err := finalTx.ToHex()
	if err != nil {
		t.Fatal(err)
	}
	txHash, err := e.Broadcast(txHex)
	if err != nil {
		t.Fatal(err)
	}

	println(txHash)

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
