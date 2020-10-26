package transform

import (
	"sort"
	"strings"

	"github.com/itsmurugappan/kubernetes-resource-builder/pkg/kubernetes"
)

func GetStringMap(kvs []kubernetes.KV, currentMap map[string]string) map[string]string {
	if len(kvs) > 0 && kvs[0].Key != "" {
		if currentMap == nil {
			currentMap = make(map[string]string)
		}
		for _, kv := range kvs {
			currentMap[kv.Key] = kv.Value
		}
	}
	return currentMap
}

func GetKVfromMap(m map[string]string) []kubernetes.KV {
	if len(m) > 0 {
		var keys []string
		var kvs []kubernetes.KV
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			kvs = append(kvs, kubernetes.KV{k, m[k]})
		}
		return kvs
	}
	return nil
}

func ConstructMapFromOFPStyleEnvString(s string) map[string]string {
	retMap := make(map[string]string)
	for _, c := range strings.Split(s, ",") {
		d := strings.Split(c, "|")
		retMap[d[0]] = d[1]
	}
	return retMap
}
