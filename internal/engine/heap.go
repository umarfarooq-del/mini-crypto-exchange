package engine

import "mini-crypto-exchange/internal/models"

// BuyHeap is a max-heap for buy orders (higher price first, earlier time first)
type BuyHeap []*models.Order

func (h BuyHeap) Len() int {
	return len(h)
}

func (h BuyHeap) Less(i, j int) bool {
	// Max-heap: higher price comes first
	if h[i].Price != h[j].Price {
		return h[i].Price > h[j].Price
	}
	// If prices equal, earlier time comes first (FIFO)
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h BuyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *BuyHeap) Push(x interface{}) {
	*h = append(*h, x.(*models.Order))
}

func (h *BuyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// SellHeap is a min-heap for sell orders (lower price first, earlier time first)
type SellHeap []*models.Order

func (h SellHeap) Len() int {
	return len(h)
}

func (h SellHeap) Less(i, j int) bool {
	// Min-heap: lower price comes first
	if h[i].Price != h[j].Price {
		return h[i].Price < h[j].Price
	}
	// If prices equal, earlier time comes first (FIFO)
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h SellHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *SellHeap) Push(x interface{}) {
	*h = append(*h, x.(*models.Order))
}

func (h *SellHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
