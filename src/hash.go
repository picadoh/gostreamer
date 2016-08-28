package main

import (
	"hash/fnv"
	"fmt"
)

func hash(value interface{}, module int) uint32 {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%s", value)))
	return h.Sum32() % uint32(module)
}
