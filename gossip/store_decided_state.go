package gossip

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/skyhighblockchain/push-base/inter/idx"
	"github.com/skyhighblockchain/push-base/inter/pos"

	"github.com/skyhighblockchain/skyhigh/gossip/blockproc"
	"github.com/skyhighblockchain/skyhigh/skyhigh"
)

const sKey = "s"

type BlockEpochState struct {
	BlockState *blockproc.BlockState
	EpochState *blockproc.EpochState
}

// SetBlockEpochState stores the latest block and epoch state in memory
func (s *Store) SetBlockEpochState(bs blockproc.BlockState, es blockproc.EpochState) {
	bs, es = bs.Copy(), es.Copy()
	s.cache.BlockEpochState.Store(&BlockEpochState{&bs, &es})
}

func (s *Store) getBlockEpochState() BlockEpochState {
	if v := s.cache.BlockEpochState.Load(); v != nil {
		return *v.(*BlockEpochState)
	}
	v, ok := s.rlp.Get(s.table.BlockEpochState, []byte(sKey), &BlockEpochState{}).(*BlockEpochState)
	if !ok {
		log.Crit("Epoch state reading failed: genesis not applied")
	}
	s.cache.BlockEpochState.Store(v)
	return *v
}

// FlushBlockEpochState stores the latest epoch and block state in DB
func (s *Store) FlushBlockEpochState() {
	s.rlp.Set(s.table.BlockEpochState, []byte(sKey), s.getBlockEpochState())
}

// GetBlockState retrieves the latest block state
func (s *Store) GetBlockState() blockproc.BlockState {
	return *s.getBlockEpochState().BlockState
}

// GetEpochState retrieves the latest epoch state
func (s *Store) GetEpochState() blockproc.EpochState {
	return *s.getBlockEpochState().EpochState
}

func (s *Store) GetBlockEpochState() (blockproc.BlockState, blockproc.EpochState) {
	bes := s.getBlockEpochState()
	return *bes.BlockState, *bes.EpochState
}

// GetEpoch retrieves the current epoch
func (s *Store) GetEpoch() idx.Epoch {
	return s.GetEpochState().Epoch
}

// GetValidators retrieves current validators
func (s *Store) GetValidators() *pos.Validators {
	return s.GetEpochState().Validators
}

// GetEpochValidators retrieves the current epoch and validators atomically
func (s *Store) GetEpochValidators() (*pos.Validators, idx.Epoch) {
	es := s.GetEpochState()
	return es.Validators, es.Epoch
}

// GetLatestBlockIndex retrieves the current block number
func (s *Store) GetLatestBlockIndex() idx.Block {
	return s.GetBlockState().LastBlock.Idx
}

// GetRules retrieves current network rules
func (s *Store) GetRules() skyhigh.Rules {
	return s.GetEpochState().Rules
}

// GetEpochRules retrieves current network rules and epoch atomically
func (s *Store) GetEpochRules() (skyhigh.Rules, idx.Epoch) {
	es := s.GetEpochState()
	return es.Rules, es.Epoch
}
