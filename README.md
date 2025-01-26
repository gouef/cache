<img align=right width="168" src="docs/gouef_logo.png">

# Cache
`Cache` is a lightweight, flexible, and efficient caching library for Go, designed to support multiple backends, including [in-memory](docs/Memory.md), [file-based](docs/File.md), and [Redis](docs/Redis.md). It provides a unified interface for caching operations, making it simple to switch between storage solutions without changing your application logic.

Whether you need a fast, ephemeral cache for high-performance applications or a persistent cache for long-term storage, this library has you covered.

With a focus on simplicity and ease of use, Cache enables developers to manage cached data effortlessly, improving application performance and reducing latency.

This package implements [github.com/gouef/standards](https://github.com/gouef/standards).

[![Static Badge](https://img.shields.io/badge/Github-gouef%2Fcache-blue?style=for-the-badge&logo=github&link=github.com%2Fgouef%2Fcache)](https://github.com/gouef/cache)

[![GoDoc](https://pkg.go.dev/badge/github.com/gouef/cache.svg)](https://pkg.go.dev/github.com/gouef/cache)
[![GitHub stars](https://img.shields.io/github/stars/gouef/cache?style=social)](https://github.com/gouef/cache/stargazers)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouef/cache)](https://goreportcard.com/report/github.com/gouef/cache)
[![codecov](https://codecov.io/github/gouef/cache/branch/main/graph/badge.svg?token=YUG8EMH6Q8)](https://codecov.io/github/gouef/cache)

## Versions
![Stable Version](https://img.shields.io/github/v/release/gouef/cache?label=Stable&labelColor=green)
![GitHub Release](https://img.shields.io/github/v/release/gouef/cache?label=RC&include_prereleases&filter=*rc*&logoSize=diago)
![GitHub Release](https://img.shields.io/github/v/release/gouef/cache?label=Beta&include_prereleases&filter=*beta*&logoSize=diago)

## Installation

To use this package in your project, add it using Go modules:

```bash
go get -u github.com/gouef/cache
```

## Usages
- [File](docs/File.md)
- [Memory](docs/Memory.md)
- [Redis](docs/Redis.md)
- [Storage](docs/Storage.md)

## Contributing

Read [Contributing](CONTRIBUTING.md)

## Contributors

<div>
<span>
  <a href="https://github.com/JanGalek"><img src="https://raw.githubusercontent.com/gouef/cache/refs/heads/contributors-svg/.github/contributors/JanGalek.svg" alt="JanGalek" /></a>
</span>
<span>
  <a href="https://github.com/actions-user"><img src="https://raw.githubusercontent.com/gouef/cache/refs/heads/contributors-svg/.github/contributors/actions-user.svg" alt="actions-user" /></a>
</span>
</div>

