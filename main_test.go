package main

import "testing"

func Test_maskJSON(t *testing.T) {
	cases := []struct {
		test string
		in   string
		out  string
	}{
		{
			test: "simple",
			in:   `{"email":"test@example.com","password":"foobarbaz"}`,
			out:  `{"email":"test@example.com","password":"[FILTERED]"}`,
		},
		{
			test: "nested",
			in:   `{"token":"quxquux","user":{"email":"test@example.com","password":"foobarbaz"}}`,
			out:  `{"token":"[FILTERED]","user":{"email":"test@example.com","password":"[FILTERED]"}}`,
		},
		{
			test: "array",
			in:   `{"users":[{"email":"test@example.com","password":"foobarbaz"},{"email":"test2@example.com","password":"quxquux"}]}`,
			out:  `{"users":[{"email":"test@example.com","password":"[FILTERED]"},{"email":"test2@example.com","password":"[FILTERED]"}]}`,
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			out, err := Mask([]byte(c.in))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			if got, want := string(out), c.out; got != want {
				t.Errorf("Mask() returned %q, want %q", got, want)
			}
		})
	}
}
