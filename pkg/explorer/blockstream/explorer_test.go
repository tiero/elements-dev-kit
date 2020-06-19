package blockstream

import (
	"testing"
)

const (
	baseURL = "https://blockstream.info/liquid/api"
	address = "ex1qs95hf3q25uh9u3gjacwuzwt49h5ef7e5ujw00v"
	hash    = "e32b095696c00ae94b95a2f74cc6ddf23f9791381f332a64423e9187339fcb8b"
	badHash = "02b082113e35d5386285094c2829e7e2963fa0b5369fb7f4b79c4c90877dcd3d"
)

var blockexplorer = NewExplorer(baseURL)

func TestPingBlockstream(t *testing.T) {
	expectedStatus := 200
	status := blockexplorer.Ping()
	if status != expectedStatus {
		t.Fatalf("Got: %d, expected: %d\n", status, expectedStatus)
	}
}

func TestGetUtxoFromBlockstream(t *testing.T) {
	utxos, err := blockexplorer.GetUnspents(address)
	if err != nil {
		t.Fatal(err)
	}
	if len(utxos) <= 0 {
		t.Fatal("Got empty utxo list, expected not empty")
	}
	for _, utxo := range utxos[:3] {
		t.Log("utxo_hash:", utxo.Hash())
		t.Log("utxo_index:", utxo.Index())
		t.Log("utxo_value:", utxo.Value())
	}
}

func TestGetTxBlockstream(t *testing.T) {
	tx, err := blockexplorer.GetTransaction(hash)
	if err != nil {
		t.Fatal(err)
	}
	if tx == nil {
		t.Fatal("Got empty transaction, expected not empty")
	}
	t.Log("tx_hash:", tx.Hash())
	t.Log("tx_version:", tx.Version())
	t.Log("tx_size:", tx.Size())
	t.Log("tx_weight:", tx.Weight())
	t.Log("tx_locktime:", tx.LockTime())
	t.Log("tx_fees:", tx.Fees())
	t.Log("tx_confirmed:", tx.Confirmed())
	t.Log("tx_input_size:", len(tx.Inputs()))
	t.Log("tx_output_size:", len(tx.Outputs()))
}

func TestGetTxShouldFail(t *testing.T) {
	expectedError := "Transaction not found"
	_, err := blockexplorer.GetTransaction(badHash)
	if err == nil {
		t.Fatal("Should have failed before")
	}
	if err.Error() != expectedError {
		t.Fatalf("Got error: %s, expected: %s", err, expectedError)
	}
}

func TestEstimateFees(t *testing.T) {
	estimation, err := blockexplorer.EstimateFees()
	if err != nil {
		t.Fatal(err)
	}
	if estimation == nil {
		t.Fatal("Got empty estimation, expected not empyt")
	}

	t.Log("high_fee_per_byte:", estimation.High())
	t.Log("medium_fee_per_byte:", estimation.Medium())
	t.Log("low_fee_per_byte:", estimation.Low())
}
