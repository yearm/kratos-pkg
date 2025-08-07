package gjson

import (
	"encoding/json"

	"github.com/yearm/kratos-pkg/utils/bytesconv"
)

// MustMarshal converts a value to its JSON byte representation, ignoring any errors.
func MustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

// MustMarshalToString converts a value to its JSON string representation, ignoring any errors.
func MustMarshalToString(v any) string {
	b, _ := json.Marshal(v)
	return bytesconv.BytesToString(b)
}
