package docker

import (
	"fmt"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/eirini/models/cf"
	"github.com/containers/image/types"
	"github.com/docker/distribution/reference"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

//go:generate counterfeiter . ImageMetadataFetcher
type ImageMetadataFetcher func(string, types.SystemContext) (*v1.ImageConfig, error)

func (f ImageMetadataFetcher) Fetch(dockerRef string, sysCtx types.SystemContext) (*v1.ImageConfig, error) {
	return f(dockerRef, sysCtx)
}

type Stager struct {
	ImageMetadataFetcher ImageMetadataFetcher
}

func (s Stager) Stage(stagingGUID string, request cf.StagingRequest) error {
	named, err := reference.ParseNormalizedNamed(request.Lifecycle.DockerLifecycle.Image)
	if err != nil {
		return err
	}
	dockerRef := fmt.Sprintf("//%s", named.Name())
	s.ImageMetadataFetcher.Fetch(dockerRef, types.SystemContext{})

	// call fetchmetadata(d0ckerref)
	// call CompleteStaging
	return nil
}
func (s Stager) CompleteStaging(*models.TaskCallbackResponse) error {
	// create CC response
	// call CC
	return nil
}
