package scpoa

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func Test_CompileSigner(t *testing.T) {
	ctx := context.Background()

	userKey, _ := crypto.GenerateKey()
	userAddr := crypto.PubkeyToAddress(userKey.PublicKey)
	userIntialBalance := big.NewInt(1000000000)
	// ---------------------------------------------
	// START BLOCKCHAIN SIMULATOR
	// ---------------------------------------------
	alloc := make(core.GenesisAlloc)
	alloc[userAddr] = core.GenesisAccount{Balance: userIntialBalance}
	blockchain := backends.NewSimulatedBackend(alloc)

	// ---------------------------------------------
	// COMPILE AND DEPLOY SIGNER
	// ---------------------------------------------
	randomKey1, _ := crypto.GenerateKey()
	randomKey2, _ := crypto.GenerateKey()
	randomAddr1 := crypto.PubkeyToAddress(randomKey1.PublicKey)
	randomAddr2 := crypto.PubkeyToAddress(randomKey2.PublicKey)

	signerABI, signerBytes := compileSignerContract()
	signers := []common.Address{userAddr, randomAddr1}
	signerContractTx := compileAndDeployContract(ctx, blockchain, userKey, signerBytes, signerABI, nil, uint64(3000000), signers)
	blockchain.Commit()
	signerContractAddr, err := bind.WaitDeployed(ctx, blockchain, signerContractTx)
	if err != nil {
		t.Fatal("ERROR while waiting for contract deployment")
	}

	var isSigner bool

	callContract(ctx, blockchain, signerABI, signerContractAddr, userAddr, "getSigner", &isSigner, randomAddr1)
	t.Log(isSigner)
	callContract(ctx, blockchain, signerABI, signerContractAddr, userAddr, "getSigner", &isSigner, randomAddr2)
	t.Log(isSigner)

}
