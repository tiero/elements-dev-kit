package partial

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/tiero/ocean/internal/bufferutil"
	"github.com/tiero/ocean/pkg/keypair"
	"github.com/vulpemventures/go-elements/confidential"
	"github.com/vulpemventures/go-elements/network"
	"github.com/vulpemventures/go-elements/payment"
	"github.com/vulpemventures/go-elements/pset"
	"github.com/vulpemventures/go-elements/transaction"
)

// Partial defines a Partial Signed Elements Transaction
type Partial struct {
	Data    *pset.Pset
	Network *network.Network
}

// WitnessUtxo defines a witness utxo
type WitnessUtxo struct {
	Asset  string
	Value  uint64
	Script []byte
}

// ConfidentialWitnessUtxo defines a confidential witness utxo
type ConfidentialWitnessUtxo struct {
	AssetCommitment string
	ValueCommitment string
	Script          []byte
	Nonce           []byte
	RangeProof      []byte
	SurjectionProof []byte
}

//NewPartial returns a Partial instance with an empty pset in Partial.Data and the selected Network
func NewPartial(net *network.Network) *Partial {
	currentNetwork := &network.Liquid
	if net != nil {
		currentNetwork = net
	}
	emptyPset, _ := pset.New([]*transaction.TxInput{}, []*transaction.TxOutput{}, 2, 0)
	return &Partial{Data: emptyPset, Network: currentNetwork}
}

//AddInput adds an utxo to a Partial Signed Elements Transaction
func (p *Partial) AddInput(hash string, index uint32, witnessUtxo *WitnessUtxo, nonWitnessUtxo []byte) error {
	updater, err := pset.NewUpdater(p.Data)
	if err != nil {
		return err
	}

	inputHash, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	inputHash = bufferutil.ReverseBytes(inputHash)
	inputIndex := index
	input := transaction.NewTxInput(inputHash, inputIndex)

	updater.AddInput(input)
	lastAdded := len(updater.Upsbt.Inputs) - 1

	err = updater.AddInSighashType(txscript.SigHashAll, lastAdded)
	if err != nil {
		return err
	}

	if witnessUtxo != nil {
		elementsValue, err := confidential.SatoshiToElementsValue(witnessUtxo.Value)
		if err != nil {
			return err
		}
		elementsAsset, err := AssetHashToBytes(witnessUtxo.Asset, false)
		if err != nil {
			return err
		}
		witnessUtxo := &transaction.TxOutput{Asset: elementsAsset, Value: elementsValue[:], Script: witnessUtxo.Script}
		updater.AddInWitnessUtxo(witnessUtxo, lastAdded)
		p.Data = updater.Upsbt
		return nil
	}

	if nonWitnessUtxo != nil {
		return errors.New("Not yet implemented. Only segwit inputs supported")
	}

	return errors.New("Either witnessUtxo or nonWitnessUtxo is missing")
}

//AddBlindedInput adds an utxo to a Partial Signed Elements Transaction
func (p *Partial) AddBlindedInput(hash string, index uint32, witnessUtxo *ConfidentialWitnessUtxo, nonWitnessUtxo []byte) error {
	updater, err := pset.NewUpdater(p.Data)
	if err != nil {
		return err
	}

	inputHash, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}
	inputHash = bufferutil.ReverseBytes(inputHash)
	inputIndex := index
	input := transaction.NewTxInput(inputHash, inputIndex)

	updater.AddInput(input)
	lastAdded := len(updater.Upsbt.Inputs) - 1

	err = updater.AddInSighashType(txscript.SigHashAll, lastAdded)
	if err != nil {
		return err
	}

	if witnessUtxo != nil {
		valueCommitment, err := hex.DecodeString(witnessUtxo.ValueCommitment)
		if err != nil {
			return err
		}
		assetCommitment, err := hex.DecodeString(witnessUtxo.AssetCommitment)
		if err != nil {
			return err
		}
		witnessUtxo := &transaction.TxOutput{
			Asset:           assetCommitment,
			Value:           valueCommitment,
			Script:          witnessUtxo.Script,
			Nonce:           witnessUtxo.Nonce,
			RangeProof:      witnessUtxo.RangeProof,
			SurjectionProof: witnessUtxo.SurjectionProof,
		}
		updater.AddInWitnessUtxo(witnessUtxo, lastAdded)
		p.Data = updater.Upsbt
		return nil
	}

	if nonWitnessUtxo != nil {
		return errors.New("Not yet implemented. Only segwit inputs supported")
	}

	return errors.New("Either witnessUtxo or nonWitnessUtxo is missing")

}

// AddOutput adds an output to a Partial Signed Elements Transaction
func (p *Partial) AddOutput(asset string, value uint64, script []byte, blinded bool) error {
	updater, err := pset.NewUpdater(p.Data)
	if err != nil {
		return err
	}

	elementsValue, err := confidential.SatoshiToElementsValue(value)
	if err != nil {
		return err
	}
	elementsAsset, err := AssetHashToBytes(asset, blinded)
	if err != nil {
		return err
	}
	output := transaction.NewTxOutput(elementsAsset, elementsValue[:], script)

	updater.AddOutput(output)
	p.Data = updater.Upsbt
	return nil
}

//BlindWithKeys unblinds all the inputs and blinds all the outputs with the provided arrays of keys
func (p *Partial) BlindWithKeys(blindingPrivateKeys [][]byte, blindingPublicKeys [][]byte) error {
	blinder, err := pset.NewBlinder(
		p.Data,
		blindingPrivateKeys,
		blindingPublicKeys,
		nil,
		nil)
	if err != nil {
		return err
	}
	err = blinder.Blind()
	if err != nil {
		return err
	}

	return nil
}

//SignWithPrivateKey signs a witness input with a provided EC private key
func (p *Partial) SignWithPrivateKey(index int, keyPair *keypair.KeyPair) error {
	updater, err := pset.NewUpdater(p.Data)
	if err != nil {
		return err
	}

	if index > (len(updater.Upsbt.Inputs) - 1) {
		return errors.New("index out of range")
	}

	currInput := updater.Upsbt.Inputs[index]
	if currInput.NonWitnessUtxo != nil {
		return errors.New("Only segwit input supported")
	}

	var witHash [32]byte
	script := currInput.WitnessUtxo.Script
	if script[0] == txscript.OP_0 {
		println("native")
		prevoutPayment, err := payment.FromScript(script, p.Network, nil)
		if err != nil {
			return err
		}
		// legacy Script
		legacyScript := append(append([]byte{0x76, 0xa9, 0x14}, prevoutPayment.Hash...), []byte{0x88, 0xac}...)
		witHash = updater.Upsbt.UnsignedTx.HashForWitnessV0(index, legacyScript, currInput.WitnessUtxo.Value[:], txscript.SigHashAll)
	}

	if script[0] == txscript.OP_HASH160 {
		println("wrapped")
		println(hex.EncodeToString(script))
		witHash, err = updater.Upsbt.UnsignedTx.HashForSignature(index, script, txscript.SigHashAll)
		if err != nil {
			return err
		}
	}

	sig, err := keyPair.PrivateKey.Sign(witHash[:])
	if err != nil {
		return fmt.Errorf("PrivateKey Sign: %w", err)
	}

	sigWithHashType := append(sig.Serialize(), byte(txscript.SigHashAll))
	println(index)
	_, err = updater.Sign(index, sigWithHashType, keyPair.PublicKey.SerializeCompressed(), nil, nil)
	if err != nil {
		return fmt.Errorf("Updater Sign: %w", err)
	}

	p.Data = updater.Upsbt
	return nil
}

//AssetHashToBytes reverse decode from hex string and reverse it adding a 0x01 byte for ublinded asset
func AssetHashToBytes(hash string, blinded bool) ([]byte, error) {
	firstByte := byte(0x01)
	if blinded {
		firstByte = byte(0x00)
	}
	assetBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	assetBytes = append([]byte{firstByte}, bufferutil.ReverseBytes(assetBytes)...)
	return assetBytes, nil
}
