package largewitness

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
)

const Size = 2000000

type Circuit struct {
	P [Size]frontend.Variable `gnark:",public"`
	Q frontend.Variable
}

func (circuit *Circuit) Define(curveID ecc.ID, cs *frontend.ConstraintSystem) error {
	var _p [Size]interface{}
	for i := 0; i < Size; i++ {
		_p[i] = circuit.P[i]
	}

	sum := cs.Add(_p[0], _p[1], _p[2:]...)

	cs.AssertIsEqual(circuit.Q, sum)

	return nil
}
