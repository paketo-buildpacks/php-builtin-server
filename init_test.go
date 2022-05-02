package phpbuiltinserver_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitBuiltinServer(t *testing.T) {
	suite := spec.New("phpbuiltinserver", spec.Report(report.Terminal{}))
	suite("Build", testBuild)
	suite("Detect", testDetect, spec.Sequential())
	suite.Run(t)
}
