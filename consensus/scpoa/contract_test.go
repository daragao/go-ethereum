package scpoa

import (
	"context"
	"math/big"
	"testing"

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
	signerABI, signerBytes := compileSignerContract()

	signers := []common.Address{userAddr}
	compileAndDeployContract(ctx, blockchain, userKey, signerBytes, signerABI, nil, uint64(3000000), signers)

}
