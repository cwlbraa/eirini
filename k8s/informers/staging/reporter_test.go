package staging_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/k8s/informers/staging"
	. "code.cloudfoundry.org/eirini/k8s/informers/staging"
)

var _ = Describe("FailedStagingReporter", func() {

	var (
		reporter StagingReporter
		server   *ghttp.Server
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		reporter = staging.FailedStagingReporter{
			Client: &http.Client{},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Reporting status to Eirini", func() {

		var (
			thePod                 *v1.Pod
			statusFailed, statusOK v1.ContainerStatus
		)

		BeforeEach(func() {
			thePod = &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "i'm not feeling well",
					Annotations: map[string]string{},
					Labels: map[string]string{
						k8s.LabelStagingGUID: "the-stage-guid",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Env: []v1.EnvVar{
								{
									Name:  "COMPLETION_CALLBACK",
									Value: "internal_cc_staging_endpoint.io/stage/build_completed",
								},
								{
									Name:  "EIRINI_ADDRESS",
									Value: server.URL(),
								},
							},
						},
					},
					InitContainers: []v1.Container{},
				},
			}

			statusFailed = v1.ContainerStatus{
				Name: "failing-container",
				State: v1.ContainerState{
					Waiting: &v1.ContainerStateWaiting{
						Reason: "ErrImagePull",
					},
				},
			}

			statusOK = v1.ContainerStatus{
				Name: "starting-container",
				State: v1.ContainerState{
					Waiting: &v1.ContainerStateWaiting{
						Reason: "PodInitializing",
					},
				},
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/stage/the-stage-guid/completed"),
					ghttp.VerifyJSON(`{
  					"task_guid": "the-stage-guid",
  					"failed": true,
						"failure_reason": "Container failing-container failed: ErrImagePull",
  					"result": "",
  					"annotation": "{\"lifecycle\":\"\",\"completion_callback\":\"internal_cc_staging_endpoint.io/stage/build_completed\"}",
  					"created_at": 0
  				}`),
				))

		})

		It("should report the correct container failure reason to Eirini", func() {
			thePod.Status = v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					statusFailed,
				},
			}
			reporter.Report(thePod)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("should silently ignore happy containers", func() {
			thePod.Status = v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					statusOK,
				},
			}
			reporter.Report(thePod)
			Expect(server.ReceivedRequests()).To(HaveLen(0))
		})

		It("should detect failing InitContainers", func() {
			thePod.Status = v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					statusOK,
				},
				InitContainerStatuses: []v1.ContainerStatus{
					statusFailed,
				},
			}
			reporter.Report(thePod)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		/*
			XContext("When pod init container cannot start", func() {
			})
		*/
	})

})
