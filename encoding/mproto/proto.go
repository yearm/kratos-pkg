package mproto

import (
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Map ...
func Map(m proto.Message, opts ...protojson.MarshalOptions) map[string]interface{} {
	return toMap(m, opts...)
}

// Maps ...
func Maps(m interface{}, opts ...protojson.MarshalOptions) []map[string]interface{} {
	ms := gconv.Interfaces(m)
	result := make([]map[string]interface{}, 0, len(ms))
	for _, _m := range ms {
		if pm, ok := _m.(proto.Message); ok {
			result = append(result, toMap(pm, opts...))
		}
	}
	return result
}

// KeyMaps ...
func KeyMaps(key string, m proto.Message, opts ...protojson.MarshalOptions) []map[string]interface{} {
	_m := toMap(m, opts...)
	ms, ok := _m[key].([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	return gconv.Maps(ms)
}

// toMap ...
func toMap(m proto.Message, opts ...protojson.MarshalOptions) map[string]interface{} {
	var marshalOptions = protojson.MarshalOptions{EmitUnpopulated: true}
	if len(opts) > 0 {
		marshalOptions = opts[0]
	}
	mb, _ := marshalOptions.Marshal(m)
	return gconv.Map(mb)
}
