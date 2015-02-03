# pkg

[![Build Status](https://travis-ci.org/mewkiz/pkg.svg?branch=master)](https://travis-ci.org/mewkiz/pkg)
[![Coverage Status](https://img.shields.io/coveralls/mewkiz/pkg.svg)](https://coveralls.io/r/mewkiz/pkg?branch=master)
[![GoDoc](https://godoc.org/github.com/mewkiz/pkg?status.svg)](https://godoc.org/github.com/mewkiz/pkg)

The pkg project provides packages for various utility functions and commonly used features.

## Mature packages

Packages which reach a certain level of maturity are moved to dedicated repositories under the [mewpkg] organization. This enables API stability through the use of [Semantic Versioning](http://semver.org/).

[mewpkg]: https://github.com/mewpkg/

The following packages have been moved:

- [bits]: provides bit reading operations and binary decoding algorithms.
- [hashutil]: provides utility interfaces for hash functions.
	- [crc8][hashutil/crc8]: implements the 8-bit cyclic redundancy check, or CRC-8, checksum.
	- [crc16][hashutil/crc16]: implements the 16-bit cyclic redundancy check, or CRC-16, checksum.

[bits]: http://godoc.org/github.com/mewpkg/bits
[hashutil]: http://godoc.org/github.com/mewpkg/hashutil
[hashutil/crc8]: http://godoc.org/github.com/mewpkg/hashutil/crc8
[hashutil/crc16]: http://godoc.org/github.com/mewpkg/hashutil/crc16

## Documentation

Documentation provided by GoDoc.

- [bufioutil]: implements utility functions for buffered I/O.
- [bytesutil]: implements some bytes utility functions.
- [dbg]: implements formatted I/O which can be enabled or disabled at runtime.
- [errorsutil]: implements some errors utility functions.
- [errutil]: implements some error utility functions.
- [geometry]: implements basic geometric types and operations.
- [goutil]: implements some golang relevant utility functions.
- [htmlutil]: implements some html utility functions.
- [httputil]: implements some http utility functions.
- [imgutil]: implements some image utility functions.
- [osutil]: implements some os utility functions.
- [pathutil]: implements path utility functions.
- [proxy]: provides proxy server utility functions.
- [readerutil]: implements io.Reader utility functions.
- [stringsutil]: implements some strings utility functions.
- [term]: implements colored output and size measurements for terminals.

[bufioutil]: http://godoc.org/github.com/mewkiz/pkg/bufioutil
[bytesutil]: http://godoc.org/github.com/mewkiz/pkg/bytesutil
[dbg]: http://godoc.org/github.com/mewkiz/pkg/dbg
[errorsutil]: http://godoc.org/github.com/mewkiz/pkg/errorsutil
[errutil]: http://godoc.org/github.com/mewkiz/pkg/errutil
[geometry]: http://godoc.org/github.com/mewkiz/pkg/geometry
[goutil]: http://godoc.org/github.com/mewkiz/pkg/goutil
[htmlutil]: http://godoc.org/github.com/mewkiz/pkg/htmlutil
[httputil]: http://godoc.org/github.com/mewkiz/pkg/httputil
[imgutil]: http://godoc.org/github.com/mewkiz/pkg/imgutil
[osutil]: http://godoc.org/github.com/mewkiz/pkg/osutil
[pathutil]: http://godoc.org/github.com/mewkiz/pkg/pathutil
[proxy]: http://godoc.org/github.com/mewkiz/pkg/proxy
[readerutil]: http://godoc.org/github.com/mewkiz/pkg/readerutil
[stringsutil]: http://godoc.org/github.com/mewkiz/pkg/stringsutil
[term]: http://godoc.org/github.com/mewkiz/pkg/term

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
