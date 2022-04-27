package phpbuiltinserver

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		if server, ok := os.LookupEnv("BP_PHP_SERVER"); ok {
			if server != "php-server" {
				return packit.DetectResult{}, packit.Fail
			}
		}

		requirements := []packit.BuildPlanRequirement{
			{
				Name: "php",
				Metadata: map[string]interface{}{
					"launch": true,
				},
			},
		}

		composerJsonPath := filepath.Join(context.WorkingDir, "composer.json")

		if value, found := os.LookupEnv("COMPOSER"); found {
			composerJsonPath = filepath.Join(context.WorkingDir, value)
		}

		if exists, err := fs.Exists(composerJsonPath); err != nil {
			return packit.DetectResult{}, err
		} else if exists {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name: "composer-packages",
				Metadata: map[string]interface{}{
					"launch": true,
				},
			})
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: requirements,
			},
		}, nil
	}
}
