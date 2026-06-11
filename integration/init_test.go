package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var (
	buildpack        string
	phpDistBuildpack string
	buildpackInfo    struct {
		Buildpack struct {
			ID   string
			Name string
		}
	}
	Config struct {
		PhpDist string `json:"php-dist"`
	}
)

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	integrationFile, err := os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(integrationFile).Decode(&Config)).To(Succeed())
	Expect(integrationFile.Close()).To(Succeed())

	buildpackFile, err := os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(buildpackFile).Decode(&buildpackInfo)
	Expect(err).NotTo(HaveOccurred())
	Expect(buildpackFile.Close()).To(Succeed())

	root, err := filepath.Abs("./..")
	Expect(err).ToNot(HaveOccurred())

	buildpackStore := occam.NewBuildpackStore()
	targetedBuildpackStore := buildpackStore.WithTarget("linux/" + runtime.GOARCH)

	buildpack, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).ToNot(HaveOccurred())

	phpDistBuildpack, err = targetedBuildpackStore.Get.
		Execute(Config.PhpDist)
	Expect(err).ToNot(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Default", testDefault)
	suite.Run(t)
}
