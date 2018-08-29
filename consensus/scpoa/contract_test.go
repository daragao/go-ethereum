package scpoa

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func Test_CompileSigner(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
	callContract(ctx, blockchain, signerABI, userAddr, signerContractAddr, "getSigner", &isSigner, randomAddr1)
	if !isSigner {
		t.Fatal("Does not have random address1")
	}
	callContract(ctx, blockchain, signerABI, userAddr, signerContractAddr, "getSigner", &isSigner, randomAddr2)
	if isSigner {
		t.Fatal("Has random address2")
	}

}

func Test_AddTxToSCPoa(t *testing.T) {
	randomKey1, _ := crypto.GenerateKey()
	randomKey2, _ := crypto.GenerateKey()
	randomAddr1 := crypto.PubkeyToAddress(randomKey1.PublicKey)
	randomAddr2 := crypto.PubkeyToAddress(randomKey2.PublicKey)

	signers := []common.Address{randomAddr1, randomAddr2}
	// extraVanity global variable 32
	// extraSeal global variable 65
	// Create the genesis block with the initial set of signers
	genesis := &core.Genesis{
		ExtraData: make([]byte, extraVanity+common.AddressLength*len(signers)+extraSeal),
	}
	for j, signer := range signers {
		copy(genesis.ExtraData[extraVanity+j*common.AddressLength:], signer[:])
	}
	// Create a pristine blockchain with the genesis injected
	db := ethdb.NewMemDatabase()
	genesis.Commit(db)

	scpoaConfig := params.SCPoaConfig{Epoch: 3, Period: 10}
	scpoa := New(&scpoaConfig, db)

	// Assemble a chain of headers from the cast votes
	header := &types.Header{
		Number: big.NewInt(0),
		Time:   big.NewInt(0),
		Extra:  make([]byte, extraVanity+extraSeal),
		//Coinbase: randomAddr1,
	}

	//headers[j].ParentHash = headers[j-1].Hash()

	// sign header and append signature to extra data
	sig, err := crypto.Sign(sigHash(header).Bytes(), randomKey1)
	if err != nil {
		t.Fatal("Failed to sign header:", err)
	}
	copy(header.Extra[len(header.Extra)-65:], sig)

	headers := []*types.Header{header}

	chainReader := testerChainReader{db: db}
	snap, err := scpoa.snapshot(
		&chainReader,
		header.Number.Uint64(),
		header.Hash(),
		headers,
	)
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	header0 := chainReader.GetHeaderByNumber(0)
	block0 := chainReader.GetBlock(header0.Hash(), 0)
	t.Logf("header: %+v\nsnapshot: %+v", block0, snap)
}
