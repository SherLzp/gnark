package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/examples/largewitness"
	"github.com/consensys/gnark/frontend"
)

const circuitName = "large"

//go:generate go run generate.go
func main() {
	log.Println("starting...")
	defer log.Println("done")
	circuitDir := filepath.Join("bn254", circuitName)
	os.MkdirAll(circuitDir, 0777)

	log.Println("compiling", circuitName)
	var circuit largewitness.Circuit
	r1cs, _ := frontend.Compile(ecc.BN254, backend.GROTH16, &circuit)

	{
		f, _ := os.Create(filepath.Join(circuitDir, circuitName+".r1cs"))
		r1cs.WriteTo(f)
		f.Close()
	}

	log.Println("groth16.Setup()")
	pk, vk, _ := groth16.Setup(r1cs)
	{
		f, _ := os.Create(filepath.Join(circuitDir, circuitName+".pk"))
		pk.WriteTo(f)
		f.Close()
	}
	{
		f, _ := os.Create(filepath.Join(circuitDir, circuitName+".vk"))
		vk.WriteTo(f)
		f.Close()
	}
}
