package sealmodule

import (
	"math/big"

	"github.com/skyhighblockchain/push-base/inter/idx"
	"github.com/skyhighblockchain/push-base/inter/pos"
	"github.com/skyhighblockchain/push-base/push"

	"github.com/skyhighblockchain/skyhigh/gossip/blockproc"
)

type SkyhighEpochsSealerModule struct{}

func New() *SkyhighEpochsSealerModule {
	return &SkyhighEpochsSealerModule{}
}

func (m *SkyhighEpochsSealerModule) Start(block blockproc.BlockCtx, bs blockproc.BlockState, es blockproc.EpochState) blockproc.SealerProcessor {
	return &SkyhighEpochsSealer{
		block: block,
		es:    es,
		bs:    bs,
	}
}

type SkyhighEpochsSealer struct {
	block blockproc.BlockCtx
	es    blockproc.EpochState
	bs    blockproc.BlockState
}

func (s *SkyhighEpochsSealer) EpochSealing() bool {
	sealEpoch := s.bs.EpochGas >= s.es.Rules.Epochs.MaxEpochGas
	sealEpoch = sealEpoch || (s.block.Time-s.es.EpochStart) >= s.es.Rules.Epochs.MaxEpochDuration
	sealEpoch = sealEpoch || s.bs.AdvanceEpochs > 0
	return sealEpoch || s.bs.EpochCheaters.Len() != 0
}

func (p *SkyhighEpochsSealer) Update(bs blockproc.BlockState, es blockproc.EpochState) {
	p.bs, p.es = bs, es
}

// SealEpoch is called after pre-internal transactions are executed
func (s *SkyhighEpochsSealer) SealEpoch() (blockproc.BlockState, blockproc.EpochState) {
	// Select new validators
	oldValidators := s.es.Validators
	builder := pos.NewBigBuilder()
	for v, profile := range s.bs.NextValidatorProfiles {
		builder.Set(v, profile.Weight)
	}
	newValidators := builder.Build()
	s.es.Validators = newValidators
	s.es.ValidatorProfiles = s.bs.NextValidatorProfiles.Copy()

	// Build new []ValidatorEpochState and []ValidatorBlockState
	newValidatorEpochStates := make([]blockproc.ValidatorEpochState, newValidators.Len())
	newValidatorBlockStates := make([]blockproc.ValidatorBlockState, newValidators.Len())
	for newValIdx := idx.Validator(0); newValIdx < newValidators.Len(); newValIdx++ {
		// default values
		newValidatorBlockStates[newValIdx] = blockproc.ValidatorBlockState{
			Originated: new(big.Int),
		}
		// inherit validator's state from previous epoch, if any
		valID := newValidators.GetID(newValIdx)
		if !oldValidators.Exists(valID) {
			// new validator
			newValidatorBlockStates[newValIdx].LastBlock = s.block.Idx
			newValidatorBlockStates[newValIdx].LastOnlineTime = s.block.Time
			continue
		}
		oldValIdx := oldValidators.GetIdx(valID)
		newValidatorBlockStates[newValIdx] = s.bs.ValidatorStates[oldValIdx]
		newValidatorBlockStates[newValIdx].DirtyGasRefund = 0
		newValidatorBlockStates[newValIdx].Uptime = 0
		newValidatorEpochStates[newValIdx].GasRefund = s.bs.ValidatorStates[oldValIdx].DirtyGasRefund
		newValidatorEpochStates[newValIdx].PrevEpochEvent = s.bs.ValidatorStates[oldValIdx].LastEvent
	}
	s.es.ValidatorStates = newValidatorEpochStates
	s.bs.ValidatorStates = newValidatorBlockStates
	s.es.Validators = newValidators

	// dirty data becomes active
	s.es.PrevEpochStart = s.es.EpochStart
	s.es.EpochStart = s.block.Time
	s.es.Rules = s.bs.DirtyRules.Copy()
	s.es.EpochStateRoot = s.bs.FinalizedStateRoot

	s.bs.EpochGas = 0
	s.bs.EpochCheaters = push.Cheaters{}
	newEpoch := s.es.Epoch + 1
	s.es.Epoch = newEpoch

	if s.bs.AdvanceEpochs > 0 {
		s.bs.AdvanceEpochs--
	}

	return s.bs, s.es
}
