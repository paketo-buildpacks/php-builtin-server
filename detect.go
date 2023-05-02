package phpbuiltinserver

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		webDir := context.WorkingDir
		if wd, ok := os.LookupEnv("BP_PHP_WEB_DIR"); ok {
			webDir = filepath.Join(context.WorkingDir, wd)
		}

		files, err := filepath.Glob(filepath.Join(webDir, "*.php"))
		if err != nil {
			return packit.DetectResult{}, err
		}

		if len(files) == 0 {
			return packit.DetectResult{}, packit.Fail.WithMessage(fmt.Sprintf("no *.php files found at: %s", webDir))
		}

		if server, ok := os.LookupEnv("BP_PHP_SERVER"); ok {
			if server != "php-server" {
				return packit.DetectResult{}, packit.Fail.WithMessage("BP_PHP_SERVER is not set to 'php-server'")
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
