package tarantool

import (
	"fmt"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/tarantool` module instances for each VU.
type RootModule struct{}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	metrics, err := registerMetrics(vu.InitEnv().Registry)
	if err != nil {
		common.Throw(vu.Runtime(), fmt.Errorf("failed to register GRPC module metrics: %w", err))
	}

	return &ModuleInstance{vu: vu, metrics: metrics}
}

type ModuleInstance struct {
	vu      modules.VU
	metrics *instanceMetrics
}

// Exports implements the modules.Instance interface and returns
// the exports of the JS module.
func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{Named: map[string]interface{}{
		"Client": mi.NewClient,
	}}
}
