package mergeomap

import "github.com/iancoleman/orderedmap"

var (
	MaxDepth = 32
)

func Merge(dst, src *orderedmap.OrderedMap) *orderedmap.OrderedMap {
	return merge(dst, src, 0)
}

func merge(dst, src *orderedmap.OrderedMap, depth int) *orderedmap.OrderedMap {
	if depth > MaxDepth {
		panic("too deep!")
	}
	for key, srcVal := range src.Values() {
		dstVal, ok := dst.Values()[key]
		if ok {
			srcMap, srcMapOk := mapify(srcVal)
			dstMap, dstMapOk := mapify(dstVal)
			if srcMapOk && dstMapOk {
				srcVal = merge(dstMap, srcMap, depth+1)
			}
			dst.Values()[key] = srcVal
		} else {
			dst.Set(key, srcVal)
		}
	}
	return dst
}

func mapify(i interface{}) (*orderedmap.OrderedMap, bool) {
	switch iv := i.(type) {
	case orderedmap.OrderedMap:
		m := orderedmap.New()
		for _, k := range iv.Keys() {
			v, _ := iv.Get(k)
			m.Set(k, v)
		}
		return m, true
	case *orderedmap.OrderedMap:
		m := orderedmap.New()
		for _, k := range iv.Keys() {
			v, _ := iv.Get(k)
			m.Set(k, v)
		}
		return m, true
	default:
		return nil, false
	}
}
