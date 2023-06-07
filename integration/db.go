package integration

import (
	"errors"
	"strings"

	"github.com/skyhighblockchain/push-base/hash"
	"github.com/skyhighblockchain/push-base/inter/dag"
	"github.com/skyhighblockchain/push-base/kvdb"
	"github.com/skyhighblockchain/push-base/kvdb/leveldb"
	"github.com/skyhighblockchain/push-base/kvdb/memorydb"
	"github.com/skyhighblockchain/push-base/utils/cachescale"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/skyhighblockchain/skyhigh/gossip"
)

func DBProducer(chaindataDir string, scale cachescale.Func) kvdb.IterableDBProducer {
	if chaindataDir == "inmemory" || chaindataDir == "" {
		return memorydb.NewProducer("")
	}

	return leveldb.NewProducer(chaindataDir, func(name string) int {
		return dbCacheSize(name, scale.I)
	})
}

func CheckDBList(names []string) error {
	if len(names) == 0 {
		return nil
	}
	namesMap := make(map[string]bool)
	for _, name := range names {
		namesMap[name] = true
	}
	if !namesMap["gossip"] {
		return errors.New("gossip DB is not found")
	}
	if !namesMap["push"] {
		return errors.New("push DB is not found")
	}
	if !namesMap["genesis"] {
		return errors.New("genesis DB is not found")
	}
	return nil
}

func dbCacheSize(name string, scale func(int) int) int {
	if name == "gossip" {
		return scale(128 * opt.MiB)
	}
	if name == "push" {
		return scale(4 * opt.MiB)
	}
	if strings.HasPrefix(name, "push-") {
		return scale(8 * opt.MiB)
	}
	if strings.HasPrefix(name, "gossip-") {
		return scale(8 * opt.MiB)
	}
	return scale(2 * opt.MiB)
}

func dropAllDBs(producer kvdb.IterableDBProducer) {
	names := producer.Names()
	for _, name := range names {
		db, err := producer.OpenDB(name)
		if err != nil {
			continue
		}
		_ = db.Close()
		db.Drop()
	}
}

func dropAllDBsIfInterrupted(rawProducer kvdb.IterableDBProducer) {
	names := rawProducer.Names()
	if len(names) == 0 {
		return
	}
	// if flushID is not written, then previous genesis processing attempt was interrupted
	for _, name := range names {
		db, err := rawProducer.OpenDB(name)
		if err != nil {
			return
		}
		flushID, err := db.Get(FlushIDKey)
		_ = db.Close()
		if flushID != nil || err != nil {
			return
		}
	}
	dropAllDBs(rawProducer)
}

type GossipStoreAdapter struct {
	*gossip.Store
}

func (g *GossipStoreAdapter) GetEvent(id hash.Event) dag.Event {
	e := g.Store.GetEvent(id)
	if e == nil {
		return nil
	}
	return e
}

type DummyFlushableProducer struct {
	kvdb.DBProducer
}

func (p *DummyFlushableProducer) NotFlushedSizeEst() int {
	return 0
}

func (p *DummyFlushableProducer) Flush(_ []byte) error {
	return nil
}
