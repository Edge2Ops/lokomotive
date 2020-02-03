package flatcarlinuxupdateoperator

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/pkg/errors"

	"github.com/kinvolk/lokoctl/pkg/assets"
	"github.com/kinvolk/lokoctl/pkg/components"
	"github.com/kinvolk/lokoctl/pkg/util/walkers"
)

const componentName = "flatcar-linux-update-operator"

func init() {
	components.Register(componentName, &component{})
}

type component struct{}

func (c *component) LoadConfig(configBody *hcl.Body, evalContext *hcl.EvalContext) hcl.Diagnostics {
	if configBody == nil {
		// This component has no configuration, so don't complain when there is no configuration defined.
		return nil
	}
	return gohcl.DecodeBody(*configBody, evalContext, c)
}

func (c *component) RenderManifests() (map[string]string, error) {
	ret := make(map[string]string)
	walk := walkers.DumpingWalker(ret, ".yaml")
	if err := assets.Assets.WalkFiles(fmt.Sprintf("/components/%s/manifests", componentName), walk); err != nil {
		return nil, errors.Wrap(err, "failed to walk assets")
	}

	return ret, nil
}

func (c *component) Metadata() components.Metadata {
	return components.Metadata{
		Namespace: "reboot-coordinator",
		Helm:      &components.HelmMetadata{},
	}
}
