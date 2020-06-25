package explorer

// Explorer interface defines method for a block explorer
type Explorer interface {
	Ping() int
	GetUnspents(address string) ([]Utxo, error)
	GetTransaction(hash string) (Transaction, error)
	GetTransactionHex(hash string) (string, error)
	Broadcast(tx string) (string, error)
	EstimateFees() (Estimation, error)
}

// Utxo defines the unspent from the explorer
type Utxo interface {
	Hash() string
	Index() uint32
	Value() uint64
	Asset() string
	ValueCommitment() string
	AssetCommitment() string
	Nonce() []byte
	Script() []byte
	RangeProof() []byte
	SurjectionProof() []byte
}

// Transaction interface defines what data a transaction must include
type Transaction interface {
	Hash() string
	Version() int
	LockTime() int
	Size() int
	Weight() int
	Confirmed() bool
	Fees() int
	Inputs() []TxInput
	Outputs() []TxOutput
}

// TxInput defines what data a tx input must include
type TxInput interface {
	Hash() string
	Index() int
	ScriptSig() string
	Sequence() int
	OutputValue() int
	Address() string
}

// TxOutput interface defines what data a tx output must include
type TxOutput interface {
	Value() int
	Address() string
	ScriptPubKey() string
	ScriptPubKeyType() string
}

// Estimation interface defines what data an estimation response must include
type Estimation interface {
	Low() float64
	Medium() float64
	High() float64
}
