package blockstream

import "github.com/tiero/ocean/pkg/explorer"

type transaction struct {
	TxHash      string     `json:"txid"`
	TxVersion   int        `json:"version"`
	TxLocktime  int        `json:"locktime"`
	TxSize      int        `json:"size"`
	TxWeight    int        `json:"weight"`
	TxConfirmed bool       `json:"status.confirmed"`
	TxFees      int        `json:"fee"`
	TxInputs    []txInput  `json:"vin"`
	TxOutputs   []txOutput `json:"vout"`
}

type txInput struct {
	InputHash      string `json:"txid"`
	InputIndex     int    `json:"vout"`
	InputScriptSig string `json:"scriptsig"`
	InputSequence  int    `json:"sequence"`
	InputValue     int    `json:"prevout.value"`
	InputAddress   string `json:"prevout.scriptpubkey_address"`
}

type txOutput struct {
	OutputValue            int    `json:"value"`
	OutputAddress          string `json:"scriptpubkey_address"`
	OutputScriptPubKey     string `json:"scriptpubkey"`
	OutputScriptPubKeyType string `json:"scriptpubkey_type"`
}

func (t transaction) Hash() string {
	return t.TxHash
}

func (t transaction) Version() int {
	return t.TxVersion
}

func (t transaction) LockTime() int {
	return t.TxLocktime
}

func (t transaction) Size() int {
	return t.TxSize
}

func (t transaction) Weight() int {
	return t.TxWeight
}

func (t transaction) Confirmed() bool {
	return t.TxConfirmed
}

func (t transaction) Fees() int {
	return t.TxFees
}

func (t transaction) Inputs() []explorer.TxInput {
	inputs := make([]explorer.TxInput, len(t.TxInputs))
	for i, in := range t.TxInputs {
		inputs[i] = in
	}
	return inputs
}

func (t transaction) Outputs() []explorer.TxOutput {
	outputs := make([]explorer.TxOutput, len(t.TxOutputs))
	for i, out := range t.TxOutputs {
		outputs[i] = out
	}
	return outputs
}

func (i txInput) Hash() string {
	return i.InputHash
}

func (i txInput) Index() int {
	return i.InputIndex
}

func (i txInput) ScriptSig() string {
	return i.InputScriptSig
}

func (i txInput) Sequence() int {
	return i.InputSequence
}

func (i txInput) OutputValue() int {
	return i.InputValue
}

func (i txInput) Address() string {
	return i.InputAddress
}

func (o txOutput) Value() int {
	return o.OutputValue
}

func (o txOutput) Address() string {
	return o.OutputAddress
}

func (o txOutput) ScriptPubKey() string {
	return o.OutputScriptPubKey
}

func (o txOutput) ScriptPubKeyType() string {
	return o.OutputScriptPubKeyType
}
