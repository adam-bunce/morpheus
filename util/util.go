package util

import (
	"fmt"
)

func DeepCopyMap[K comparable, V any](m map[K]V) map[K]V {
	mapCopy := make(map[K]V, len(m))

	for k, v := range m {
		mapCopy[k] = v
	}

	return mapCopy
}

func Display(string string) {
	fmt.Println(">", string)
}
