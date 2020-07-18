package apollo

import (
	"github.com/micro/go-micro/v2/util/log"
	"strings"
	"github.com/micro/go-micro/v2/config/encoder"
)

func makeMap(e encoder.Encoder, kv map[string]string) (map[string]interface{}, error) {	
	data := make(map[string]interface{})

	// consul guarantees lexicographic order, so no need to sort
	for k, v := range kv {
		pathString := k
		if pathString == "" {
			continue
		}
		
		var err error

		// set target at the root
		target := data
		path := strings.Split(k, ".")
		// find (or create) the leaf node we want to put this value at
		for _, dir := range path[:len(path)-1] {
			if _, ok := target[dir]; !ok {
				target[dir] = make(map[string]interface{})
			}
			target = target[dir].(map[string]interface{})
		}

		leafDir := path[len(path)-1]

		var vv interface{}
		err = e.Decode([]byte(v), &vv)
		if err != nil {
			log.Error(err)
		}

		target[leafDir] = vv
	}

	return data, nil
}