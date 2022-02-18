package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	phpbuiltinserver "github.com/paketo-buildpacks/php-builtin-server"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout)
	packit.Run(
		phpbuiltinserver.Detect(),
		phpbuiltinserver.Build(logger),
	)
}
