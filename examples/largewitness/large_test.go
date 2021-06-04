package largewitness

import (
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

func TestLargeWitness(t *testing.T) {
	assert := groth16.NewAssert(t)

	var circuit Circuit

	// compiles our circuit into a R1CS
	r1cs, err := frontend.Compile(ecc.BN254, backend.GROTH16, &circuit)
	assert.NoError(err)

	{
		var witness Circuit
		for i := 0; i < Size; i++ {
			witness.P[i].Assign(i)
		}
		witness.Q.Assign(42)

		assert.ProverFailed(r1cs, &witness)
	}

	{
		var witness Circuit
		for i := 0; i < Size; i++ {
			witness.P[i].Assign(1)
		}
		witness.Q.Assign(Size)

		assert.ProverSucceeded(r1cs, &witness)
	}

}
