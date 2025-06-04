package proto

import (
	"encoding/json"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ToMap converts a protobuf message into a map[string]interface{}.
func ToMap(pm proto.Message, opts ...protojson.MarshalOptions) (map[string]interface{}, error) {
	return toMap(pm, opts...)
}

// ToMaps converts a slice of protobuf messages into a slice of map[string]interface{}.
func ToMaps(pms []proto.Message, opts ...protojson.MarshalOptions) ([]map[string]interface{}, error) {
	ms := make([]map[string]interface{}, 0, len(pms))
	for _, pm := range pms {
		m, err := toMap(pm, opts...)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// ToKeyMaps extracts a nested slice of map[string]interface{} from a protobuf message.
func ToKeyMaps(pm proto.Message, key string, opts ...protojson.MarshalOptions) ([]map[string]interface{}, error) {
	m, err := toMap(pm, opts...)
	if err != nil {
		return nil, err
	}
	ms, ok := m[key].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}
	return lo.FilterMap(ms, func(item interface{}, _ int) (map[string]interface{}, bool) {
		v, ok := item.(map[string]interface{})
		return v, ok
	}), nil
}

// toMap converts a protobuf message into a map[string]interface{}.
func toMap(pm proto.Message, opts ...protojson.MarshalOptions) (map[string]interface{}, error) {
	var marshalOptions = protojson.MarshalOptions{EmitUnpopulated: true}
	if len(opts) > 0 {
		marshalOptions = opts[0]
	}
	mb, err := marshalOptions.Marshal(pm)
	if err != nil {
		return nil, errors.Wrap(err, "marshalOptions.Marshal failed")
	}
	m := make(map[string]interface{})
	if err = json.Unmarshal(mb, &m); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal failed")
	}
	return m, nil
}
