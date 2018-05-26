package main

import (
	"fmt"
	"github.com/golang-collections/go-datastructures/bitarray"
	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	bits bitarray.BitArray
}

func NewBloomFilter() *BloomFilter {
	// a sparse bit array has worse perf but saves us having to worry about length
	ba := bitarray.NewSparseBitArray()
	return &BloomFilter{ba}
}

func (f *BloomFilter) getPositions(data []byte) []uint64 {
	// we generate two positions for each value
	// the more positions you generate, the slower the filter is, but the less chance of false positives
	p1 := murmur3.Sum64WithSeed(data, uint32(1))
	p2 := murmur3.Sum64WithSeed(data, uint32(2))
	return []uint64{p1, p2}
}

func (f *BloomFilter) Set(data []byte) *BloomFilter {
	for _, pos := range f.getPositions(data) {
		f.bits.SetBit(pos)
	}
	return f
}

func (f *BloomFilter) Test(data []byte) bool {
	for _, pos := range f.getPositions(data) {
		hasBit, _ := f.bits.GetBit(pos)
		if !hasBit {
			return false
		}
	}
	return true
}

type ArrayWithBloomFilter struct {
	array  []string
	filter *BloomFilter
}

func NewArrayWithBloomFilter() *ArrayWithBloomFilter {
	arr := make([]string, 0)
	bf := NewBloomFilter()
	return &ArrayWithBloomFilter{arr, bf}
}

func (a *ArrayWithBloomFilter) Set(value string) {
	a.filter.Set([]byte(value))
	a.array = append(a.array, value)
}

func (a *ArrayWithBloomFilter) Test(value string) bool {
	hasElement := a.filter.Test([]byte(value))
	if !hasElement {
		// we know the array doesn't have the element
		return false
	} else {
		// check manually - this will be slow
		// we need to check because a bloom filter only guarantees no false
		// negatives. false positives are still possible (especially with only two hashing functions)
		for _, el := range a.array {
			if el == value {
				return true
			}
		}
		return false
	}
}

func main() {

	bf := NewBloomFilter()

	data := []byte("test")
	fmt.Println("Should be true:")
	fmt.Println(bf.Set(data).Test(data))
	fmt.Println("Should be false:")
	fmt.Println(bf.Test([]byte("test2")))

    arr := NewArrayWithBloomFilter()
    arr.Set("test")
	fmt.Println("Should be true:")
    fmt.Println(arr.Test("test"))
	fmt.Println("Should be false:")
    fmt.Println(arr.Test("test2"))

}
