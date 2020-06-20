# 🌊 ocean
Go package for building Elements/Liquid wallets.


## Usage 

### Install

```sh
$ go get github.com/tiero/ocean
```

### Quickstart

```go
import "github.com/tiero/ocean/partial"
// Create an empty PSET
pset := NewPartial()

// Add a segwit input and a unblinded output
pset.AddInput(hash, index , witnessUtxo, nil)
pset.AddOutput(asset, value, script, false) 
// Generate a KeyPair object
privKeyHex := "bfb96a215dfb07d1a193464174b9ea8e91f2a15bba79800dea838add330f6d86"
keyPair, _ := keypair.FromPrivateKey(privKeyHex)
//Sign the input with a private key
pset.SignWithPrivateKey(0, keyPair)
```

## Development

### Clone

```sh
$ git clone https://github.com/tiero/ocean
```

### Run test 

```sh
$ go test -race ./...
```
