# Bloom

A toy implementation of a bloom filter in Go. It's toy because it uses a sparse bitset, which makes the code much simpler at the cost of speed, and because it uses a fixed number of hashing functions. This is for people who want to learn what a bloom filter is, not for people who want to put a bloom filter into production.

It includes an implementation of an array-like data structure that uses a bloom filter to quickly check if an element is in the array. The idea is to give a short, useful example of the kind of thing a bloom filter is actually _for_.

The code is ideally commented aggressively enough that someone who knows Go but has never heard of a bloom filter can understand it.
