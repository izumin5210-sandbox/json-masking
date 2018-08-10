package main

import "testing"

var (
	cases = []struct {
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
)

func testMask(tb testing.TB, maskFn func([]byte) ([]byte, error), in, out string) {
	tb.Helper()
	masked, err := maskFn([]byte(in))
	if err != nil {
		tb.Fatalf("Unexpected error %v", err)
	}
	if got, want := string(masked), out; got != want {
		tb.Errorf("Mask() returned %q, want %q", got, want)
	}
}

func BenchmarkMask_EncodingJSON(b *testing.B) {
	for _, c := range cases {
		b.Run(c.test, func(b *testing.B) {
			testMask(b, MaskWithEncodingJSON, c.in, c.out)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MaskWithEncodingJSON([]byte(c.in))
			}
		})
	}
}

func BenchmarkMask_DJSON(b *testing.B) {
	for _, c := range cases {
		b.Run(c.test, func(b *testing.B) {
			testMask(b, MaskWithDJSON, c.in, c.out)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MaskWithDJSON([]byte(c.in))
			}
		})
	}
}
