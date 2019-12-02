package tests

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ethereum/go-ethereum/params/convert"
	"github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/parity"
)

var paritySpecsDir = filepath.Join("..", "params", "parity.json.d")

func paritySpecPath(name string) string {
	p := filepath.Join(paritySpecsDir, name)
	if fi, err := os.Open(p); err == nil {
		fi.Close()
		return p
	} else if os.IsNotExist(err) {
		// This is an ugly HACK because tests function are sometimes called from
		// other packages that are nested more deeply, eg. eth/tracers.
		// This is a workaround for that.
		// And it sucks.
		p = filepath.Join("..", paritySpecsDir, name)
	}
	return p
}

var MapForkNameChainspecFileState = map[string]string{
	"Frontier":             "frontier_test.json",
	"Homestead":            "homestead_test.json",
	"EIP150":               "eip150_test.json",
	"EIP158":               "eip161_test.json",
	"Byzantium":            "byzantium_test.json",
	"Constantinople":       "constantinople_test.json",
	"ConstantinopleFix":    "st_peters_test.json",
	"EIP158ToByzantiumAt5": "transition_test.json",
	"Istanbul":             "istanbul_test.json",
	"ETC_Atlantis":         "classic_atlantis_test.json",
	"ETC_Agharta":          "classic_agharta_test.json",
}

var mapForkNameChainspecFileDifficulty = map[string]string{
	"Ropsten":           "ropsten_difficulty_test.json",
	"Morden":            "morden_difficulty_test.json",
	"Frontier":          "frontier_difficulty_test.json",
	"Homestead":         "homestead_difficulty_test.json",
	"Byzantium":         "byzantium_difficulty_test.json",
	"MainNetwork":       "mainnetwork_difficulty_test.json",
	"CustomMainNetwork": "custom_mainnetwork_difficulty_test.json",
	"Constantinople":    "constantinople_difficulty_test.json",
	"difficulty.json":   "difficulty_json_difficulty_test.json",
	"ETC_Atlantis":      "classic_atlantis_difficulty_test.json",
	"ETC_Agharta":       "classic_agharta_difficulty_test.json",
}

func readConfigFromSpecFile(name string) (spec common.ChainConfigurator, sha1sum []byte, err error) {
	spec = &parity.ParityChainSpec{}
	if fi, err := os.Open(name); os.IsNotExist(err) {
		return nil, nil, err
	} else {
		fi.Close()
	}
	b, err := ioutil.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("%s err: %s\n%s", name, err, b))
	}
	err = json.Unmarshal(b, spec)
	if err != nil {
		if jsonError, ok := err.(*json.SyntaxError); ok {
			line, character, lcErr := lineAndCharacter(string(b), int(jsonError.Offset))
			fmt.Fprintf(os.Stderr, "test failed with error: Cannot parse JSON schema due to a syntax error at line %d, character %d: %v\n", line, character, jsonError.Error())
			if lcErr != nil {
				fmt.Fprintf(os.Stderr, "Couldn't find the line and character position of the error due to error %v\n", lcErr)
			}
		}
		if jsonError, ok := err.(*json.UnmarshalTypeError); ok {
			line, character, lcErr := lineAndCharacter(string(b), int(jsonError.Offset))
			fmt.Fprintf(os.Stderr, "test failed with error: The JSON type '%v' cannot be converted into the Go '%v' type on struct '%s', field '%v'. See input file line %d, character %d\n", jsonError.Value, jsonError.Type.Name(), jsonError.Struct, jsonError.Field, line, character)
			if lcErr != nil {
				fmt.Fprintf(os.Stderr, "test failed with error: Couldn't find the line and character position of the error due to error %v\n", lcErr)
			}
		}
		panic(fmt.Sprintf("%s err: %s\n%s", name, err, b))
	}
	bb := sha1.Sum(b)
	return spec, bb[:], nil
}

func init() {

	if os.Getenv(MG_CHAINCONFIG_FEATURE_EQ_KEY) != "" {
		log.Println("Setting equivalent fork feature chain configurations")

		for i, config := range Forks {
			mgc := &paramtypes.MultiGethChainConfig{}
			if err := convert.Convert(config, mgc); common.IsFatalUnsupportedErr(err) {
				panic(err)
			}
			Forks[i] = mgc
		}

		for k, v := range difficultyChainConfigurations {
			mgc := &paramtypes.MultiGethChainConfig{}
			if err := convert.Convert(v, mgc); common.IsFatalUnsupportedErr(err) {
				panic(err)
			}
			difficultyChainConfigurations[k] = mgc
		}

	} else if os.Getenv(MG_CHAINCONFIG_CHAINSPEC_KEY) != "" {
		log.Println("Setting chain configurations from Parity chainspecs")

		for k, v := range MapForkNameChainspecFileState {
			config, sha1sum, err := readConfigFromSpecFile(paritySpecPath(v))
			if os.IsNotExist(err) {
				wd, wde := os.Getwd()
				if wde != nil {
					panic(wde)
				}
				panic(fmt.Sprintf("failed to find chainspec, wd: %s", wd))
			} else if err != nil {
				panic(err)
			}
			chainspecRefsState[k] = chainspecRef{filepath.Base(v), sha1sum}
			if diffs := convert.Equal(reflect.TypeOf((*common.ChainConfigurator)(nil)), Forks[k], config); len(diffs) != 0 {
				log.Println(k, v, len(diffs), "diffs")
				for _, diff := range diffs {
					log.Println(diff)
				}
				//panic("not same configs")
			}
			if err := common.Equivalent(Forks[k], config); err != nil {
				log.Println("Not equivalent configs", "err", err)
				panic("")
			}
			Forks[k] = config
		}

		for k, v := range mapForkNameChainspecFileDifficulty {
			config, sha1sum, err := readConfigFromSpecFile(paritySpecPath(v))
			if os.IsNotExist(err) && os.Getenv(MG_GENERATE_DIFFICULTY_TESTS_KEY) != "" {
				log.Println("Will generate chainspec file for", k, v)
			} else if len(sha1sum) == 0 {
				panic("zero sum game")
			} else {
				chainspecRefsDifficulty[k] = chainspecRef{filepath.Base(v), sha1sum}
				difficultyChainConfigurations[k] = config
			}
		}
	}
}

//func convertMetaForkBlocksToFeatures(config *paramtypes.MultiGethChainConfig) {
//	if config.HomesteadBlock != nil {
//		config.EIP2FBlock = config.HomesteadBlock
//		config.EIP7FBlock = config.HomesteadBlock
//		config.HomesteadBlock = nil
//	}
//	if config.EIP158Block != nil {
//		config.EIP160FBlock = config.EIP158Block
//		config.EIP161FBlock = config.EIP158Block
//		config.EIP170FBlock = config.EIP158Block
//		config.EIP158Block = nil
//	}
//	if config.ByzantiumBlock != nil {
//		// Difficulty adjustment to target mean block time including uncles
//		// https://github.com/ethereum/EIPs/issues/100
//		config.EIP100FBlock = config.ByzantiumBlock
//		// Opcode REVERT
//		// https://eips.ethereum.org/EIPS/eip-140
//		config.EIP140FBlock = config.ByzantiumBlock
//		// Precompiled contract for bigint_modexp
//		// https://github.com/ethereum/EIPs/issues/198
//		config.EIP198FBlock = config.ByzantiumBlock
//		// Opcodes RETURNDATACOPY, RETURNDATASIZE
//		// https://github.com/ethereum/EIPs/issues/211
//		config.EIP211FBlock = config.ByzantiumBlock
//		// Precompiled contract for pairing check
//		// https://github.com/ethereum/EIPs/issues/212
//		config.EIP212FBlock = config.ByzantiumBlock
//		// Precompiled contracts for addition and scalar multiplication on the elliptic curve alt_bn128
//		// https://github.com/ethereum/EIPs/issues/213
//		config.EIP213FBlock = config.ByzantiumBlock
//		// Opcode STATICCALL
//		// https://github.com/ethereum/EIPs/issues/214
//		config.EIP214FBlock = config.ByzantiumBlock
//		// Metropolis diff bomb delay and reducing block reward
//		// https://github.com/ethereum/EIPs/issues/649
//		// note that this is closely related to EIP100.
//		// In fact, EIP100 is bundled in
//		config.EIP649FBlock = config.ByzantiumBlock
//		// Transaction receipt status
//		// https://github.com/ethereum/EIPs/issues/658
//		config.EIP658FBlock = config.ByzantiumBlock
//		// NOT CONFIGURABLE: prevent overwriting contracts
//		// https://github.com/ethereum/EIPs/issues/684
//		// EIP684FBlock *big.Int `json:"eip684BFlock,omitempty"`
//
//		config.ByzantiumBlock = nil
//	}
//	if config.ConstantinopleBlock != nil {
//		// Opcodes SHR, SHL, SAR
//		// https://eips.ethereum.org/EIPS/eip-145
//		config.EIP145FBlock = config.ConstantinopleBlock
//		// Opcode CREATE2
//		// https://eips.ethereum.org/EIPS/eip-1014
//		config.EIP1014FBlock = config.ConstantinopleBlock
//		// Opcode EXTCODEHASH
//		// https://eips.ethereum.org/EIPS/eip-1052
//		config.EIP1052FBlock = config.ConstantinopleBlock
//		// Constantinople difficulty bomb delay and block reward adjustment
//		// https://eips.ethereum.org/EIPS/eip-1234
//		config.EIP1234FBlock = config.ConstantinopleBlock
//		// Net gas metering
//		// https://eips.ethereum.org/EIPS/eip-1283
//		config.EIP1283FBlock = config.ConstantinopleBlock
//
//		config.ConstantinopleBlock = nil
//	}
//	if config.IstanbulBlock != nil {
//		config.EIP152FBlock = config.IstanbulBlock
//		config.EIP1108FBlock = config.IstanbulBlock
//		config.EIP1344FBlock = config.IstanbulBlock
//		config.EIP1884FBlock = config.IstanbulBlock
//		config.EIP2028FBlock = config.IstanbulBlock
//		config.EIP2200FBlock = config.IstanbulBlock
//		config.IstanbulBlock = nil
//	}
//}

// https://adrianhesketh.com/2017/03/18/getting-line-and-character-positions-from-gos-json-unmarshal-errors/
func lineAndCharacter(input string, offset int) (line int, character int, err error) {
	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return 0, 0, fmt.Errorf("Couldn't find offset %d within the input.", offset)
	}

	// Humans tend to count from 1.
	line = 1
	for i, b := range input {
		if b == lf {
			line++
			character = 0
		}
		character++
		if i == offset {
			break
		}
	}
	return line, character, nil
}
