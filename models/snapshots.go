package models

import (
	"log"

	"github.com/Kucoin/kucoin-go-sdk"
)

type content interface {
	kucoin.TickersModel | TickerPrices
}
type SnapshotsContainer[T content] map[int]T

func NewSnapshotsContainter[T content]() SnapshotsContainer[T] {
	return SnapshotsContainer[T]{}
}

// Appending snap while removing left-hand value (max snaps = 2)
func (cont SnapshotsContainer[T]) AddSnapshotAndReplace(list T) {
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

func (cont SnapshotsContainer[T]) AddSnapshot(list T) {
	cont[len(cont)] = list
}

func (cont SnapshotsContainer[T]) AddSnapshotAtIndex(list T, number int) {
	cont[number] = list
}

func (cont SnapshotsContainer[T]) DeleteSnapshot(number int) {
	if _, ok := cont[number]; ok {
		delete(cont, number)
	} else {
		log.Println("Key doesn't exist in the snapshot")
	}

}

func (cont SnapshotsContainer[T]) ClearSnapshots() {
	for key := range cont {
		delete(cont, key)
	}
}
