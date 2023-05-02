package phpbuiltinserver_test

import (
	"fmt"
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

		Expect(os.WriteFile(filepath.Join(workingDir, "some-file.php"), []byte{}, 0644)).To(Succeed())
		detect = phpbuiltinserver.Detect()
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when the working dir contains a .php file", func() {
		it("detection passes", func() {
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
				},
			}))
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
				t.Setenv("COMPOSER", "some/other-file.json")
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
					t.Setenv("BP_PHP_SERVER", "php-server")
				})

				it("detection passes", func() {
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
						},
					}))
				})
			})

			context("set to something else", func() {
				it.Before(func() {
					t.Setenv("BP_PHP_SERVER", "different-server")
				})

				it("detection fails", func() {
					_, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					Expect(err).To(MatchError(packit.Fail.WithMessage("BP_PHP_SERVER is not set to 'php-server'")))
				})
			})
		})
	})

	context("when BP_PHP_WEB_DIR is set", func() {
		it.Before(func() {
			t.Setenv("BP_PHP_WEB_DIR", "web-dir")
			Expect(os.Mkdir(filepath.Join(workingDir, "web-dir"), os.ModePerm)).To(Succeed())
		})

		it.After(func() {
			Expect(os.RemoveAll(filepath.Join(workingDir, "web-dir"))).To(Succeed())
		})

		context("the web dir contains a .php file", func() {
			it.Before(func() {
				Expect(os.WriteFile(filepath.Join(workingDir, "web-dir", "some-file.php"), []byte{}, os.ModePerm)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Remove(filepath.Join(workingDir, "web-dir", "some-file.php"))).To(Succeed())
			})
			it("detection passes", func() {
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
					},
				}))
			})
		})

		context("the web dir does not contain a .php file", func() {
			it("detection fails", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(fmt.Sprintf("no *.php files found at: %s", filepath.Join(workingDir, "web-dir"))))
			})
		})
	})

	context("when there is no .php file in the working directory", func() {
		it.Before(func() {
			Expect(os.Remove(filepath.Join(workingDir, "some-file.php"))).To(Succeed())
		})

		it("detection fails", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(fmt.Sprintf("no *.php files found at: %s", workingDir)))
		})
	})

	context("failure cases", func() {
		context("the web directory .php file path cannot be globbed", func() {
			it.Before(func() {
				t.Setenv("BP_PHP_WEB_DIR", "\\")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("syntax error in pattern")))
			})
		})

		context("$COMPOSER is set to a non-existent file", func() {
			it.Before(func() {
				Expect(os.Mkdir(filepath.Join(workingDir, "composer-dir"), 0000)).To(Succeed())
				t.Setenv("COMPOSER", filepath.Join("composer-dir", "some-composer"))

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
