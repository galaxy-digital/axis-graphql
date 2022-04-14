/*
Package rpc implements bridge to Lachesis full node API interface.

We recommend using local IPC for fast and the most efficient inter-process communication between the API server
and an Opera/Lachesis node. Any remote RPC connection will work, but the performance may be significantly degraded
by extra networking overhead of remote RPC calls.

You should also consider security implications of opening Lachesis RPC interface for remote access.
If you considering it as your deployment strategy, you should establish encrypted channel between the API server
and Lachesis RPC interface with connection limited to specified endpoints.

We strongly discourage opening Lachesis RPC interface for unrestricted Internet access.
*/
package rpc

import (
	"axis-graphql/internal/types"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// BlockTypeLatest represents the latest available block in blockchain.
const (
	BlockTypeLatest   = "latest"
	BlockTypeEarliest = "earliest"
)

// MustBlockHeight returns the current block height
// of the blockchain. It returns nil if the block height can not be pulled.
func (axis *AxisBridge) MustBlockHeight() *big.Int {
	var val hexutil.Big
	// axis_blockNumber
	if err := axis.rpc.Call(&val, "axis_blockNumber"); err != nil {
		axis.log.Errorf("failed block height check; %s", err.Error())
		return nil
	}
	return val.ToInt()
}

// BlockHeight returns the current block height of the AXIS blockchain.
func (axis *AxisBridge) BlockHeight() (*hexutil.Big, error) {
	// keep track of the operation
	axis.log.Debugf("checking current block height")

	// call for data
	var height hexutil.Big

	err := axis.rpc.Call(&height, "axis_blockNumber") // axis_blockNumber
	if err != nil {
		axis.log.Error("block height could not be obtained")
		return nil, err
	}

	// inform and return
	axis.log.Debugf("current block height is %s", height.String())
	return &height, nil
}

// Block returns information about a blockchain block by encoded hex number, or by a type tag.
// For tag based loading use predefined BlockType contacts.
func (axis *AxisBridge) Block(numTag *string) (*types.Block, error) {
	// keep track of the operation
	axis.log.Debugf("loading details of block num/tag %s", *numTag)

	// call for data
	var block types.Block
	err := axis.rpc.Call(&block, "axis_getBlockByNumber", numTag, false)
	if err != nil {
		axis.log.Error("block could not be extracted")
		return nil, err
	}

	// detect block not found situation; block number is zero and the hash is also zero
	/* if uint64(block.Number) == 0 && block.Hash.Big().Cmp(big.NewInt(0)) == 0 {
		axis.log.Debugf("block [%s] not found", *numTag)
		return nil, fmt.Errorf("block not found")
	} */

	// keep track of the operation
	axis.log.Debugf("block #%d found at mark %s",
		uint64(block.Number), time.Unix(int64(block.TimeStamp), 0).String())
	return &block, nil
}

// BlockByHash returns information about a blockchain block by hash.
func (axis *AxisBridge) BlockByHash(hash *string) (*types.Block, error) {
	// keep track of the operation
	axis.log.Debugf("loading details of block %s", *hash)

	// call for data
	var block types.Block
	err := axis.rpc.Call(&block, "axis_getBlockByHash", hash, false)
	if err != nil {
		axis.log.Error("block could not be extracted")
		return nil, err
	}

	// detect block not found situation
	if uint64(block.Number) == 0 {
		axis.log.Debugf("block [%s] not found", *hash)
		return nil, fmt.Errorf("block not found")
	}

	// inform and return
	axis.log.Debugf("block #%d found at mark %s by hash %s",
		uint64(block.Number), time.Unix(int64(block.TimeStamp), 0).String(), *hash)
	return &block, nil
}
