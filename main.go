package main

import (
	"fmt"
	"strconv"
)

// A bloom filter is an array of bits, a function for adding elements, and a function for testing if an element has probably been added
type BloomFilter struct {
	bits [99]bool // Every bloom filter begins with every bit set to 0: [0,0,0,0,0...]
}

// We need a function that takes an element and returns two positions between 0 and 99
// This function must be deterministic: every time you run it with the same data, you have to get the same positions
// In real life you'd use a proper hashing function, but here we just hack up our own
func (f *BloomFilter) getPositions(data []byte) []int {
	p1 := 0
	p2 := 0
	for _, b := range data {
		p1 += int(b >> 1)
		p2 += int(b >> 2)
	}
	p1, _ = strconv.Atoi(strconv.Itoa(p1)[:2])
	p2, _ = strconv.Atoi(strconv.Itoa(p2)[:2])
	return []int{p1, p2}
}

// Adding an element to a bloom filter means setting a fixed number of bits to 1 in the bit array
// Bits may never be set back to 0, under any circumstances
func (f *BloomFilter) Set(data []byte) *BloomFilter {
	for _, pos := range f.getPositions(data) {
		f.bits[pos] = true
	}
	return f
}

// To test if an element has been added to the bloom filter, we generate the bits that would have been
// set to 1 when the element was added. If any of these bits are 0, we know for a fact that the element
// has not been added
// Note that the converse does not apply. If all the bits are 1, the element may still not have been added
// if adding other elements has flipped the same bits
func (f *BloomFilter) Test(data []byte) bool {
	for _, pos := range f.getPositions(data) {
		hasBit := f.bits[pos]
		if !hasBit {
			return false
		}
	}
	return true
}

// That's it! That's a functioning bloom filter in three tiny functions

// ---------------------------------------------------------------------

// Here's an example of a useful data structure that uses a bloom filter. It's an array of strings where every string that
// gets added gets also added to a bloom filter. We can thus check if a string belongs to the array by first
// asking the bloom filter. We only iterate over the array if the bloom filter can't rule the string out. For a
// very large array, this could save a lot of time!
type ArrayWithBloomFilter struct {
	array  []string
	filter *BloomFilter
}

func NewArrayWithBloomFilter() *ArrayWithBloomFilter {
	arr := make([]string, 0)
	bf := BloomFilter{}
	return &ArrayWithBloomFilter{arr, &bf}
}

func (a *ArrayWithBloomFilter) Set(value string) {
	a.filter.Set([]byte(value))      // Add the element to the bloom filter
	a.array = append(a.array, value) // Add the element to the array
}

func (a *ArrayWithBloomFilter) Test(value string) bool {
	hasElement := a.filter.Test([]byte(value))
	if !hasElement {
		// We know the array doesn't have the element, since a bloom filter guarantees
		// no false negatives
		return false
	} else {
		// Since a bloom filter doesn't guarantee no false positives, we need to check manually
		// This will be a slow operation for a large array
		for _, el := range a.array {
			if el == value {
				return true
			}
		}
		return false
	}
}

func main() {
	arr := NewArrayWithBloomFilter()
	arr.Set("test")

	fmt.Println("Should be true:")
	fmt.Println(arr.Test("test")) // This will iterate over the array
	fmt.Println("Should be false:")

	// This will not iterate over the array, unless
	// our two hashing functions happen to return the same two positions
	// as they did for the value "test"
	fmt.Println(arr.Test("test2"))

}
