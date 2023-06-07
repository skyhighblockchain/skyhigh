package evmstore

import (
	"github.com/skyhighblockchain/push-base/hash"
	"github.com/skyhighblockchain/push-base/kvdb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/skyhighblockchain/skyhigh/evmcore"
	"github.com/skyhighblockchain/skyhigh/skyhigh"
)

func (s *Store) applyRawEvmItems(db kvdb.Iteratee) (err error) {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	batch := s.table.Evm.NewBatch()
	defer batch.Reset()
	for it.Next() {
		err = batch.Put(it.Key(), it.Value())
		if err != nil {
			return err
		}
		if batch.ValueSize() > kvdb.IdealBatchSize {
			err = batch.Write()
			if err != nil {
				return err
			}
			batch.Reset()
		}
	}
	return batch.Write()
}

// ApplyGenesis writes initial state.
func (s *Store) ApplyGenesis(g skyhigh.Genesis, startingRoot hash.Hash) (evmBlock *evmcore.EvmBlock, err error) {
	// apply raw EVM storage
	err = s.applyRawEvmItems(g.RawEvmItems)
	if err != nil {
		return nil, err
	}
	// state
	statedb, err := s.StateDB(startingRoot)
	if err != nil {
		return nil, err
	}
	return evmcore.ApplyGenesis(statedb, g, 128*opt.MiB)
}
