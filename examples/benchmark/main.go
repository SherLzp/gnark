package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	bls12381fr "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage is ./benchmark [nbConstraints list]")
		os.Exit(-1)
	}
	ns := strings.Split(os.Args[1], ",")
	curveIDs := []ecc.ID{ecc.BN254, ecc.BLS12_381}

	// write to stdout
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(benchData{}.headers()); err != nil {
		panic(err)
	}

	for _, curveID := range curveIDs {
		for _, _n := range ns {
			n, err := strconv.Atoi(_n)
			if err != nil {
				panic(err)
			}
			// generate dummy circuit and witness
			pk, r1cs := generateCircuit(n, curveID)
			witness := generateSolution(n, curveID)

			// measure proving time
			start := time.Now()
			// p := profile.Start(profile.TraceProfile, profile.ProfilePath("."), profile.NoShutdownHook)
			_, _ = groth16.Prove(r1cs, pk, &witness)
			// p.Stop()

			took := time.Since(start)

			// check memory usage, max ram requested from OS
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			bData := benchData{
				Curve:          curveID.String(),
				NbCores:        runtime.NumCPU(),
				NbCoefficients: r1cs.GetNbCoefficients(),
				NbConstraints:  r1cs.GetNbConstraints(),
				NbWires:        0, // TODO @gbotrel fixme
				RunTime:        took.Milliseconds(),
				MaxRAM:         (m.Sys / 1024 / 1024),
				Throughput:     int(float64(r1cs.GetNbConstraints()) / took.Seconds()),
			}
			bData.ThroughputPerCore = bData.Throughput / bData.NbCores

			if err := w.Write(bData.values()); err != nil {
				panic(err)
			}
			w.Flush()
		}
	}

}

// benchCircuit is a simple circuit that checks X*X*X*X*X... == Y
type benchCircuit struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
	n int
}

func (circuit *benchCircuit) Define(curveID ecc.ID, cs *frontend.ConstraintSystem) error {
	for i := 0; i < circuit.n; i++ {
		circuit.X = cs.Mul(circuit.X, circuit.X)
	}
	cs.AssertIsEqual(circuit.X, circuit.Y)
	return nil
}

func generateCircuit(nbConstraints int, curveID ecc.ID) (groth16.ProvingKey, frontend.CompiledConstraintSystem) {
	var circuit benchCircuit
	circuit.n = nbConstraints

	r1cs, err := frontend.Compile(curveID, backend.GROTH16, &circuit)
	if err != nil {
		panic(err)
	}

	// dummy setup will not compute a verifying key and just sets random value in the proving key
	pk, _ := groth16.DummySetup(r1cs)
	return pk, r1cs
}

func generateSolution(nbConstraints int, curveID ecc.ID) (witness benchCircuit) {
	witness.n = nbConstraints
	witness.X.Assign(2)

	switch curveID {
	case ecc.BN254:
		// compute expected Y
		var expectedY bn254fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < nbConstraints; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}

		witness.Y.Assign(expectedY)
	case ecc.BLS12_381:
		// compute expected Y
		var expectedY bls12381fr.Element
		expectedY.SetInterface(2)
		for i := 0; i < nbConstraints; i++ {
			expectedY.Mul(&expectedY, &expectedY)
		}

		witness.Y.Assign(expectedY)
	default:
		panic("not implemented")
	}

	return
}

type benchData struct {
	Curve             string
	NbConstraints     int
	NbWires           int
	NbCoefficients    int
	MaxRAM            uint64
	RunTime           int64
	NbCores           int
	Throughput        int
	ThroughputPerCore int
}

func (bData benchData) headers() []string {
	return []string{"curve", "nbConstraints", "nbWires", "nbCoefficients", "ram(mb)", "time(ms)", "nbCores", "throughput(constraints/s)", "througputPerCore(constraints/s)"}
}
func (bData benchData) values() []string {
	return []string{
		bData.Curve,
		strconv.Itoa(int(bData.NbConstraints)),
		strconv.Itoa(int(bData.NbWires)),
		strconv.Itoa(bData.NbCoefficients),
		strconv.Itoa(int(bData.MaxRAM)),
		strconv.Itoa(int(bData.RunTime)),
		strconv.Itoa(bData.NbCores),
		strconv.Itoa(bData.Throughput),
		strconv.Itoa(bData.ThroughputPerCore),
	}
}
