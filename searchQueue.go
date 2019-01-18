package SimpleRTree

import (
	"unsafe"
	"math"
	"encoding/binary"
)

type searchQueueItem struct {
	distance float64 // It is important that distance goes first to avoid offset computation
	node   uintptr   // if nil item carries node
	px, py float64 // points are not stored in nodes so we need to track them explicitely
}

var rSEARCH_QUEUE_ITEM_SIZE = unsafe.Sizeof(searchQueueItem{})
const rSEARCH_QUEUE_ITEM_SIZE_INT = 32 // int(unsafe.Sizeof(searchQueueItem{}))
var rSEARCH_QUEUE_ITEM_DISTANCE_OFFSET = int(unsafe.Offsetof(searchQueueItem{}.distance)) // This is 0
var rSEARCH_QUEUE_ITEM_NODE_OFFSET = int(unsafe.Offsetof(searchQueueItem{}.node))


type searchQueue []byte

func (sq searchQueue) Len() int {
	return len(sq)
}

func (sq searchQueue) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
	sq[i+ 1], sq[j + 1] = sq[j + 1], sq[i +1]
	sq[i+ 2], sq[j + 2] = sq[j + 2], sq[i +2]
	sq[i+ 3], sq[j + 3] = sq[j + 3], sq[i +3]
	sq[i+ 4], sq[j + 4] = sq[j + 4], sq[i +4]
	sq[i+ 5], sq[j + 5] = sq[j + 5], sq[i +5]
	sq[i+ 6], sq[j + 6] = sq[j + 6], sq[i +6]
	sq[i+ 7], sq[j + 7] = sq[j + 7], sq[i +7]
	sq[i+ 8], sq[j + 8] = sq[j + 8], sq[i +8]
	sq[i+ 9], sq[j + 9] = sq[j + 9], sq[i +9]
	sq[i+ 10], sq[j + 10] = sq[j + 10], sq[i +10]
	sq[i+ 11], sq[j + 11] = sq[j + 11], sq[i +11]
	sq[i+ 12], sq[j + 12] = sq[j + 12], sq[i +12]
	sq[i+ 13], sq[j + 13] = sq[j + 13], sq[i +13]
	sq[i+ 14], sq[j + 14] = sq[j + 14], sq[i +14]
	sq[i+ 15], sq[j + 15] = sq[j + 15], sq[i +15]
	sq[i+ 16], sq[j + 16] = sq[j + 16], sq[i +16]
	sq[i+ 17], sq[j + 17] = sq[j + 17], sq[i +17]
	sq[i+ 18], sq[j + 18] = sq[j + 18], sq[i +18]
	sq[i+ 19], sq[j + 19] = sq[j + 19], sq[i +19]
	sq[i+ 20], sq[j + 20] = sq[j + 20], sq[i +20]
	sq[i+ 21], sq[j + 21] = sq[j + 21], sq[i +21]
	sq[i+ 22], sq[j + 22] = sq[j + 22], sq[i +22]
	sq[i+ 23], sq[j + 23] = sq[j + 23], sq[i +23]
	sq[i+ 24], sq[j + 24] = sq[j + 24], sq[i +24]
	sq[i+ 25], sq[j + 25] = sq[j + 25], sq[i +25]
	sq[i+ 26], sq[j + 26] = sq[j + 26], sq[i +26]
	sq[i+ 27], sq[j + 27] = sq[j + 27], sq[i +27]
	sq[i+ 28], sq[j + 28] = sq[j + 28], sq[i +28]
	sq[i+ 29], sq[j + 29] = sq[j + 29], sq[i +29]
	sq[i+ 30], sq[j + 30] = sq[j + 30], sq[i +30]
	sq[i+ 31], sq[j + 31] = sq[j + 31], sq[i +31]
}

func (sq searchQueue) PreparePop () {
	n := sq.Len() - rSEARCH_QUEUE_ITEM_SIZE_INT
	for j := 0; j < n; j+= rSEARCH_QUEUE_ITEM_SIZE_INT {
		if sq.Less(j, n) {
			sq.Swap(n, j)
		}
	}
}

func (sq searchQueue) Less(i, j int) bool {
	// We want to pop element with smaller distance
	return math.Float64frombits(binary.LittleEndian.Uint64(sq[i:])) < math.Float64frombits(binary.LittleEndian.Uint64(sq[j:]))
}

func (sq searchQueue) setItem (i int, item searchQueueItem) {
	binary.LittleEndian.PutUint64(sq[i:], math.Float64bits(item.distance))
	binary.LittleEndian.PutUint64(sq[i + 8:], uint64(item.node))
	binary.LittleEndian.PutUint64(sq[i + 16:], math.Float64bits(item.px))
	binary.LittleEndian.PutUint64(sq[i + 24:], math.Float64bits(item.py))
}
func (sq searchQueue) itemAt (i int) searchQueueItem {
	var item searchQueueItem
	item.distance = math.Float64frombits(binary.LittleEndian.Uint64(sq[i:]))
	item.node = uintptr(binary.LittleEndian.Uint64(sq[i + 8:]))
	item.px = math.Float64frombits(binary.LittleEndian.Uint64(sq[i+ 16:]))
	item.py = math.Float64frombits(binary.LittleEndian.Uint64(sq[i+ 24:]))
	return item
}

