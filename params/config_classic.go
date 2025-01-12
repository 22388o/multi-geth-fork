// Copyright 2019 The multi-geth Authors
// This file is part of the multi-geth library.
//
// The multi-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The multi-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the multi-geth library. If not, see <http://www.gnu.org/licenses/>.
package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/types/multigeth"
)

var (
	// ClassicChainConfig is the chain parameters to run a node on the Classic main network.
	ClassicChainConfig = &multigeth.MultiGethChainConfig{
		NetworkID: 1,
		Ethash:    new(ctypes.EthashConfig),
		ChainID:   big.NewInt(61),

		EIP2FBlock: big.NewInt(1150000),
		EIP7FBlock: big.NewInt(1150000),

		//DAOForkBlock:        big.NewInt(1920000),

		EIP150Block: big.NewInt(2500000),

		EIP155Block:  big.NewInt(3000000),
		EIP160FBlock: big.NewInt(3000000),

		// EIP158~
		EIP161FBlock: big.NewInt(8772000),
		EIP170FBlock: big.NewInt(8772000),

		// Byzantium eq
		EIP100FBlock: big.NewInt(8772000),
		EIP140FBlock: big.NewInt(8772000),
		EIP198FBlock: big.NewInt(8772000),
		EIP211FBlock: big.NewInt(8772000),
		EIP212FBlock: big.NewInt(8772000),
		EIP213FBlock: big.NewInt(8772000),
		EIP214FBlock: big.NewInt(8772000),
		EIP658FBlock: big.NewInt(8772000),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(9573000),
		EIP1014FBlock: big.NewInt(9573000),
		EIP1052FBlock: big.NewInt(9573000),

		// Istanbul eq, aka Aztlan
		// ECIP-1061
		EIP152FBlock:  big.NewInt(10_500_839),
		EIP1108FBlock: big.NewInt(10_500_839),
		EIP1344FBlock: big.NewInt(10_500_839),
		EIP1884FBlock: nil,
		EIP2028FBlock: big.NewInt(10_500_839),
		EIP2200FBlock: big.NewInt(10_500_839), // RePetersburg (=~ re-1283)

		// ECIP-1078, aka Phoenix Fix
		EIP2200DisableFBlock: big.NewInt(10_500_839),
		EIP1283FBlock:        big.NewInt(10_500_839),
		EIP1706FBlock:        big.NewInt(10_500_839),
		ECIP1080FBlock:       big.NewInt(10_500_839),

		DisposalBlock:      big.NewInt(5900000),
		ECIP1017FBlock:     big.NewInt(5000000),
		ECIP1017EraRounds:  big.NewInt(5000000),
		ECIP1010PauseBlock: big.NewInt(3000000),
		ECIP1010Length:     big.NewInt(2000000),
		RequireBlockHashes: map[uint64]common.Hash{
			1920000: common.HexToHash("0x94365e3a8c0b35089c1d1195081fe7489b528a84b22199c916180db8b28ade7f"),
			2500000: common.HexToHash("0xca12c63534f565899681965528d536c52cb05b7c48e269c2a6cb77ad864d878a"),
		},
	}

	DisinflationRateQuotient = big.NewInt(4)      // Disinflation rate quotient for ECIP1017
	DisinflationRateDivisor  = big.NewInt(5)      // Disinflation rate divisor for ECIP1017
	ExpDiffPeriod            = big.NewInt(100000) // Exponential diff period for diff bomb & ECIP1010
)
