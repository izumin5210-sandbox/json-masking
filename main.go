package main

import "encoding/json"

var (
	keywordSet = map[string]struct{}{
		"password":              struct{}{},
		"password_confirmation": struct{}{},
		"token":                 struct{}{},
	}
)

func Mask(data []byte) ([]byte, error) {
	var v interface{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(mask(v))
	if err != nil {
		return nil, err
	}

	return out, nil
}

func mask(in interface{}) interface{} {
	switch in := in.(type) {
	case map[string]interface{}:
		for k, v := range in {
			if _, ok := keywordSet[k]; ok {
				in[k] = "[FILTERED]"
			} else {
				in[k] = mask(v)
			}
		}
		return in
	case []interface{}:
		masked := make([]interface{}, 0, len(in))
		for _, v := range in {
			masked = append(masked, mask(v))
		}
		return masked
	}
	return in
}
