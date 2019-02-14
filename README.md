# fuid [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0) [![release](https://img.shields.io/github/release/kpango/fuid.svg?style=flat-square)](https://github.com/kpango/fuid/releases/latest) [![CircleCI](https://circleci.com/gh/kpango/fuid.svg)](https://circleci.com/gh/kpango/fuid) [![codecov](https://codecov.io/gh/kpango/fuid/branch/master/graph/badge.svg?token=2CzooNJtUu&style=flat-square)](https://codecov.io/gh/kpango/fuid) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/890d1b3e9bef4b9e9219894e80b4085f)](https://www.codacy.com/app/i.can.feel.gravity/fuid?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=kpango/fuid&amp;utm_campaign=Badge_Grade) [![Go Report Card](https://goreportcard.com/badge/github.com/kpango/fuid)](https://goreportcard.com/report/github.com/kpango/fuid) [![GolangCI](https://golangci.com/badges/github.com/kpango/fuid.svg?style=flat-square)](https://golangci.com/r/github.com/kpango/fuid) [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/kpango/fuid) [![GoDoc](http://godoc.org/github.com/kpango/fuid?status.svg)](http://godoc.org/github.com/kpango/fuid) [![DepShield Badge](https://depshield.sonatype.org/badges/kpango/fuid/depshield.svg)](https://depshield.github.io)

fuid is simple and fast uuid string library forked from rs/xid

## Requirement
Go 1.11

## Installation
```shell
go get github.com/kpango/fuid
```

## Example
```go
	fuid.String()
```


## Benchmarks
fuid vs rs/xid vs satori/go.uuid vs google/uuid
```lstv
go test -count=5 -run=NONE -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/kpango/fuid
BenchmarkFUID-8         	30000000	        41.4 ns/op	      32 B/op	       1 allocs/op
BenchmarkFUID-8         	30000000	        41.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkFUID-8         	50000000	        39.8 ns/op	      32 B/op	       1 allocs/op
BenchmarkFUID-8         	50000000	        40.7 ns/op	      32 B/op	       1 allocs/op
BenchmarkFUID-8         	30000000	        49.4 ns/op	      32 B/op	       1 allocs/op
BenchmarkXID-8          	20000000	        57.8 ns/op	      32 B/op	       1 allocs/op
BenchmarkXID-8          	30000000	        62.3 ns/op	      32 B/op	       1 allocs/op
BenchmarkXID-8          	20000000	        58.8 ns/op	      32 B/op	       1 allocs/op
BenchmarkXID-8          	20000000	        60.2 ns/op	      32 B/op	       1 allocs/op
BenchmarkXID-8          	30000000	        59.5 ns/op	      32 B/op	       1 allocs/op
BenchmarkSatoriUUID-8   	 3000000	       486 ns/op	      16 B/op	       1 allocs/op
BenchmarkSatoriUUID-8   	 3000000	       490 ns/op	      16 B/op	       1 allocs/op
BenchmarkSatoriUUID-8   	 3000000	       502 ns/op	      16 B/op	       1 allocs/op
BenchmarkSatoriUUID-8   	 3000000	       492 ns/op	      16 B/op	       1 allocs/op
BenchmarkSatoriUUID-8   	 3000000	       501 ns/op	      16 B/op	       1 allocs/op
BenchmarkGoogleUUID-8   	 3000000	       499 ns/op	      16 B/op	       1 allocs/op
BenchmarkGoogleUUID-8   	 3000000	       491 ns/op	      16 B/op	       1 allocs/op
BenchmarkGoogleUUID-8   	 3000000	       491 ns/op	      16 B/op	       1 allocs/op
BenchmarkGoogleUUID-8   	 3000000	       485 ns/op	      16 B/op	       1 allocs/op
BenchmarkGoogleUUID-8   	 3000000	       485 ns/op	      16 B/op	       1 allocs/op
PASS
ok  	github.com/kpango/fuid	37.002s
```

## Contribution
1. Fork it ( https://github.com/kpango/fuid/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Author
[kpango](https://github.com/kpango)

## LICENSE
fuid released under MIT license, refer [LICENSE](https://github.com/kpango/fuid/blob/master/LICENSE) file.
