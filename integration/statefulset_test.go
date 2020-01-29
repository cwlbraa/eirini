package statefulsets_test

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/eirini/integration/util"
	. "code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("StatefulSet Manager", func() {

	var (
		desirer opi.Desirer
		odinLRP *opi.LRP
		thorLRP *opi.LRP
	)

	BeforeEach(func() {
		odinLRP = createLRP("ödin")
		thorLRP = createLRP("thor")
	})

	AfterEach(func() {
		cleanupStatefulSet(odinLRP)
		cleanupStatefulSet(thorLRP)
		Eventually(func() []appsv1.StatefulSet {
			return listAllStatefulSets(odinLRP, thorLRP)
		}, timeout).Should(BeEmpty())
	})

	JustBeforeEach(func() {
		logger := lagertest.NewTestLogger("test")
		desirer = NewStatefulSetDesirer(
			clientset,
			namespace,
			"registry-secret",
			"rootfsversion",
			logger,
		)
	})

	Context("When creating a StatefulSet", func() {

		JustBeforeEach(func() {
			err := desirer.Desire(odinLRP)
			Expect(err).ToNot(HaveOccurred())
			err = desirer.Desire(thorLRP)
			Expect(err).ToNot(HaveOccurred())
		})

		// join all tests in a single with By()
		It("should create a StatefulSet object", func() {
			statefulset := getStatefulSet(odinLRP)
			Expect(statefulset.Name).To(ContainSubstring(odinLRP.GUID))
			Expect(statefulset.Spec.Template.Spec.Containers[0].Command).To(Equal(odinLRP.Command))
			Expect(statefulset.Spec.Template.Spec.Containers[0].Image).To(Equal(odinLRP.Image))
			Expect(statefulset.Spec.Replicas).To(Equal(int32ptr(odinLRP.TargetInstances)))
			Expect(statefulset.Annotations[AnnotationOriginalRequest]).To(Equal(odinLRP.LRP))
		})

		It("should create all associated pods", func() {
			var pods []string
			Eventually(func() []string {
				pods = podNamesFromPods(listPods(odinLRP.LRPIdentifier))
				return pods
			}, timeout).Should(HaveLen(2))
			Expect(pods[0]).To(ContainSubstring(odinLRP.GUID))
			Expect(pods[1]).To(ContainSubstring(odinLRP.GUID))
		})

		It("should create a pod disruption budget for the lrp", func() {
			statefulset := getStatefulSet(odinLRP)
			pdb, err := podDisruptionBudgets().Get(statefulset.Name, v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(pdb).NotTo(BeNil())
		})

		Context("when the lrp has 1 instance", func() {
			BeforeEach(func() {
				odinLRP.TargetInstances = 1
			})
			It("should not create a pod disruption budget for the lrp", func() {
				statefulset := getStatefulSet(odinLRP)
				_, err := podDisruptionBudgets().Get(statefulset.Name, v1.GetOptions{})
				Expect(err).To(MatchError(ContainSubstring("not found")))
			})
		})

		Context("when additional app info is provided", func() {
			BeforeEach(func() {
				odinLRP.OrgName = "odin-org"
				odinLRP.OrgGUID = "odin-org-guid"
				odinLRP.SpaceName = "odin-space"
				odinLRP.SpaceGUID = "odin-space-guid"
			})

			DescribeTable("sets appropriate annotations to statefulset", func(key, value string) {
				statefulset := getStatefulSet(odinLRP)
				Expect(statefulset.Annotations).To(HaveKeyWithValue(key, value))
			},
				Entry("SpaceName", AnnotationSpaceName, "odin-space"),
				Entry("SpaceGUID", AnnotationSpaceGUID, "odin-space-guid"),
				Entry("OrgName", AnnotationOrgName, "odin-org"),
				Entry("OrgGUID", AnnotationOrgGUID, "odin-org-guid"),
			)

			It("sets appropriate labels to statefulset", func() {
				statefulset := getStatefulSet(odinLRP)
				Expect(statefulset.Labels).To(HaveKeyWithValue(LabelGUID, odinLRP.LRPIdentifier.GUID))
				Expect(statefulset.Labels).To(HaveKeyWithValue(LabelVersion, odinLRP.LRPIdentifier.Version))
				Expect(statefulset.Labels).To(HaveKeyWithValue(LabelSourceType, "APP"))
			})

		})

		Context("when we create the same StatefulSet again", func() {
			It("should error", func() {
				err := desirer.Desire(odinLRP)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("When deleting a LRP", func() {
		var statefulsetName string

		JustBeforeEach(func() {
			err := desirer.Desire(odinLRP)
			Expect(err).ToNot(HaveOccurred())

			statefulsetName = getStatefulSet(odinLRP).Name

			err = desirer.Stop(odinLRP.LRPIdentifier)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should delete the StatefulSet object", func() {
			Eventually(func() []appsv1.StatefulSet {
				return listStatefulSets("odin")
			}, timeout).Should(BeEmpty())
		})

		It("should delete the associated pods", func() {
			Eventually(func() []corev1.Pod {
				return listPods(odinLRP.LRPIdentifier)
			}, timeout).Should(BeEmpty())
		})

		It("should delete the pod disruption budget for the lrp", func() {
			Eventually(func() error {
				_, err := podDisruptionBudgets().Get(statefulsetName, v1.GetOptions{})
				return err
			}, timeout).Should(MatchError(ContainSubstring("not found")))
		})

		Context("when the lrp has only 1 instance", func() {
			BeforeEach(func() {
				odinLRP.TargetInstances = 1
			})

			It("should delete the pod disruption budget for the lrp", func() {
				Eventually(func() error {
					_, err := podDisruptionBudgets().Get(statefulsetName, v1.GetOptions{})
					return err
				}, timeout).Should(MatchError(ContainSubstring("not found")))
			})
		})
	})

	Context("When updating a LRP", func() {
		var (
			instancesBefore int
			instancesAfter  int
		)

		JustBeforeEach(func() {
			odinLRP.TargetInstances = instancesBefore
			Expect(desirer.Desire(odinLRP)).To(Succeed())

			odinLRP.TargetInstances = instancesAfter
			Expect(desirer.Update(odinLRP)).To(Succeed())
		})

		Context("when scaling up from 1 to 2 instances", func() {
			BeforeEach(func() {
				instancesBefore = 1
				instancesAfter = 2
			})

			It("should create a pod disruption budget for the lrp", func() {
				statefulset := getStatefulSet(odinLRP)
				pdb, err := podDisruptionBudgets().Get(statefulset.Name, v1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(pdb).NotTo(BeNil())
			})
		})

		Context("when scaling up from 2 to 3 instances", func() {
			BeforeEach(func() {
				instancesBefore = 2
				instancesAfter = 3
			})

			It("should keep the existing pod disruption budget for the lrp", func() {
				statefulset := getStatefulSet(odinLRP)
				pdb, err := podDisruptionBudgets().Get(statefulset.Name, v1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(pdb).NotTo(BeNil())
			})
		})

		Context("when scaling down from 2 to 1 instances", func() {
			BeforeEach(func() {
				instancesBefore = 2
				instancesAfter = 1
			})

			It("should delete the pod disruption budget for the lrp", func() {
				statefulset := getStatefulSet(odinLRP)
				_, err := podDisruptionBudgets().Get(statefulset.Name, v1.GetOptions{})
				Expect(err).To(MatchError(ContainSubstring("not found")))
			})
		})

		Context("when scaling down from 1 to 0 instances", func() {
			BeforeEach(func() {
				instancesBefore = 1
				instancesAfter = 0
			})

			It("should keep the lrp without a pod disruption budget", func() {
				statefulset := getStatefulSet(odinLRP)
				_, err := podDisruptionBudgets().Get(statefulset.Name, v1.GetOptions{})
				Expect(err).To(MatchError(ContainSubstring("not found")))
			})
		})
	})

	Context("When getting an app", func() {
		numberOfInstancesFn := func() int {
			lrp, err := desirer.Get(odinLRP.LRPIdentifier)
			Expect(err).ToNot(HaveOccurred())
			return lrp.RunningInstances
		}

		JustBeforeEach(func() {
			err := desirer.Desire(odinLRP)
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly reports the running instances", func() {
			Eventually(numberOfInstancesFn, timeout).Should(Equal(2))
			Consistently(numberOfInstancesFn, "10s").Should(Equal(2))
		})

		Context("When one of the instances if failing", func() {
			BeforeEach(func() {
				odinLRP = createLRP("odin")
				odinLRP.Health = opi.Healtcheck{
					Type: "port",
					Port: 3000,
				}
				odinLRP.Command = []string{
					"/bin/sh",
					"-c",
					`if [ $(echo $HOSTNAME | sed 's|.*-\(.*\)|\1|') -eq 1 ]; then
	exit;
else
	while true; do
		nc -lk -p 3000 -e echo just a server;
	done;
fi;`,
				}
			})

			It("correctly reports the running instances", func() {
				Eventually(numberOfInstancesFn, timeout).Should(Equal(1), fmt.Sprintf("pod %#v did not start", odinLRP.LRPIdentifier))
				Consistently(numberOfInstancesFn, "10s").Should(Equal(1), fmt.Sprintf("pod %#v did not keep running", odinLRP.LRPIdentifier))
			})
		})
	})

})

func int32ptr(i int) *int32 {
	i32 := int32(i)
	return &i32
}

func createLRP(name string) *opi.LRP {
	guid := util.RandomString()
	routes, err := json.Marshal([]cf.Route{{Hostname: "foo.example.com", Port: 8080}})
	Expect(err).ToNot(HaveOccurred())
	return &opi.LRP{
		Command: []string{
			"/bin/sh",
			"-c",
			"while true; do echo hello; sleep 10;done",
		},
		AppName:         name,
		SpaceName:       "space-foo",
		TargetInstances: 2,
		Image:           "busybox",
		AppURIs:         string(routes),
		LRPIdentifier:   opi.LRPIdentifier{GUID: guid, Version: "version_" + guid},
		LRP:             "metadata",
		DiskMB:          2047,
	}
}
