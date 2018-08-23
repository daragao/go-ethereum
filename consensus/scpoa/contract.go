package scpoa

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"log"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// ---------------------------------------------------------
// USED ONLY FOR TESTING
// ---------------------------------------------------------

// TODO
//   SignerDeploy
//   GetSigner from Signers address using tranaction and local key

func newTx(
	ctx context.Context,
	backend bind.ContractBackend,
	from, to *common.Address,
	amount *big.Int,
	gasLimit uint64,
	payloadBytecode []byte,
) *types.Transaction {

	nonce, err := backend.PendingNonceAt(ctx, *from) // uint64(0)
	if err != nil {
		log.Fatal("Error getting pending nonce ", err)
	}
	gasPrice, err := backend.SuggestGasPrice(ctx) //new(big.Int)
	if err != nil {
		log.Fatal("Error suggesting gas price ", err)
	}

	// create contract transaction NewContractCreation is the same has NewTransaction with `to` == nil
	// tx := types.NewTransaction(nonce, nil, amount, gasLimit, gasPrice, payloadBytecode)
	var tx *types.Transaction
	if to == nil {
		tx = types.NewContractCreation(nonce, amount, gasLimit, gasPrice, payloadBytecode)
	} else {
		tx = types.NewTransaction(nonce, *to, amount, gasLimit, gasPrice, payloadBytecode)
	}
	return tx
}

// method created just to easily sign a tranasaction
func signTx(tx *types.Transaction, userKey *ecdsa.PrivateKey) *types.Transaction {
	signer := types.HomesteadSigner{} // this functions makes it easier to change signer if needed
	signedTx, err := types.SignTx(tx, signer, userKey)
	if err != nil {
		log.Fatal("Error signing tx: ", err)
	}
	return signedTx
}

func generateContractPayload(bytecode []byte, contractABIStr string, constructorArgs ...interface{}) []byte {
	abiContract, err := abi.JSON(strings.NewReader(contractABIStr))
	if err != nil {
		log.Fatal("ERROR reading contract ABI ", err)
	}
	packedABI, err := abiContract.Pack("", constructorArgs...)
	if err != nil {
		log.Fatal("ERROR packing ABI ", err)
	}
	payloadBytecode := append(bytecode, packedABI...)
	return payloadBytecode
}

func compileAndDeployContract(
	ctx context.Context,
	backend bind.ContractBackend,
	userKey *ecdsa.PrivateKey,
	binBytes []byte,
	abiStr string,
	amount *big.Int,
	gasLimit uint64,
	constructorArgs ...interface{},
) *types.Transaction {
	payload := generateContractPayload(binBytes, abiStr, constructorArgs...)
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	tx := newTx(ctx, backend, &userAddr, nil, amount, gasLimit, payload)
	signedTx := signTx(tx, userKey)

	err := backend.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("ERROR sending contract deployment transaction")
	}
	return signedTx
}

func callContract(
	ctx context.Context,
	client bind.ContractCaller,
	abiStr string,
	from, to common.Address,
	methodName string,
	out interface{},
	args ...interface{},
) {
	abiContract, err := abi.JSON(strings.NewReader(string(abiStr)))
	if err != nil {
		log.Fatal("ERROR reading contract ABI ", err)
	}

	input, err := abiContract.Pack(methodName, args...)
	if err != nil {
		log.Fatal("ERROR packing the method name for the contract call: ", err)
	}
	msg := ethereum.CallMsg{From: from, To: &to, Data: input}
	output, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		log.Fatal("ERROR calling the Ion Contract", err)
	}
	log.Printf("OUTPUT: %#v", output)
	err = abiContract.Unpack(out, methodName, output)
	if err != nil {
		log.Fatal("ERROR upacking the call: ", err)
	}
}

func compileSignerContract() (string, []byte) {
	basePath := "./" // os.Getenv("GOPATH") + "/src/github.com/clearmatics/ion/contracts/"
	scpoaContractPath := basePath + "SCPoa.sol"

	contracts, err := compiler.CompileSolidity("", scpoaContractPath)
	if err != nil {
		log.Fatal("ERROR failed to compile SCPoa.sol:", err)
	}

	signerContract := contracts[scpoaContractPath+":Signers"]
	signerABIBytes, err := json.Marshal(signerContract.Info.AbiDefinition)
	if err != nil {
		log.Fatal("ERROR marshaling Signer contract ABI:", err)
	}
	signerBytecode := common.Hex2Bytes(signerContract.Code[2:])
	return string(signerABIBytes), signerBytecode
}
