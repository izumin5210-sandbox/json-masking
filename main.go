package main

import (
	"encoding/json"

	"github.com/a8m/djson"
	"github.com/pquerna/ffjson/ffjson"
)

var (
	keywordSet = map[string]struct{}{
		"password":              struct{}{},
		"password_confirmation": struct{}{},
		"token":                 struct{}{},
	}
	filteredPlaceholder = "[FILTERED]"
)

func MaskWithEncodingJSON(data []byte) ([]byte, error) {
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
				in[k] = filteredPlaceholder
			} else {
				in[k] = mask(v)
			}
		}
		return in
	case []interface{}:
		for i := 0; i < len(in); i++ {
			in[i] = mask(in[i])
		}
		return in
	}
	return in
}

func MaskWithDJSON(data []byte) ([]byte, error) {
	v, err := djson.Decode(data)
	if err != nil {
		return nil, err
	}

	out, err := json.Marshal(mask(v))
	if err != nil {
		return nil, err
	}

	return out, nil
}

func MaskWithFFJSON(data []byte) ([]byte, error) {
	var v interface{}

	err := ffjson.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	out, err := ffjson.Marshal(mask(v))
	if err != nil {
		return nil, err
	}

	return out, nil
}
