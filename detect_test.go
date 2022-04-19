package phpbuiltinserver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	phpbuiltinserver "github.com/paketo-buildpacks/php-builtin-server"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		workingDir string
		detect     packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())
		detect = phpbuiltinserver.Detect()
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
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

	context("there is a composer.json file", func() {
		it.Before(func() {
			Expect(os.WriteFile(filepath.Join(workingDir, "composer.json"), []byte(""), os.ModePerm)).To(Succeed())
		})
		it.After(func() {
			Expect(os.RemoveAll(filepath.Join(workingDir, "composer.json"))).To(Succeed())
		})

		it("detection passes, and requires composer-packages", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
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
					{
						Name: "composer-packages",
						Metadata: map[string]interface{}{
							"launch": true,
						},
					},
				},
			}))
		})
	})

	context("$COMPOSER is set to an existent file", func() {
		it.Before(func() {
			Expect(os.Setenv("COMPOSER", "some/other-file.json")).To(Succeed())
			Expect(os.Mkdir(filepath.Join(workingDir, "some"), os.ModeDir|os.ModePerm)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(workingDir, "some", "other-file.json"), []byte(""), os.ModePerm)).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("COMPOSER")).To(Succeed())
			Expect(os.RemoveAll(filepath.Join(workingDir, "some"))).To(Succeed())
		})

		it("detection passes, and requires composer-packages", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
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
					{
						Name: "composer-packages",
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

	context("failure cases", func() {
		context("$COMPOSER is set to a non-existent file", func() {
			it.Before(func() {
				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
				Expect(os.Setenv("COMPOSER", filepath.Join(workingDir, "other-file.json"))).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("COMPOSER")).To(Succeed())
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})
	})
}
