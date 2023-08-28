# xk6-tarantool

This is a Tarantool client library for [k6](https://github.com/grafana/k6),
implemented as an extension using the [xk6](https://github.com/grafana/xk6) system.

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

```shell
GOFLAGS=-tags=go_tarantool_ssl_disable CGO_ENABLED=0 \
  go run go.k6.io/xk6/cmd/xk6@latest build master --with github.com/tarantool/xk6-tarantool
```

## Usage

This extension exposes a [promise](https://javascript.info/promise-basics)-based API. As opposed to most other current
k6 modules and extensions, who operate in a synchronous manner,
xk6-tarantool operates in an asynchronous manner. In practice, this means that using the Tarantool client's methods won't
block the execution of the test,
and that the test will continue to run even if the Tarantool client is not ready to respond to the request.

## API

xk6-tarantool exposes a subset of Tarantool's commands the core team judged relevant in the context of k6 scripts.

| Tarantool Command | Module function signature                                    | Description | Returns |
|-------------------|:-------------------------------------------------------------|:------------|:--------|
| **call**          | `call(function_name: string, args: any[]) => Promise<any[]>` |             |         |
