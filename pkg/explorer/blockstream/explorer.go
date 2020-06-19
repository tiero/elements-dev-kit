package blockstream

import (
	"encoding/json"
	"fmt"
	"net/http"

	uhttp "github.com/tiero/ocean/internal/http"
	"github.com/tiero/ocean/pkg/explorer"
)

const service = "blockstream"

type blockstream struct {
	baseURL string
}

// NewExplorer returns a blockstream implementation of Explorer interface
// @param baseURL <string>: Blockstream API base URL
func NewExplorer(baseURL string) explorer.Explorer {
	bs := &blockstream{}
	bs.baseURL = baseURL
	return bs
}

// Ping is used to test that service API is up and running
func (bs *blockstream) Ping() int {
	url := fmt.Sprintf("%s/blocks/tip/height", bs.baseURL)
	status, _, _ := uhttp.NewHTTPRequest("GET", url, "", nil)
	return status
}

func (bs *blockstream) GetUnspents(address string) ([]explorer.Utxo, error) {
	url := fmt.Sprintf("%s/address/%s/utxo", bs.baseURL, address)
	status, resp, err := uhttp.NewHTTPRequest("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf(resp)
	}

	var out []utxo
	if err := json.Unmarshal([]byte(resp), &out); err != nil {
		return nil, err
	}

	unspents := make([]explorer.Utxo, len(out))
	for i, o := range out {
		unspents[i] = o
	}

	return unspents, nil
}

func (bs *blockstream) GetTransaction(hash string) (explorer.Transaction, error) {
	url := fmt.Sprintf("%s/tx/%s", bs.baseURL, hash)
	status, resp, err := uhttp.NewHTTPRequest("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf(resp)
	}

	out := &transaction{}
	if err := json.Unmarshal([]byte(resp), out); err != nil {
		return nil, err
	}

	return *out, nil
}

func (bs *blockstream) Broadcast(tx string) (string, error) {
	url := fmt.Sprintf("%s/tx", bs.baseURL)
	status, resp, err := uhttp.NewHTTPRequest("POST", url, tx, nil)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", fmt.Errorf(resp)
	}
	return resp, nil
}

func (bs *blockstream) EstimateFees() (explorer.Estimation, error) {
	url := fmt.Sprintf("%s/fee-estimates", bs.baseURL)
	status, resp, err := uhttp.NewHTTPRequest("GET", url, "", nil)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf(resp)
	}

	out := &estimation{}
	if err := json.Unmarshal([]byte(resp), out); err != nil {
		return nil, err
	}

	return *out, nil
}

type utxo struct {
	TxHash  string `json:"txid"`
	TxIndex int    `json:"vout"`
	TxValue int    `json:"value"`
	TxAsset string `json:"asset"`
}

func (u utxo) Hash() string {
	return u.TxHash
}

func (u utxo) Index() uint32 {
	return uint32(u.TxIndex)
}

func (u utxo) Value() uint64 {
	return uint64(u.TxValue)
}

func (u utxo) Asset() string {
	return u.TxAsset
}
