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

var _ = Describe("CCFailedStagingReporter", func() {

	var (
		reporter StagingReporter
		server   *ghttp.Server
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		/*
			handlers := []http.HandlerFunc{
				ghttp.VerifyJSON(`{
					"task_guid": "the-stage-guid",
					"failed": true,
					"failure_reason": "fix this to be more descriptive",
					"result": "",
					"annotation": "{\"completion_callback\": \"internal_cc_staging_endpoint.io/stage/build_completed\"}"
				}`),
			}
		*/

		//		server.RouteToHandler("PUT", "/stage/the-stage-guid/completed", ghttp.CombineHandlers(handlers...))

		reporter = staging.CCFailedStagingReporter{
			Client: http.Client{},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Reporting status to Eirini", func() {

		var (
			thePod *v1.Pod
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
				},
			}

		})

		It("should report the correct failure to CC", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/stage/the-stage-guid/completed"),
					ghttp.VerifyJSON(`{
  					"task_guid": "the-stage-guid",
  					"failed": true,
  					"failure_reason": "ErrImagePull",
  					"result": "",
  					"annotation": "{\"lifecycle\":\"\",\"completion_callback\":\"internal_cc_staging_endpoint.io/stage/build_completed\"}",
  					"created_at": 0
  				}`),
				))

			thePod.Status = v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					{
						State: v1.ContainerState{
							Waiting: &v1.ContainerStateWaiting{
								Reason: "ErrImagePull",
							},
						},
					},
				},
			}
			reporter.Report(thePod)
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("should silently ignore happy Pods", func() {
			thePod.Status = v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					{
						State: v1.ContainerState{
							Waiting: &v1.ContainerStateWaiting{
								Reason: "PodInitializing",
							},
						},
					},
				},
			}
			reporter.Report(thePod)
			Expect(server.ReceivedRequests()).To(HaveLen(0))
		})

		/*
			XContext("When pod init container cannot start", func() {
			})
		*/
	})

})
