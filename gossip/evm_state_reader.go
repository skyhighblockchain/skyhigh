package gossip

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/skyhighblockchain/push-base/hash"
	"github.com/skyhighblockchain/push-base/inter/idx"

	"github.com/skyhighblockchain/skyhigh/evmcore"
	"github.com/skyhighblockchain/skyhigh/gossip/gasprice"
	"github.com/skyhighblockchain/skyhigh/inter"
)

var (
	big3 = big.NewInt(3)
	big4 = big.NewInt(4)
)

type EvmStateReader struct {
	*ServiceFeed

	store *Store
	gpo   *gasprice.Oracle
}

func (s *Service) GetEvmStateReader() *EvmStateReader {
	return &EvmStateReader{
		ServiceFeed: &s.feed,
		store:       s.store,
		gpo:         s.gpo,
	}
}

// MinGasPrice returns current hard lower bound for gas price
func (r *EvmStateReader) MinGasPrice() *big.Int {
	return r.store.GetRules().Economy.MinGasPrice
}

// RecommendedMinGasPrice returns current soft lower bound for gas price
func (r *EvmStateReader) RecommendedMinGasPrice() *big.Int {
	est := new(big.Int).Set(r.gpo.SuggestPrice())
	est.Mul(est, big3)
	est.Div(est, big4)

	min := r.MinGasPrice()
	if min.Cmp(est) > 0 {
		return min
	}

	return est
}

func (r *EvmStateReader) MaxGasLimit() uint64 {
	rules := r.store.GetRules()
	maxEmptyEventGas := rules.Economy.Gas.EventGas +
		uint64(rules.Dag.MaxParents-rules.Dag.MaxFreeParents)*rules.Economy.Gas.ParentGas +
		uint64(rules.Dag.MaxExtraData)*rules.Economy.Gas.ExtraDataGas
	if rules.Economy.Gas.MaxEventGas < maxEmptyEventGas {
		return 0
	}
	return rules.Economy.Gas.MaxEventGas - maxEmptyEventGas
}

func (r *EvmStateReader) Config() *params.ChainConfig {
	return r.store.GetRules().EvmChainConfig()
}

func (r *EvmStateReader) CurrentBlock() *evmcore.EvmBlock {
	n := r.store.GetLatestBlockIndex()

	return r.getBlock(hash.Event{}, n, true)
}

func (r *EvmStateReader) CurrentHeader() *evmcore.EvmHeader {
	n := r.store.GetLatestBlockIndex()

	return r.getBlock(hash.Event{}, n, false).Header()
}

func (r *EvmStateReader) GetHeader(h common.Hash, n uint64) *evmcore.EvmHeader {
	return r.getBlock(hash.Event(h), idx.Block(n), false).Header()
}

func (r *EvmStateReader) GetBlock(h common.Hash, n uint64) *evmcore.EvmBlock {
	return r.getBlock(hash.Event(h), idx.Block(n), true)
}

func (r *EvmStateReader) getBlock(h hash.Event, n idx.Block, readTxs bool) *evmcore.EvmBlock {
	block := r.store.GetBlock(n)
	if block == nil {
		return nil
	}
	if (h != hash.Event{}) && (h != block.Atropos) {
		return nil
	}
	if readTxs {
		if cached := r.store.EvmStore().GetCachedEvmBlock(n); cached != nil {
			return cached
		}
	}

	var transactions types.Transactions
	if readTxs {
		transactions = make(types.Transactions, 0, len(block.Txs)+len(block.InternalTxs)+len(block.Events)*10)
		for _, txid := range block.InternalTxs {
			tx := r.store.evm.GetTx(txid)
			if tx == nil {
				log.Crit("Internal tx not found", "tx", txid.String())
				continue
			}
			transactions = append(transactions, tx)
		}
		for _, txid := range block.Txs {
			tx := r.store.evm.GetTx(txid)
			if tx == nil {
				log.Crit("Tx not found", "tx", txid.String())
				continue
			}
			transactions = append(transactions, tx)
		}
		for _, id := range block.Events {
			e := r.store.GetEventPayload(id)
			if e == nil {
				log.Crit("Block event not found", "event", id.String())
				continue
			}
			transactions = append(transactions, e.Txs()...)
		}

		transactions = inter.FilterSkippedTxs(transactions, block.SkippedTxs)
	} else {
		transactions = make(types.Transactions, 0)
	}

	var prev hash.Event
	if n != 0 {
		prev = r.store.GetBlock(n - 1).Atropos
	}
	evmHeader := evmcore.ToEvmHeader(block, n, prev)

	var evmBlock *evmcore.EvmBlock
	if readTxs {
		evmBlock = evmcore.NewEvmBlock(evmHeader, transactions)
		r.store.EvmStore().SetCachedEvmBlock(n, evmBlock)
	} else {
		// not completed block here
		evmBlock = &evmcore.EvmBlock{
			EvmHeader: *evmHeader,
		}
	}

	return evmBlock
}

func (r *EvmStateReader) StateAt(root common.Hash) (*state.StateDB, error) {
	return r.store.evm.StateDB(hash.Hash(root))
}

func (r *EvmStateReader) TxExists(txid common.Hash) bool {
	return r.store.EvmStore().GetTxPosition(txid) != nil
}
