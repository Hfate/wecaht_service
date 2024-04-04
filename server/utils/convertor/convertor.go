package convertor

import (
	"github.com/deckarep/golang-set"
	"github.com/spf13/cast"
)

func StringListToSet(list []string) mapset.Set {
	set := mapset.NewThreadUnsafeSet()
	if list == nil {
		return set
	}
	for _, s := range list {
		set.Add(s)
	}

	return set
}

func Intersect(a []string, b []string) []string {
	aSet := StringListToSet(a)
	bSet := StringListToSet(b)
	interSet := aSet.Intersect(bSet)
	return StringSetToList(interSet)
}

func Difference(a []string, b []string) []string {
	aSet := StringListToSet(a)
	bSet := StringListToSet(b)
	interSet := aSet.Difference(bSet)
	return StringSetToList(interSet)
}

func Uint64ListToSet(list []uint64) mapset.Set {
	set := mapset.NewThreadUnsafeSet()

	for _, s := range list {
		set.Add(s)
	}

	return set
}

func IntListToSet(list []int) mapset.Set {
	set := mapset.NewThreadUnsafeSet()

	for _, s := range list {
		set.Add(s)
	}

	return set
}

func Uint64SetToList(set mapset.Set) []uint64 {
	var res []uint64
	for _, s := range set.ToSlice() {
		res = append(res, cast.ToUint64(s))
	}

	return res
}

func StringSetToList(set mapset.Set) []string {
	var res []string
	if set == nil {
		return []string{}
	}
	for _, s := range set.ToSlice() {
		res = append(res, cast.ToString(s))
	}

	return res
}

func StringSetToListE(set mapset.Set) ([]string, error) {
	res := make([]string, 0)

	for _, it := range set.ToSlice() {
		s, err := cast.ToStringE(it)

		if err != nil {
			return res, err
		}

		res = append(res, s)
	}

	return res, nil
}
