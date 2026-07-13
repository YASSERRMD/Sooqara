package pipeline

import "encoding/json"

// jsonMarshal serializes a value to a JSON string.
func jsonMarshal(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// jsonUnmarshal deserializes a JSON string into a value.
func jsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
