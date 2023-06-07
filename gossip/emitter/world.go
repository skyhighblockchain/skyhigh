package emitter

import (
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	notify "github.com/ethereum/go-ethereum/event"
	"github.com/skyhighblockchain/push-base/hash"
	"github.com/skyhighblockchain/push-base/inter/idx"
	"github.com/skyhighblockchain/push-base/inter/pos"

	"github.com/skyhighblockchain/skyhigh/evmcore"
	"github.com/skyhighblockchain/skyhigh/inter"
	"github.com/skyhighblockchain/skyhigh/skyhigh"
	"github.com/skyhighblockchain/skyhigh/valkeystore"
	"github.com/skyhighblockchain/skyhigh/vecmt"
)

var (
	ErrNotEnoughGasPower = errors.New("not enough gas power")
)

type (
	// External world
	External interface {
		sync.Locker
		Reader

		Check(e *inter.EventPayload, parents inter.Events) error
		Process(*inter.EventPayload) error
		Broadcast(*inter.EventPayload)
		Build(*inter.MutableEventPayload, func()) error
		DagIndex() *vecmt.Index

		IsBusy() bool
		IsSynced() bool
		PeersNum() int
	}

	// aliases for mock generator
	Signer   valkeystore.SignerI
	TxSigner types.Signer

	// World is an emitter's environment
	World struct {
		External
		TxPool   TxPool
		Signer   valkeystore.SignerI
		TxSigner types.Signer
	}
)

// Reader is a callback for getting events from an external storage.
type Reader interface {
	GetLatestBlockIndex() idx.Block
	GetEpochValidators() (*pos.Validators, idx.Epoch)
	GetEvent(hash.Event) *inter.Event
	GetEventPayload(hash.Event) *inter.EventPayload
	GetLastEvent(epoch idx.Epoch, from idx.ValidatorID) *hash.Event
	GetHeads(idx.Epoch) hash.Events
	GetGenesisTime() inter.Timestamp
	GetRules() skyhigh.Rules
	GetRecommendedGasPrice() *big.Int
}

type TxPool interface {
	// Has returns an indicator whether txpool has a transaction cached with the
	// given hash.
	Has(hash common.Hash) bool
	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending() (map[common.Address]types.Transactions, error)

	// SubscribeNewTxsNotify should return an event subscription of
	// NewTxsNotify and send events to the given channel.
	SubscribeNewTxsNotify(chan<- evmcore.NewTxsNotify) notify.Subscription

	// Count returns the total number of transactions
	Count() int
}
