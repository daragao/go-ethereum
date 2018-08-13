package scpoa

import "testing"

func Test_CompileSigner(t *testing.T) {
	signerABI, signerBytes := compileSignerContract()
	t.Log(signerABI)
	t.Log(signerBytes)
}
