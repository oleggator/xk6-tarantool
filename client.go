package tarantool

import (
	"errors"
	"github.com/dop251/goja"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

type Client struct {
	vu   modules.VU
	pool *pool.ConnectionPool

	addrs    []string
	opts     *tarantool.Opts
	promises chan *Promise
}

type Opts struct {
	Addrs    []string `js:"addrs"`
	User     string   `js:"user"`
	Password string   `js:"user"`
	Readers  int      `js:"readers"`
}

func (mi *ModuleInstance) NewClient(call goja.ConstructorCall) *goja.Object {
	rt := mi.vu.Runtime()

	var opts Opts
	if err := rt.ExportTo(call.Arguments[0], &opts); err != nil {
		common.Throw(rt, errors.New("unable to parse options object"))
	}

	client := &Client{
		vu:       mi.vu,
		promises: make(chan *Promise, 1024),
		addrs:    opts.Addrs,
		opts: &tarantool.Opts{
			User: opts.User,
			Pass: opts.Password,
		},
	}

	if opts.Readers == 0 {
		opts.Readers = 1
	}
	for i := 0; i < opts.Readers; i++ {
		go client.reader()
	}

	return rt.ToValue(client).ToObject(rt)
}

func (c *Client) reader() {
	for promise := range c.promises {
		resp, err := promise.tarantoolFuture.Get()
		if err != nil {
			promise.reject(err)
			continue
		}

		promise.resolve(resp.Data)
	}
}

type Promise struct {
	tarantoolFuture *tarantool.Future
	resolve         func(any)
	reject          func(any)
}

func (c *Client) do(req tarantool.Request) *goja.Promise {
	runtime := c.vu.Runtime()
	p, resolve, reject := runtime.NewPromise()

	if err := c.connect(); err != nil {
		reject(err)
		return p
	}

	enqueueToRunOnEventLoop := c.vu.RegisterCallback()
	// resolve and reject should be called inside the event loop because they are not thread-safe
	resolveFn := func(i any) {
		enqueueToRunOnEventLoop(func() error {
			resolve(i)
			return nil
		})
	}
	rejectFn := func(i any) {
		enqueueToRunOnEventLoop(func() error {
			reject(i)
			return nil
		})
	}

	tntPromise := &Promise{
		tarantoolFuture: c.pool.Do(req, pool.ANY),
		resolve:         resolveFn,
		reject:          rejectFn,
	}
	select {
	case <-c.vu.Context().Done():
		reject(c.vu.Context().Err())
	case c.promises <- tntPromise:
	}

	return p
}

// connect establishes the client's connection to the target
// tarantool instance(s).
func (c *Client) connect() error {
	// A nil VU state indicates we are in the init context.
	// As a general convention, k6 should not perform IO in the
	// init context. Thus, the Connect method will error if
	// called in the init context.
	// https://github.com/grafana/k6/issues/2719#issuecomment-1280033675
	vuState := c.vu.State()
	if vuState == nil {
		return common.NewInitContextError("connecting to a tarantool server in the init context is not supported")
	}

	// If the pool is already instantiated, it is safe
	// to assume that the connection is already established.
	if c.pool != nil {
		return nil
	}

	opts := *c.opts
	// use k6's lib.DialerContexter function as tarantool's
	// client Dialer
	// opts.Dialer = vuState.Dialer

	poolConn, err := pool.Connect(c.addrs, opts)
	if err != nil {
		return err
	}

	// Replace the internal tarantool client instance with a new
	// one using our custom options.
	c.pool = poolConn

	return nil
}
