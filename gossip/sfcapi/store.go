package sfcapi

import (
	"github.com/skyhighblockchain/push-base/kvdb"
	"github.com/skyhighblockchain/push-base/kvdb/table"

	"github.com/skyhighblockchain/skyhigh/logger"
	"github.com/skyhighblockchain/skyhigh/utils/rlpstore"
)

// Store is a node persistent storage working over physical key-value database.
type Store struct {
	mainDB kvdb.Store
	table  struct {
		GasPowerRefund kvdb.Store `table:"R"`

		Validators  kvdb.Store `table:"1"`
		Stakers     kvdb.Store `table:"2"`
		Delegations kvdb.Store `table:"3"`

		DelegationOldRewards        kvdb.Store `table:"6"`
		StakerOldRewards            kvdb.Store `table:"7"`
		StakerDelegationsOldRewards kvdb.Store `table:"8"`
	}

	rlp rlpstore.Helper

	logger.Instance
}

// NewStore creates store over key-value db.
func NewStore(mainDB kvdb.Store) *Store {
	s := &Store{
		mainDB:   mainDB,
		Instance: logger.MakeInstance(),
		rlp:      rlpstore.Helper{logger.MakeInstance()},
	}

	table.MigrateTables(&s.table, s.mainDB)

	return s
}

// Close closes underlying database.
func (s *Store) Close() {
	table.MigrateTables(&s.table, nil)

	_ = s.mainDB.Close()
}
