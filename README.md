# iexcloud

Go library for accessing the IEX Cloud API.

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![License Badge][license badge]][license]

## Overview

[iexcloud][] provides a Go interface to the [IEX Cloud API][iexcloudio]. To
access the [IEX Cloud API][iexcloudio] an account and token are required. The
goal is for iexcloud to be compatible with the v1 version of the IEX Cloud API.
There were some changes from the beta version to v1 of the API, so things may
still be in flux for this library.

## Installation

```bash
$ go get github.com/woodstock-tokyo/iexcloud/v2
```

### Testing

Prior to submitting a [pull request][], please run:

```bash
$ make check
```

To update and view the test coverage report:

```bash
$ make cover
```

## License

[iexcloud][] is released under the MIT license. Please see the
[LICENSE][] file for more information.

[iexcloudio]: https://iexcloud.io
[iexcloud]: https://github.com/woodstock-tokyo/iexcloud
[godoc badge]: https://godoc.org/github.com/woodstock-tokyo/iexcloud?status.svg
[godoc link]: https://godoc.org/github.com/woodstock-tokyo/iexcloud
[implementation]: https://github.com/woodstock-tokyo/iexcloud/blob/master/implementation.md
[license]: https://github.com/woodstock-tokyo/iexcloud/blob/master/LICENSE
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/woodstock-tokyo/iexcloud
[report card]: https://goreportcard.com/report/github.com/woodstock-tokyo/iexcloud
