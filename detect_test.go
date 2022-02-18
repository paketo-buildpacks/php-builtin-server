package phpbuiltinserver_test

import (
	"os"
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

	context("when no configuration is set", func() {
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
	})

	context("BP_PHP_SERVER", func() {
		context("set to php-server", func() {
			it.Before(func() {
				os.Setenv("BP_PHP_SERVER", "php-server")
			})
			it.After(func() {
				os.Unsetenv("BP_PHP_SERVER")
			})

			it("detection passes", func() {
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
		})

		context("set to something else", func() {
			it.Before(func() {
				os.Setenv("BP_PHP_SERVER", "different-server")
			})
			it.After(func() {
				os.Unsetenv("BP_PHP_SERVER")
			})

			it("detection fails", func() {
				_, err := detect(packit.DetectContext{})
				Expect(err).To(MatchError(packit.Fail))
			})
		})
	})
}
