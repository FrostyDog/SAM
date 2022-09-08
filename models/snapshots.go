package models

import (
	"log"

	"github.com/Kucoin/kucoin-go-sdk"
)

type SnapshotsContainer map[int]kucoin.TickersModel

func NewSnapshotsContainter() SnapshotsContainer {
	return SnapshotsContainer{}
}

// Appending snap while removing left-hand value (max snaps = 2)
func (cont SnapshotsContainer) AddSnapshotAndReplace(list kucoin.TickersModel) {
	switch len(cont) {
	case 0:
		cont[0] = list
	case 1:
		cont[1] = list
	case 2:
		cont[0] = cont[1]
		cont[1] = list
	}
}

func (cont SnapshotsContainer) AddSnapshot(list kucoin.TickersModel) {
	cont[len(cont)] = list
}

func (cont SnapshotsContainer) AddSnapshotAtIndex(list kucoin.TickersModel, number int) {
	cont[number] = list
}

func (cont SnapshotsContainer) DeleteSnapshot(number int) {
	if _, ok := cont[number]; ok {
		delete(cont, number)
	} else {
		log.Println("Key doesn't exist in the snapshot")
	}

}

func (cont SnapshotsContainer) ClearSnapshots() {
	for key := range cont {
		delete(cont, key)
	}
}
