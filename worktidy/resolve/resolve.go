package resolve

import (
	"fmt"

	"github.com/ymd-h/go/sets"
)

type (
	DepsSets = map[string]sets.ISet[string]
)


func Resolve(unresolved DepsSets) (DepsSets, error) {
	resolved := DepsSets{}

	for mod, deps := range unresolved {
		stack := deps.ToSlice()

		for len(stack) > 0 {
			d := stack[0]
			stack = stack[1:]

			if u, ok := unresolved[d]; ok {
				for _, ud := range u.ToSlice() {
					if !deps.Has(ud) {
						deps.Add(ud)
						stack = append(stack, ud)
					}
				}
			} else {
				return nil, fmt.Errorf("Unknown Dependency: %s", d)
			}
		}

		resolved[mod] = deps
	}

	return resolved, nil
}
