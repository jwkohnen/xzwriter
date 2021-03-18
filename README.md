# xzwriter: a trivial Go library for compressing with XZ.

[![Apache License v2.0](https://img.shields.io/badge/license-Apache%20License%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.txt)
[![GoDoc](https://godoc.org/github.com/jwkohnen/xzwriter?status.svg)](https://godoc.org/github.com/jwkohnen/xzwriter)

Package xzwriter provides a writer XZWriter that pipes through an external XZ
compressor.

Expects the Tukaani XZ tool in $PATH. See the XZ Utils home page:
<http://tukaani.org/xz/>

## License
As this is a trivial convenience package, you should probably not import this. Anyway
licensed under the Apache License, Version 2.0.

## See also
https://github.com/ulikunitz/xz
