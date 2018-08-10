package main

import (
	"encoding/json"
	"testing"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Account struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

var (
	cases = []struct {
		test string
		in   interface{}
		out  string
	}{
		{
			test: "simple",
			in:   &User{Email: "test@example.com", Password: "foobarbaz"},
			out:  `{"email":"test@example.com","password":"[FILTERED]"}`,
		},
		{
			test: "nested",
			in:   &Account{Token: "quxquux", User: &User{Email: "test@example.com", Password: "foobarbaz"}},
			out:  `{"token":"[FILTERED]","user":{"email":"test@example.com","password":"[FILTERED]"}}`,
		},
		{
			test: "array",
			in: struct {
				Users []*User `json:"users"`
			}{Users: []*User{{Email: "test@example.com", Password: "foobarbaz"}, {Email: "test2@example.com", Password: "quxquux"}}},
			out: `{"users":[{"email":"test@example.com","password":"[FILTERED]"},{"email":"test2@example.com","password":"[FILTERED]"}]}`,
		},
	}
)

func testAndBenchMask(b *testing.B, maskFn func([]byte) ([]byte, error), in interface{}, out string) {
	inJSON, err := json.Marshal(in)
	if err != nil {
		b.Fatalf("Unexpected error %v", err)
	}

	masked, err := maskFn(inJSON)
	if err != nil {
		b.Fatalf("Unexpected error %v", err)
	}
	if got, want := string(masked), out; got != want {
		b.Errorf("Mask() returned %q, want %q", got, want)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		maskFn(inJSON)
	}
}

func BenchmarkMask_EncodingJSON(b *testing.B) {
	for _, c := range cases {
		b.Run(c.test, func(b *testing.B) {
			testAndBenchMask(b, MaskWithEncodingJSON, c.in, c.out)
		})
	}
}

func BenchmarkMask_DJSON(b *testing.B) {
	for _, c := range cases {
		b.Run(c.test, func(b *testing.B) {
			testAndBenchMask(b, MaskWithDJSON, c.in, c.out)
		})
	}
}

func BenchmarkMask_FFJSON(b *testing.B) {
	for _, c := range cases {
		b.Run(c.test, func(b *testing.B) {
			testAndBenchMask(b, MaskWithFFJSON, c.in, c.out)
		})
	}
}

func BenchmarkMaskStruct(b *testing.B) {
	for _, c := range cases {
		b.Run(c.test, func(b *testing.B) {
			masked, err := MaskStruct(c.in)
			if err != nil {
				b.Fatalf("Unexpected error %v", err)
			}
			if got, want := string(masked), c.out; got != want {
				b.Errorf("Mask() returned %q, want %q", got, want)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MaskStruct(c.in)
			}
		})
	}
}
