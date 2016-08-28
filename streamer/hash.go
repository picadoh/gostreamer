package streamer

import (
	"hash/fnv"
	"fmt"
)

func Hash(value interface{}, module int) uint32 {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%s", value)))
	return h.Sum32() % uint32(module)
}
