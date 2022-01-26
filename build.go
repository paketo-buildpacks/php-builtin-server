package phpbuiltinserver

import (
	"fmt"
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Build(logger scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		webDir := context.WorkingDir
		if wd, ok := os.LookupEnv("BP_PHP_WEB_DIR"); ok {
			webDir = wd
		}

		// Use default port 80 unless $PORT is set
		command := fmt.Sprintf(`php -S 0.0.0.0:"${PORT:-80}" -t %s`, webDir)

		processes := []packit.Process{
			{
				Type:    "web",
				Command: command,
				Default: true,
			},
		}

		logger.Process("Assigning launch processes")
		for _, process := range processes {
			logger.Subprocess("%s: %s", process.Type, process.Command)
		}

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
