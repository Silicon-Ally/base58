// Command base58 is a CLI for encoding and decoding base58 data.
//
// Note: This command currently reads the whole file/input before decoding and
// will not work on infinite (or very large) streaming data.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Silicon-Ally/base58"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		decode = fs.Bool("decode", false, "If true, treat the given input as base58-encoded and decode it.")
	)
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	var (
		src []byte
		err error
	)
	switch n := fs.NArg(); n {
	case 0:
		// Read from stdin
		if src, err = ioutil.ReadAll(os.Stdin); err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	case 1:
		// Read from specified file
		if src, err = os.ReadFile(fs.Arg(0)); err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	default:
		return fmt.Errorf("unexpected number of args %d. usage ./base58 [OPTION]... [FILE]", n)
	}

	if *decode {
		dat, ok := base58.Decode(string(src))
		if !ok {
			return errors.New("invalid base58 given, failed to decode")
		}
		fmt.Print(string(dat))
		return nil
	}

	// If we're here, encode
	fmt.Print(base58.Encode(src))
	return nil
}
