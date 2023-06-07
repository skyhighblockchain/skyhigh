package gossip

// compile SFC with truffle
//go:generate bash -c "cd ../../skyhigh-sfc && git checkout c1d33c81f74abf82c0e22807f16e609578e10ad8"
//go:generate bash -c "docker run --name skyhigh-sfc-compiler -v $(pwd)/contract/solc:/src/build/contracts -v $(pwd)/../../skyhigh-sfc:/src -w /src node:10.5.0 bash -c 'export NPM_CONFIG_PREFIX=~; npm install --no-save; npm install --no-save truffle@5.1.4' && docker commit skyhigh-sfc-compiler skyhigh-sfc-compiler-image && docker rm skyhigh-sfc-compiler"
//go:generate bash -c "docker run --rm -v $(pwd)/contract/solc:/src/build/contracts -v $(pwd)/../../skyhigh-sfc:/src -w /src skyhigh-sfc-compiler-image bash -c 'export NPM_CONFIG_PREFIX=~; rm -f /src/build/contracts/*json; npm run build'"
//go:generate bash -c "cd ./contract/solc && for f in LegacySfcWrapper.json; do jq -j .bytecode $DOLLAR{f} > $DOLLAR{f%.json}.bin; jq -j .deployedBytecode $DOLLAR{f} > $DOLLAR{f%.json}.bin-runtime; jq -c .abi $DOLLAR{f} > $DOLLAR{f%.json}.abi; done"
//go:generate bash -c "docker run --rm -v $(pwd)/contract/solc:/src/build/contracts -v $(pwd)/../../skyhigh-sfc:/src -w /src skyhigh-sfc-compiler-image bash -c 'export NPM_CONFIG_PREFIX=~; sed -i s/runs:\\ 200,/runs:\\ 10000,/ /src/truffle-config.js; rm -f /src/build/contracts/*json; npm run build'"
//go:generate bash -c "cd ../../skyhigh-sfc && git checkout -- truffle-config.js; docker rmi skyhigh-sfc-compiler-image"
//go:generate bash -c "cd ./contract/solc && for f in NetworkInitializer.json NodeDriver.json NodeDriverAuth.json; do jq -j .bytecode $DOLLAR{f} > $DOLLAR{f%.json}.bin; jq -j .deployedBytecode $DOLLAR{f} > $DOLLAR{f%.json}.bin-runtime; jq -c .abi $DOLLAR{f} > $DOLLAR{f%.json}.abi; done"

// wrap LegacySfcWrapper with golang
//go:generate mkdir -p ./contract/sfc100
//go:generate go run github.com/ethereum/go-ethereum/cmd/abigen --bin=./contract/solc/LegacySfcWrapper.bin --abi=./contract/solc/LegacySfcWrapper.abi --pkg=sfc100 --type=Contract --out=contract/sfc100/contract.go
//go:generate bash -c "(echo -ne '\nvar ContractBinRuntime = \"'; cat contract/solc/LegacySfcWrapper.bin-runtime; echo '\"') >> contract/sfc100/contract.go"
// wrap NetworkInitializer with golang
//go:generate mkdir -p ./contract/netinit100
//go:generate go run github.com/ethereum/go-ethereum/cmd/abigen --bin=./contract/solc/NetworkInitializer.bin --abi=./contract/solc/NetworkInitializer.abi --pkg=netinit100 --type=Contract --out=contract/netinit100/contract.go
//go:generate bash -c "(echo -ne '\nvar ContractBinRuntime = \"'; cat contract/solc/NetworkInitializer.bin-runtime; echo '\"') >> contract/netinit100/contract.go"
// wrap NodeDriver with golang
//go:generate mkdir -p ./contract/driver100
//go:generate go run github.com/ethereum/go-ethereum/cmd/abigen --bin=./contract/solc/NodeDriver.bin --abi=./contract/solc/NodeDriver.abi --pkg=driver100 --type=Contract --out=contract/driver100/contract.go
//go:generate bash -c "(echo -ne '\nvar ContractBinRuntime = \"'; cat contract/solc/NodeDriver.bin-runtime; echo '\"') >> contract/driver100/contract.go"
// wrap NodeDriverAuth with golang
//go:generate mkdir -p ./contract/driverauth100
//go:generate go run github.com/ethereum/go-ethereum/cmd/abigen --bin=./contract/solc/NodeDriverAuth.bin --abi=./contract/solc/NodeDriverAuth.abi --pkg=driverauth100 --type=Contract --out=contract/driverauth100/contract.go
//go:generate bash -c "(echo -ne '\nvar ContractBinRuntime = \"'; cat contract/solc/NodeDriverAuth.bin-runtime; echo '\"') >> contract/driverauth100/contract.go"

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/skyhighblockchain/skyhigh/gossip/contract/driver100"
	"github.com/skyhighblockchain/skyhigh/gossip/contract/driverauth100"
	"github.com/skyhighblockchain/skyhigh/gossip/contract/netinit100"
	"github.com/skyhighblockchain/skyhigh/gossip/contract/sfc100"
	"github.com/skyhighblockchain/skyhigh/logger"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesis/driver"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesis/driverauth"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesis/evmwriter"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesis/netinit"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesis/sfc"
	"github.com/skyhighblockchain/skyhigh/utils"
)

func TestSFC(t *testing.T) {
	logger.SetTestMode(t)
	logger.SetLevel("debug")

	env := newTestEnv()
	defer env.Close()

	var (
		sfc10 *sfc100.Contract
		err   error
	)

	authDriver10, err := driverauth100.NewContract(driverauth.ContractAddress, env)
	require.NoError(t, err)
	rootDriver10, err := driver100.NewContract(driver.ContractAddress, env)
	require.NoError(t, err)

	admin := 1
	adminAddr := env.Address(admin)

	_ = true &&

		t.Run("Genesis SFC", func(t *testing.T) {
			require := require.New(t)

			exp := sfc.GetContractBin()
			got, err := env.CodeAt(nil, sfc.ContractAddress, nil)
			require.NoError(err)
			require.Equal(exp, got, "genesis SFC contract")
			require.Equal(exp, hexutil.MustDecode(sfc100.ContractBinRuntime), "genesis SFC contract version")
		}) &&

		t.Run("Genesis Driver", func(t *testing.T) {
			require := require.New(t)

			exp := driver.GetContractBin()
			got, err := env.CodeAt(nil, driver.ContractAddress, nil)
			require.NoError(err)
			require.Equal(exp, got, "genesis Driver contract")
			require.Equal(exp, hexutil.MustDecode(driver100.ContractBinRuntime), "genesis Driver contract version")
		}) &&

		t.Run("Genesis DriverAuth", func(t *testing.T) {
			require := require.New(t)

			exp := driverauth.GetContractBin()
			got, err := env.CodeAt(nil, driverauth.ContractAddress, nil)
			require.NoError(err)
			require.Equal(exp, got, "genesis DriverAuth contract")
			require.Equal(exp, hexutil.MustDecode(driverauth100.ContractBinRuntime), "genesis DriverAuth contract version")
		}) &&

		t.Run("Network initializer", func(t *testing.T) {
			require := require.New(t)

			exp := netinit.GetContractBin()
			got, err := env.CodeAt(nil, netinit.ContractAddress, nil)
			require.NoError(err)
			require.NotEmpty(exp, "genesis NetworkInitializer contract")
			require.Empty(got, "genesis NetworkInitializer should be destructed")
			require.Equal(exp, hexutil.MustDecode(netinit100.ContractBinRuntime), "genesis NetworkInitializer contract version")
		}) &&

		t.Run("Builtin EvmWriter", func(t *testing.T) {
			require := require.New(t)

			exp := []byte{0}
			got, err := env.CodeAt(nil, evmwriter.ContractAddress, nil)
			require.NoError(err)
			require.Equal(exp, got, "builtin EvmWriter contract")
		}) &&

		t.Run("Some transfers I", func(t *testing.T) {
			cicleTransfers(t, env, 1)
		}) &&

		t.Run("Check owners", func(t *testing.T) {
			require := require.New(t)

			opts := env.ReadOnly()
			opts.From = adminAddr

			isOwn, err := authDriver10.IsOwner(opts)
			require.NoError(err)
			require.True(isOwn)
		}) &&

		t.Run("SFC upgrade", func(t *testing.T) {
			require := require.New(t)

			// create new
			rr := env.ApplyBlock(nextEpoch,
				env.Contract(admin, utils.ToSkh(0), sfc100.ContractBin),
			)
			require.Equal(1, rr.Len())
			require.Equal(types.ReceiptStatusSuccessful, rr[0].Status)
			newImpl := rr[0].ContractAddress
			require.NotEqual(sfc.ContractAddress, newImpl)
			newSfcContractBinRuntime, err := env.CodeAt(nil, newImpl, nil)
			require.NoError(err)
			require.Equal(hexutil.MustDecode(sfc100.ContractBinRuntime), newSfcContractBinRuntime)

			tx, err := authDriver10.CopyCode(env.Payer(admin), sfc.ContractAddress, newImpl)
			require.NoError(err)
			env.incNonce(adminAddr)
			rr = env.ApplyBlock(sameEpoch, tx)
			require.Equal(1, rr.Len())
			require.Equal(types.ReceiptStatusSuccessful, rr[0].Status)
			got, err := env.CodeAt(nil, sfc.ContractAddress, nil)
			require.NoError(err)
			require.Equal(newSfcContractBinRuntime, got, "new SFC contract")

			sfc10, err = sfc100.NewContract(sfc.ContractAddress, env)
			require.NoError(err)
			epoch, err := sfc10.ContractCaller.CurrentEpoch(env.ReadOnly())
			require.Equal(0, epoch.Cmp(big.NewInt(3)), "current epoch %s", epoch.String())
		})

	t.Run("Direct driver", func(t *testing.T) {
		require := require.New(t)

		// create new
		anyContractBin := driver100.ContractBin
		rr := env.ApplyBlock(nextEpoch,
			env.Contract(admin, utils.ToSkh(0), anyContractBin),
		)
		require.Equal(1, rr.Len())
		require.Equal(types.ReceiptStatusSuccessful, rr[0].Status)
		newImpl := rr[0].ContractAddress

		tx, err := rootDriver10.CopyCode(env.Payer(admin), sfc.ContractAddress, newImpl)
		require.NoError(err)
		env.incNonce(adminAddr)
		rr = env.ApplyBlock(sameEpoch, tx)
		require.Equal(1, rr.Len())
		require.NotEqual(types.ReceiptStatusSuccessful, rr[0].Status)
	})

}

func cicleTransfers(t *testing.T, env *testEnv, count uint64) {
	require := require.New(t)
	accounts := len(env.validators)

	// save start balances
	balances := make([]*big.Int, accounts)
	for i := range balances {
		balances[i] = env.State().GetBalance(env.Address(i + 1))
	}

	for i := uint64(0); i < count; i++ {
		// transfers
		txs := make([]*types.Transaction, accounts)
		for i := range txs {
			from := (i)%accounts + 1
			to := (i+1)%accounts + 1
			txs[i] = env.Transfer(from, to, utils.ToSkh(100))
		}

		rr := env.ApplyBlock(sameEpoch, txs...)
		for i, r := range rr {
			fee := big.NewInt(0).Mul(new(big.Int).SetUint64(r.GasUsed), txs[i].GasPrice())
			balances[i] = big.NewInt(0).Sub(balances[i], fee)
		}
	}

	// check balances
	for i := range balances {
		require.Equal(
			balances[i],
			env.State().GetBalance(env.Address(i+1)),
			fmt.Sprintf("account%d", i),
		)
	}
}
