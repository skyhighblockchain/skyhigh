package skyhigh

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/skyhighblockchain/push-base/hash"
	"github.com/skyhighblockchain/push-base/inter/idx"

	"github.com/skyhighblockchain/skyhigh/inter"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesis"
	"github.com/skyhighblockchain/skyhigh/skyhigh/genesis/gpos"
)

type Genesis struct {
	Accounts    genesis.Accounts
	Storage     genesis.Storage
	Delegations genesis.Delegations
	Blocks      genesis.Blocks
	RawEvmItems genesis.RawEvmItems
	Validators  gpos.Validators

	FirstEpoch    idx.Epoch
	PrevEpochTime inter.Timestamp
	Time          inter.Timestamp
	ExtraData     []byte

	TotalSupply *big.Int

	DriverOwner common.Address

	Rules Rules

	Hash func() hash.Hash
}
