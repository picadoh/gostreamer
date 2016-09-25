package streamer

import (
	"hash/fnv"
	"fmt"
)

/**
The hash function generates a hash for a specific string representation of a value
and modules it to a limit.
 */
func Hash(value interface{}, module int) uint32 {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%s", value)))
	return h.Sum32() % uint32(module)
}
