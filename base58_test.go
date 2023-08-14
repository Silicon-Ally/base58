package base58

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBase58RoundTrip_FromRaw(t *testing.T) {
	tests := []struct {
		in   []byte
		want string
	}{
		{in: []byte{}, want: ""},
		{in: []byte{0}, want: "1"},
		{in: []byte{0, 0}, want: "11"},
		{in: []byte{0, 0, 0}, want: "111"},
		{in: []byte{0, 0, 0, 1}, want: "1112"},
		{in: []byte{0, 0, 0, 1, 1}, want: "1115S"},
		{in: []byte{193}, want: "4L"},
		{in: []byte{179, 43}, want: "Edp"},
		{in: []byte{0, 87, 10}, want: "17dB"},
		{in: []byte{171, 14, 39, 165}, want: "5Nbdwz"},
		{in: []byte{132, 244, 233, 50, 105}, want: "G12Fsdz"},
		{in: []byte{0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 16, 32, 64, 128, 240, 250, 255}, want: "11112drXXUifSrS46koaV2Qv"},
		{in: []byte{57, 245, 240, 179, 0, 56, 101, 159, 186, 111, 148, 194, 58, 169, 5, 106, 209, 20, 63, 230, 255, 12, 162, 79, 124, 47, 244, 218, 182, 238, 220, 6, 2, 245, 56, 112, 114, 219, 78, 146, 91, 239, 56, 2, 236, 33, 70, 228, 254, 210, 168, 193, 59, 115, 22, 223, 65, 251, 81, 3, 70, 178, 48, 217, 105, 81, 61, 150, 39, 248, 11, 121, 76, 60, 241, 89, 75, 32, 41, 49, 224, 150, 146, 255, 149, 216, 2, 44, 114, 240, 88, 146, 45, 171, 100, 179, 135, 6, 135, 221}, want: "3FejHzUxN4SmU8EpehUGWE9YnX96Q8eshf5eaRqbiLWGPHf9kni2eBgxAdNbBQsLzhrf41PvDCxedgSR2Xfnvj2jonoVmV723Fo3A1XeTjpaXs68fcjMHxg72hnQnzqCRTRcNJ29J"},
		{in: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100}, want: "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111112j"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got := Encode(test.in)
			if got != test.want {
				t.Fatalf("ToBase58(%v) = %q, want %q", test.in, got, test.want)
			}
			dec, ok := Decode(got)
			if !ok {
				t.Fatalf("input %q was invalid", got)
			}
			if diff := cmp.Diff(test.in, dec); diff != "" {
				t.Errorf("unexpected dec(enc(v)) round trip (-want +got)\n%s", diff)
			}
		})
	}
}

func TestDecodeInvalid(t *testing.T) {
	tests := []struct {
		desc string
		in   string
	}{
		{
			desc: "random text",
			in:   "this isn't base58 encoded",
		},
		{
			desc: "Capital O",
			in:   "11112drXXUifSrS46kOaV2Qv",
		},
		{
			desc: "Zero",
			in:   "11112drXXUifSrS46k0aV2Qv",
		},
		{
			desc: "Capital I",
			in:   "11112drXXUIfSrS46koaV2Qv",
		},
		{
			desc: "Lowercase l",
			in:   "11l12drXXUifSrS46koaV2Qv",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			dec, ok := Decode(test.in)
			if ok {
				t.Fatalf("Decode = %q, expected decoding failure", dec)
			}
		})
	}
}

func TestBase58RoundTrip_FromEnc(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{in: "", want: []byte{}},
		{in: "1", want: []byte{0}},
		{in: "11", want: []byte{0, 0}},
		{in: "111", want: []byte{0, 0, 0}},
		{in: "1112", want: []byte{0, 0, 0, 1}},
		{in: "1115S", want: []byte{0, 0, 0, 1, 1}},
		{in: "4L", want: []byte{193}},
		{in: "Edp", want: []byte{179, 43}},
		{in: "17dB", want: []byte{0, 87, 10}},
		{in: "5Nbdwz", want: []byte{171, 14, 39, 165}},
		{in: "G12Fsdz", want: []byte{132, 244, 233, 50, 105}},
		{in: "11112drXXUifSrS46koaV2Qv", want: []byte{0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 16, 32, 64, 128, 240, 250, 255}},
		{in: "3FejHzUxN4SmU8EpehUGWE9YnX96Q8eshf5eaRqbiLWGPHf9kni2eBgxAdNbBQsLzhrf41PvDCxedgSR2Xfnvj2jonoVmV723Fo3A1XeTjpaXs68fcjMHxg72hnQnzqCRTRcNJ29J", want: []byte{57, 245, 240, 179, 0, 56, 101, 159, 186, 111, 148, 194, 58, 169, 5, 106, 209, 20, 63, 230, 255, 12, 162, 79, 124, 47, 244, 218, 182, 238, 220, 6, 2, 245, 56, 112, 114, 219, 78, 146, 91, 239, 56, 2, 236, 33, 70, 228, 254, 210, 168, 193, 59, 115, 22, 223, 65, 251, 81, 3, 70, 178, 48, 217, 105, 81, 61, 150, 39, 248, 11, 121, 76, 60, 241, 89, 75, 32, 41, 49, 224, 150, 146, 255, 149, 216, 2, 44, 114, 240, 88, 146, 45, 171, 100, 179, 135, 6, 135, 221}},
		{in: "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111112j", want: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, ok := Decode(test.in)
			if !ok {
				t.Fatalf("input %q was invalid", got)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("unexpected dec(enc(v)) round trip (-want +got)\n%s", diff)
			}
			enc := Encode(got)
			if !ok {
				t.Fatalf("input %q was invalid", got)
			}
			if enc != test.in {
				t.Fatalf("ToBase58(%v) = %q, want %q", test.in, got, test.want)
			}
		})
	}
}

func TestBase58RoundTrip_Pseudorandom(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	for i := 0; i < 13; i++ {
		sz := 2 << i
		t.Run(fmt.Sprintf("#%d - %d", i, sz), func(t *testing.T) {
			inp := make([]byte, sz)
			n, err := r.Read(inp)
			if err != nil {
				t.Fatalf("failed to read pseudorandom bytes: %v", err)
			}
			if n != sz {
				t.Fatalf("read %d bytes, expected %d", n, sz)
			}
			out, ok := Decode(Encode(inp))
			if !ok {
				t.Fatal("Encode produced undecodable output")
			}

			if diff := cmp.Diff(inp, out); diff != "" {
				t.Errorf("unexpected dec(enc(v)) round trip (-want +got)\n%s", diff)
			}
		})
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	// This structure makes it easier to compare the performance of different approaches.
	benchmarkRoundTrip(b, Encode, Decode)
}

func benchmarkRoundTrip(b *testing.B, enc func([]byte) string, dec func(in string) ([]byte, bool)) {
	b.StopTimer()
	sz := 1024 * 16
	r := rand.New(rand.NewSource(0))
	inp := make([]byte, sz)

	for i := 0; i < b.N; i++ {
		n, err := r.Read(inp)
		if err != nil {
			b.Fatalf("failed to read pseudorandom bytes: %v", err)
		}
		if n != sz {
			b.Fatalf("read %d bytes, expected %d", n, sz)
		}

		b.StartTimer()
		out, ok := dec(enc(inp))
		if !ok {
			b.Fatal("Encode produced undecodable output")
		}
		b.StopTimer()

		if !bytes.Equal(out, inp) {
			b.Error("dec(enc(v)) round trip was not equal")
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	b.StopTimer()
	sz := 1024 * 16
	r := rand.New(rand.NewSource(0))
	inp := make([]byte, sz)

	for i := 0; i < b.N; i++ {
		n, err := r.Read(inp)
		if err != nil {
			b.Fatalf("failed to read pseudorandom bytes: %v", err)
		}
		if n != sz {
			b.Fatalf("read %d bytes, expected %d", n, sz)
		}

		b.StartTimer()
		encoded := Encode(inp)
		b.StopTimer()
		out, ok := Decode(encoded)
		if !ok {
			b.Fatal("Encode produced undecodable output")
		}

		if !bytes.Equal(out, inp) {
			b.Error("dec(enc(v)) round trip was not equal")
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	b.StopTimer()
	sz := 1024 * 16
	r := rand.New(rand.NewSource(0))
	inp := make([]byte, sz)

	for i := 0; i < b.N; i++ {
		n, err := r.Read(inp)
		if err != nil {
			b.Fatalf("failed to read pseudorandom bytes: %v", err)
		}
		if n != sz {
			b.Fatalf("read %d bytes, expected %d", n, sz)
		}

		encoded := Encode(inp)
		b.StartTimer()
		out, ok := Decode(encoded)
		if !ok {
			b.Fatal("Encode produced undecodable output")
		}
		b.StopTimer()

		if !bytes.Equal(out, inp) {
			b.Error("dec(enc(v)) round trip was not equal")
		}
	}
}
