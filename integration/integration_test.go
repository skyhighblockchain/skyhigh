package integration

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"

	"github.com/skyhighblockchain/skyhigh/integration/makegenesis"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesisstore"
)

type topology func(net *simulations.Network, nodes []enode.ID)

func TestStar(t *testing.T) {
	testSim(t, topologyStar)
}

func TestRing(t *testing.T) {
	testSim(t, topologyRing)
}

var registerGossip sync.Once

func testSim(t *testing.T, connect topology) {
	const count = 3

	// set the log level to Trace
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace,
		log.StreamHandler(os.Stderr, log.TerminalFormat(false))))

	// fake net
	fakeGenesisStore := makegenesis.FakeGenesisStore(count, big.NewInt(0), big.NewInt(10000))
	genesis := InputGenesis{
		Hash: fakeGenesisStore.Hash(),
		Read: func(store *genesisstore.Store) error {
			buf := bytes.NewBuffer(nil)
			err := fakeGenesisStore.Export(buf)
			if err != nil {
				return err
			}
			return store.Import(buf)
		},
		Close: func() error {
			return nil
		},
	}

	// register a single gossip service
	services := adapters.LifecycleConstructors{
		"gossip": func(ctx *adapters.ServiceContext, stack *node.Node) (node.Lifecycle, error) {
			g := NewIntegration(ctx, genesis, stack)
			return g, nil
		},
	}
	registerGossip.Do(func() {
		adapters.RegisterLifecycles(services)
	})

	// create the NodeAdapter
	var adapter adapters.NodeAdapter
	adapter = adapters.NewSimAdapter(services)

	// create network
	sim := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: serviceNames(services)[0],
	})

	// create and start nodes
	nodes := make([]enode.ID, count)
	for i, val := range fakeGenesisStore.GetMetadata().Validators {
		key := makegenesis.FakeKey(int(val.ID))
		id := enode.PubkeyToIDV4(&key.PublicKey)
		config := &adapters.NodeConfig{
			ID:         id,
			Name:       fmt.Sprintf("Node-%d", i),
			PrivateKey: key,
			Lifecycles: serviceNames(services),
		}

		_, err := sim.NewNodeWithConfig(config)
		if err != nil {
			panic(err)
		}

		nodes[i] = id
	}

	sim.StartAll()
	defer sim.Shutdown()

	connect(sim, nodes)

	// start
	srv := &http.Server{
		Addr:    ":8888",
		Handler: simulations.NewServer(sim),
	}
	go func() {
		log.Info("Starting simulation server on 0.0.0.0:8888...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Crit("Error starting simulation server", "err", err)
		}
	}()

	// stop
	<-time.After(5 * time.Second)

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Crit("Error stopping simulation server", "err", err)
	}
}

func topologyStar(net *simulations.Network, nodes []enode.ID) {
	if len(nodes) < 2 {
		return
	}
	err := net.ConnectNodesStar(nodes, nodes[0])
	if err != nil {
		panic(err)
	}
}

func topologyRing(net *simulations.Network, nodes []enode.ID) {
	if len(nodes) < 2 {
		return
	}
	err := net.ConnectNodesRing(nodes)
	if err != nil {
		panic(err)
	}
}

func serviceNames(services adapters.LifecycleConstructors) []string {
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}

	return names
}
