package phpbuiltinserver_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	phpbuiltinserver "github.com/paketo-buildpacks/php-builtin-server"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir  string
		workingDir string
		cnbDir     string
		buffer     *bytes.Buffer

		build packit.BuildFunc
	)

	it.Before(func() {
		var err error
		layersDir, err = os.MkdirTemp("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = os.MkdirTemp("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		buffer = bytes.NewBuffer(nil)
		logger := scribe.NewEmitter(buffer)
		build = phpbuiltinserver.Build(logger)
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("no web dir is specified", func() {
		it("returns a result that provides a PHP server start command on the workingDir", func() {
			result, err := build(packit.BuildContext{
				WorkingDir: workingDir,
				CNBPath:    cnbDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{},
				},
				Layers: packit.Layers{Path: layersDir},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: nil,
				},
				Layers: nil,
				Launch: packit.LaunchMetadata{
					Processes: []packit.Process{
						{
							Type:    "web",
							Command: "bash",
							Args: []string{
								"-c",
								fmt.Sprintf(`php -S 0.0.0.0:"${PORT:-80}" -t %s`, workingDir),
							},
							Default: true,
							Direct:  true,
						},
					},
				},
			}))
			Expect(buffer.String()).To(ContainSubstring("Some Buildpack some-version"))
			Expect(buffer.String()).To(ContainSubstring("Assigning launch processes:"))
			Expect(buffer.String()).To(ContainSubstring(fmt.Sprintf(`web (default): bash -c php -S 0.0.0.0:"${PORT:-80}" -t %s`, workingDir)))
		})
	})

	context("a web directory is specified via $BP_PHP_WEB_DIR", func() {
		it.Before(func() {
			t.Setenv("BP_PHP_WEB_DIR", "some-web-dir")
		})

		it("returns a result that provides a PHP server start command on $BP_PHP_WEB_DIR", func() {
			result, err := build(packit.BuildContext{
				WorkingDir: workingDir,
				CNBPath:    cnbDir,
				Stack:      "some-stack",
				BuildpackInfo: packit.BuildpackInfo{
					Name:    "Some Buildpack",
					Version: "some-version",
				},
				Plan: packit.BuildpackPlan{
					Entries: []packit.BuildpackPlanEntry{},
				},
				Layers: packit.Layers{Path: layersDir},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: nil,
				},
				Layers: nil,
				Launch: packit.LaunchMetadata{
					Processes: []packit.Process{
						{
							Type:    "web",
							Command: "bash",
							Args: []string{
								"-c",
								`php -S 0.0.0.0:"${PORT:-80}" -t some-web-dir`,
							},
							Default: true,
							Direct:  true,
						},
					},
				},
			}))
			Expect(buffer.String()).To(ContainSubstring("Some Buildpack some-version"))
			Expect(buffer.String()).To(ContainSubstring("Assigning launch processes:"))
			Expect(buffer.String()).To(ContainSubstring(`web (default): bash -c php -S 0.0.0.0:"${PORT:-80}" -t some-web-dir`))
		})

	})
}
