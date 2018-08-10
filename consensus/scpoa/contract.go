package scpoa

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
)

// ---------------------------------------------------------
// USED ONLY FOR TESTING
// ---------------------------------------------------------

// TODO
//   SignerDeploy
//   GetSigner from Signers address using tranaction and local key

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
	signerBytecode := common.Hex2Bytes(signerContract.Code)
	return string(signerABIBytes), signerBytecode
}
