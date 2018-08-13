package scpoa

import (
	"encoding/json"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
)

// ---------------------------------------------------------
// USED ONLY FOR TESTING
// ---------------------------------------------------------

// TODO
//   SignerDeploy
//   GetSigner from Signers address using tranaction and local key

func deploySignerContract(abiStr string, bytecode []byte, signers []common.Address) {
	// generate payload
	abiContract, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		log.Fatal("ERROR reading contract ABI ", err)
	}
	packedABI, err := abiContract.Pack("", signers)
	if err != nil {
		log.Fatal("ERROR packing ABI ", err)
	}
	payloadBytecode := append(bytecode, packedABI...)

	gasPrice := big.NewInt(0)
	gasLimit := uint64(30000)
	nonce := uint64(0)
	tx := types.NewContractCreation(nonce, nil, gasLimit, gasPrice, payloadBytecode)
	log.Println(tx)
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
