# base58

_Note: While a variant of base58 is used in Bitcoin and parts of this take inspiration from [that C++ implementation](https://github.com/bitcoin/bitcoin/blob/15db77f4dd7f1a7963398f1576580b577a1697bc/src/base58.cpp), this package has nothing to do with cryptocurrency._

[![GoDoc](https://pkg.go.dev/badge/github.com/Silicon-Ally/base58?status.svg)](https://pkg.go.dev/github.com/Silicon-Ally/base58?tab=doc)
[![CI Workflow](https://github.com/Silicon-Ally/base58/actions/workflows/test.yml/badge.svg)](https://github.com/Silicon-Ally/base58/actions?query=branch%3Amain)

`base58` is simple, zero-dependency Go library that implements a base 58 encoding scheme. The 58 characters it uses are:

```
123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz
```

This may or may not be the same set or ordering of characters as other base 58 encoding packages, make sure to double-check if you require interoperability. Generally, these characters were chosen to remove ambiguity (e.g. `0 vs O vs o`, `I vs l vs L`) and general complications posed by non-alphanumerics in various contexts when using base 64 encodings.

# When would I want to use this?

This encoding is useful when you want something more compact than a hexadecimal encoding, but without standard base64's `+` and `/`, which can be problematic in user-facing web contexts (e.g. URLs). Similarly, URL-safe base64 variants are more annoying to double click, as `-` stops the selection.

In short, you might want something like this in contexts where your encoded data may appear in a URL, or you expect a user to manually copy and paste the encoded data.

## Contributing

Contribution guidelines can be found [on our website](https://siliconally.org/oss/contributor-guidelines).
