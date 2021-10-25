package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type StorageForStateObject map[common.Hash]common.Hash

func (s StorageForStateObject) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

func (s StorageForStateObject) Copy() StorageForStateObject {
	cpy := make(StorageForStateObject)
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}
