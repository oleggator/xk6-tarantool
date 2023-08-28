package tarantool

import (
	"github.com/dop251/goja"
	"github.com/tarantool/go-tarantool/v2"
)

func (c *Client) Call(functionName string, args []any) *goja.Promise {
	req := tarantool.NewCallRequest(functionName).Args(args).Context(c.vu.Context())

	//metrics.PushIfNotDone(c.vu.Context(), c.vu.State().Samples, metrics.Sample{})

	return c.do(req)
}
