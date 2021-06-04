package main

import (
	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/examples/largewitness"
	"github.com/consensys/gnark/frontend"
)

//go:generate go run generate.go
func main() {
	var circuit largewitness.Circuit
	r1cs, _ := frontend.Compile(ecc.BN254, backend.GROTH16, &circuit)
	const name = "large"

	circuitDir := filepath.Join("bn254", name)
	os.MkdirAll(circuitDir, 0777)

	{
		f, _ := os.Create(filepath.Join(circuitDir, name+".r1cs"))
		r1cs.WriteTo(f)
		f.Close()
	}

	pk, vk, _ := groth16.Setup(r1cs)
	{
		f, _ := os.Create(filepath.Join(circuitDir, name+".pk"))
		pk.WriteTo(f)
		f.Close()
	}
	{
		f, _ := os.Create(filepath.Join(circuitDir, name+".vk"))
		vk.WriteTo(f)
		f.Close()
	}
}
