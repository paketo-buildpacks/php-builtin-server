package phpbuiltinserver_test

import (
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	phpbuiltinserver "github.com/paketo-buildpacks/php-builtin-server"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		detect packit.DetectFunc
	)

	it.Before(func() {
		detect = phpbuiltinserver.Detect()
	})

	it("detection always passes", func() {
		result, err := detect(packit.DetectContext{})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Plan).To(Equal(packit.BuildPlan{
			Provides: []packit.BuildPlanProvision{},
			Requires: []packit.BuildPlanRequirement{
				{
					Name: "php",
					Metadata: map[string]interface{}{
						"launch": true,
					},
				},
			},
		}))
	})
}
