package phpbuiltinserver

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		if server, ok := os.LookupEnv("BP_PHP_SERVER"); ok {
			if server != "php-server" {
				return packit.DetectResult{}, packit.Fail
			}
		}
		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "php",
						Metadata: map[string]interface{}{
							"launch": true,
						},
					},
				},
			},
		}, nil
	}
}
