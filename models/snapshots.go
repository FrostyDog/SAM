package models

import (
	"log"

	"github.com/Kucoin/kucoin-go-sdk"
)

type SnapshotsContainer map[int]kucoin.TickersModel

func NewSnapshotsContainter() SnapshotsContainer {
	return SnapshotsContainer{}
}

func (cont SnapshotsContainer) AddSnapshots(list kucoin.TickersModel, number int) {
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
