package antichamber

import (
	"github.com/consensys/zkevm-monorepo/prover/crypto/keccak"
	"github.com/consensys/zkevm-monorepo/prover/maths/common/smartvectors"
	"github.com/consensys/zkevm-monorepo/prover/maths/common/vector"
	"github.com/consensys/zkevm-monorepo/prover/maths/field"
	"github.com/consensys/zkevm-monorepo/prover/protocol/ifaces"
	"github.com/consensys/zkevm-monorepo/prover/protocol/wizard"
	"github.com/consensys/zkevm-monorepo/prover/zkevm/prover/hash/generic"
)

// TxSignature is responsible for assigning the relevant columns for transaction-Hash,
// and checking their consistency with the data coming from rlp_txn.
//
// columns for transaction-hash are native columns,
//
// columns for rlp-txn lives on the arithmetization side.
type txSignature struct {
	// we dont need it since order is preserved by projection.
	// txID     ifaces.Column
	txHashHi ifaces.Column
	txHashLo ifaces.Column
	isTxHash ifaces.Column

	// provider for keccak, Provider contain the inputs and outputs of keccak hash.
	provider generic.GenericByteModule
}

func newTxSignatures(comp *wizard.CompiledIOP, rlpTxn generic.GenDataModule, size int) *txSignature {
	createCol := createColFn(comp, NAME_TXSIGNATURE, size)

	// declare the native columns
	res := &txSignature{
		txHashHi: createCol("TX_HASH_HI"),
		txHashLo: createCol("TX_HASH_LO"),
		isTxHash: createCol("TX_IS_HASH_HI"),
	}

	res.provider = res.GetProvider(comp, rlpTxn)

	return res
}

// It builds a provider from rlp-txn (as hash input) and native columns of TxSignature (as hash output)
// the consistency check is then deferred to the keccak module.
func (txn *txSignature) GetProvider(comp *wizard.CompiledIOP, rlpTxn generic.GenDataModule) generic.GenericByteModule {
	provider := generic.GenericByteModule{}

	// pass rlp-txn as DataModule.
	provider.Data = rlpTxn

	// generate infoModule from native columns
	provider.Info = txn.buildInfoModule()

	return provider
}

// it builds an infoModule from native columns
func (txn *txSignature) buildInfoModule() generic.GenInfoModule {
	info := generic.GenInfoModule{
		// HashNum:   txn.txID,
		HashHi:   txn.txHashHi,
		HashLo:   txn.txHashLo,
		IsHashHi: txn.isTxHash,
		IsHashLo: txn.isTxHash,
	}
	return info
}

// it assign the native columns
func (txn *txSignature) assignTxSignature(run *wizard.ProverRuntime, rlpTxn generic.GenDataModule, nbEcRecover, size int) {
	n := startAt(nbEcRecover)

	hashHi := vector.Repeat(field.Zero(), n)
	hashLo := vector.Repeat(field.Zero(), n)
	isTxHash := vector.Repeat(field.Zero(), n)

	permTrace := keccak.GenerateTrace(rlpTxn.ScanStreams(run))

	var v, w field.Element
	for _, digest := range permTrace.HashOutPut {
		hi := digest[:halfDigest]
		lo := digest[halfDigest:]

		v.SetBytes(hi[:])
		w.SetBytes(lo[:])

		repeatLo := vector.Repeat(w, nbRowsPerTxSign)
		repeatHi := vector.Repeat(v, nbRowsPerTxSign)
		repeatIsTxHash := vector.Repeat(field.Zero(), nbRowsPerTxSign-1)

		hashHi = append(hashHi, repeatHi...)
		hashLo = append(hashLo, repeatLo...)
		isTxHash = append(isTxHash, field.One())
		isTxHash = append(isTxHash, repeatIsTxHash...)
	}

	run.AssignColumn(txn.txHashHi.GetColID(), smartvectors.RightZeroPadded(hashHi, size))
	run.AssignColumn(txn.txHashLo.GetColID(), smartvectors.RightZeroPadded(hashLo, size))
	run.AssignColumn(txn.isTxHash.GetColID(), smartvectors.RightZeroPadded(isTxHash, size))
}

func startAt(nbEcRecover int) int {
	return nbEcRecover * nbRowsPerEcRec
}
