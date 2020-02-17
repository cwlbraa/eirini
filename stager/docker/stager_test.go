package docker_test

import (
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/stager/docker"
	"code.cloudfoundry.org/eirini/stager/docker/dockerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DockerStager", func() {

	var (
		stager  docker.Stager
		fetcher *dockerfakes.FakeImageMetadataFetcher
	)

	Context("Stage a docker image", func() {

		BeforeEach(func() {
			fetcher = new(dockerfakes.FakeImageMetadataFetcher)
			stager = docker.Stager{
				ImageMetadataFetcher: fetcher.Spy,
			}
		})

		It("should create the correct docker image ref", func() {
			err := stager.Stage("", cf.StagingRequest{
				Lifecycle: cf.StagingLifecycle{
					DockerLifecycle: &cf.StagingDockerLifecycle{
						Image: "eirini/some-app",
					},
				},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(fetcher.CallCount()).To(Equal(1))
			ref, _ := fetcher.ArgsForCall(0)
			Expect(ref).To(Equal("//docker.io/eirini/some-app"))
		})

	})

})
