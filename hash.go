package ovoclient

import(
)

func GetPositiveHashCode(key string, maxNum int32) int32 {
	var hash1 int32 = 5381
	var hash2 int32 = 5381
	var c int32 = 0
	var index int32 = 0
	for _,r := range key {
		c = int32(r)
		if index % 2 == 0 {
	        hash1 = ((hash1 << 5) + hash1) ^ c
	    } else {
	        hash2 = ((hash2 << 5) + hash2) ^ c
	    }
		index++
	}
	hash1 = hash1 + (hash2 * 1566083941)
	if hash1 < 0 {
		hash1 = (-1) * hash1
	}
	return hash1 % maxNum;
}