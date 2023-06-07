package emitter

import (
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/params"
	"github.com/skyhighblockchain/push-base/inter/idx"

	"github.com/skyhighblockchain/skyhigh/inter/validatorpk"
	"github.com/skyhighblockchain/skyhigh/skyhigh"
)

// EmitIntervals is the configuration of emit intervals.
type EmitIntervals struct {
	Min                        time.Duration
	Max                        time.Duration
	Confirming                 time.Duration // emit time when there's no txs to originate, but at least 1 tx to confirm
	ParallelInstanceProtection time.Duration
	DoublesignProtection       time.Duration
}

type ValidatorConfig struct {
	ID     idx.ValidatorID
	PubKey validatorpk.PubKey
}

type PrevEmittedEventFile struct {
	Path     string
	SyncMode bool
}

// Config is the configuration of events emitter.
type Config struct {
	VersionToPublish string

	Validator ValidatorConfig

	EmitIntervals EmitIntervals // event emission intervals

	MaxTxsPerAddress int

	MaxParents idx.Event

	// thresholds on GasLeft
	LimitedTpsThreshold uint64
	NoTxsThreshold      uint64
	EmergencyThreshold  uint64

	TxsCacheInvalidation time.Duration

	PrevEmittedEventFile PrevEmittedEventFile
}

// DefaultConfig returns the default configurations for the events emitter.
func DefaultConfig() Config {
	return Config{
		VersionToPublish: params.VersionWithMeta(),

		EmitIntervals: EmitIntervals{
			Min:                        110 * time.Millisecond,
			Max:                        10 * time.Minute,
			Confirming:                 120 * time.Millisecond,
			DoublesignProtection:       27 * time.Minute, // should be greater than MaxEmitInterval
			ParallelInstanceProtection: 1 * time.Minute,
		},

		MaxTxsPerAddress: TxTurnNonces,

		MaxParents: 0,

		LimitedTpsThreshold: skyhigh.DefaultEventGas * 120,
		NoTxsThreshold:      skyhigh.DefaultEventGas * 30,
		EmergencyThreshold:  skyhigh.DefaultEventGas * 5,

		TxsCacheInvalidation: 200 * time.Millisecond,
	}
}

// RandomizeEmitTime and return new config
func (cfg EmitIntervals) RandomizeEmitTime(r *rand.Rand) EmitIntervals {
	config := cfg
	// value = value - 0.1 * value + 0.1 * random value
	if config.Max > 10 {
		config.Max = config.Max - config.Max/10 + time.Duration(r.Int63n(int64(config.Max/10)))
	}
	// value = value + 0.33 * random value
	if config.DoublesignProtection > 3 {
		config.DoublesignProtection = config.DoublesignProtection + time.Duration(r.Int63n(int64(config.DoublesignProtection/3)))
	}
	return config
}

// FakeConfig returns the testing configurations for the events emitter.
func FakeConfig(num int) Config {
	cfg := DefaultConfig()
	cfg.EmitIntervals.Max = 10 * time.Second // don't wait long in fakenet
	cfg.EmitIntervals.DoublesignProtection = cfg.EmitIntervals.Max / 2
	if num <= 1 {
		// disable self-fork protection if fakenet 1/1
		cfg.EmitIntervals.DoublesignProtection = 0
	}
	return cfg
}
