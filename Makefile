build:
	GOFLAGS=-tags=go_tarantool_ssl_disable \
	CGO_ENABLED=0 \
	go run go.k6.io/xk6/cmd/xk6@latest \
		build master --with github.com/tarantool/xk6-tarantool=.
