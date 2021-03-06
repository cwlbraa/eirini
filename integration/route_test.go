package statefulsets_test

import (
	"encoding/json"
	"sync"

	"code.cloudfoundry.org/eirini/k8s"
	informerroute "code.cloudfoundry.org/eirini/k8s/informers/route"
	"code.cloudfoundry.org/eirini/k8s/informers/route/event"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
	"code.cloudfoundry.org/eirini/route"
	"code.cloudfoundry.org/eirini/route/routefakes"
	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

var _ = Describe("Routes", func() {

	var (
		desirer opi.Desirer
		odinLRP *opi.LRP
	)

	AfterEach(func() {
		cleanupStatefulSet(odinLRP)
		Eventually(func() []appsv1.StatefulSet {
			return listAllStatefulSets(odinLRP, odinLRP)
		}, timeout).Should(BeEmpty())
	})

	BeforeEach(func() {
		odinLRP = createLRP("ödin")
		logger := lagertest.NewTestLogger("test")
		desirer = k8s.NewStatefulSetDesirer(
			clientset,
			namespace,
			"registry-secret",
			"rootfsversion",
			logger,
		)
	})

	Context("RouteCollector", func() {
		var collector k8s.RouteCollector

		BeforeEach(func() {
			logger := lagertest.NewTestLogger("test")
			collector = k8s.NewRouteCollector(clientset, namespace, logger)
		})

		When("an LRP is desired", func() {
			It("sends register routes message", func() {
				Expect(desirer.Desire(odinLRP)).To(Succeed())
				Eventually(func() bool {
					pods := listPods(odinLRP.LRPIdentifier)
					if len(pods) < 1 {
						return false
					}
					return podReady(pods[0])
				}, timeout).Should(BeTrue())

				routes, err := collector.Collect()
				Expect(err).ToNot(HaveOccurred())
				pods := listPods(odinLRP.LRPIdentifier)
				Expect(routes).To(ContainElement(route.Message{
					InstanceID: pods[0].Name,
					Name:       odinLRP.GUID,
					Address:    pods[0].Status.PodIP,
					Port:       8080,
					TLSPort:    0,
					Routes: route.Routes{
						RegisteredRoutes: []string{"foo.example.com"},
					},
				}))
			})

			When("one of the instances is failing", func() {
				BeforeEach(func() {
					odinLRP = createLRP("odin")
					odinLRP.Health = opi.Healtcheck{
						Type: "port",
						Port: 3000,
					}
					odinLRP.Command = []string{
						"/bin/sh",
						"-c",
						`if [ $(echo $HOSTNAME | sed 's|.*-\(.*\)|\1|') -eq 0 ]; then
	exit;
else
	while true; do
		nc -lk -p 3000 -e echo just a server;
	done;
fi;`,
					}
					err := desirer.Desire(odinLRP)
					Expect(err).ToNot(HaveOccurred())
					Eventually(func() bool {
						pods := listPods(odinLRP.LRPIdentifier)
						if len(pods) < 2 {
							return false
						}
						return podCrashed(pods[0]) && podReady(pods[1])
					}, timeout).Should(BeTrue())
				})

				It("should only return a register message for the working instance", func() {
					routes, err := collector.Collect()
					Expect(err).ToNot(HaveOccurred())
					pods := listPods(odinLRP.LRPIdentifier)
					Expect(routes).To(ContainElement(route.Message{
						InstanceID: pods[1].Name,
						Name:       odinLRP.GUID,
						Address:    pods[1].Status.PodIP,
						Port:       8080,
						TLSPort:    0,
						Routes: route.Routes{
							RegisteredRoutes: []string{"foo.example.com"},
						},
					}))
				})
			})
		})
	})

	Context("InstanceInformer", func() {
		var (
			stopChan         chan struct{}
			informerWG       sync.WaitGroup
			fakeRouteEmitter *routefakes.FakeEmitter
		)

		BeforeEach(func() {
			err := desirer.Desire(odinLRP)
			Expect(err).ToNot(HaveOccurred())
			Eventually(func() bool {
				pods := listPods(odinLRP.LRPIdentifier)
				if len(pods) < 2 {
					return false
				}
				return podReady(pods[0]) && podReady(pods[1])
			}, timeout).Should(BeTrue())

			stopChan = make(chan struct{})
			informerWG = sync.WaitGroup{}
			informerWG.Add(1)

			fakeRouteEmitter = new(routefakes.FakeEmitter)
			logger := lagertest.NewTestLogger("instance-informer-test")
			updateEventHandler := event.URIAnnotationUpdateHandler{
				Pods:         clientset.CoreV1().Pods(namespace),
				Logger:       logger,
				RouteEmitter: fakeRouteEmitter,
			}
			deleteEventHandler := event.StatefulSetDeleteHandler{
				Pods:         clientset.CoreV1().Pods(namespace),
				Logger:       logger,
				RouteEmitter: fakeRouteEmitter,
			}
			informer := &informerroute.URIChangeInformer{
				Client:        clientset,
				Cancel:        stopChan,
				Namespace:     namespace,
				UpdateHandler: updateEventHandler,
				DeleteHandler: deleteEventHandler,
			}

			go func() {
				informer.Start()
				informerWG.Done()
			}()
		})

		AfterEach(func() {
			close(stopChan)
			informerWG.Wait()
		})

		When("an app is stopped", func() {
			It("sends unregister routes message", func() {
				Expect(desirer.Stop(odinLRP.LRPIdentifier)).To(Succeed())
				pods := listPods(odinLRP.LRPIdentifier)

				Eventually(fakeRouteEmitter.EmitCallCount).Should(Equal(2))

				allArgs := []route.Message{
					fakeRouteEmitter.EmitArgsForCall(0),
					fakeRouteEmitter.EmitArgsForCall(1),
				}
				Expect(allArgs).To(ContainElement(route.Message{
					Routes: route.Routes{
						UnregisteredRoutes: []string{"foo.example.com"},
					},
					InstanceID: pods[0].Name,
					Name:       odinLRP.GUID,
					Address:    pods[0].Status.PodIP,
					Port:       8080,
					TLSPort:    0,
				}))
				Expect(allArgs).To(ContainElement(route.Message{
					Routes: route.Routes{
						UnregisteredRoutes: []string{"foo.example.com"},
					},
					InstanceID: pods[1].Name,
					Name:       odinLRP.GUID,
					Address:    pods[1].Status.PodIP,
					Port:       8080,
					TLSPort:    0,
				}))
			})
		})

		When("a new route is registered for an app", func() {
			It("should send a register route message with the new route", func() {
				routes, err := json.Marshal([]cf.Route{
					{Hostname: "foo.example.com", Port: 8080},
					{Hostname: "bar.example.com", Port: 9090},
				})
				Expect(err).ToNot(HaveOccurred())
				odinLRP.AppURIs = string(routes)
				Expect(desirer.Update(odinLRP)).To(Succeed())
				pods := listPods(odinLRP.LRPIdentifier)

				Eventually(fakeRouteEmitter.EmitCallCount).Should(Equal(4))
				allArgs := []route.Message{
					fakeRouteEmitter.EmitArgsForCall(0),
					fakeRouteEmitter.EmitArgsForCall(1),
					fakeRouteEmitter.EmitArgsForCall(2),
					fakeRouteEmitter.EmitArgsForCall(3),
				}
				Expect(allArgs).To(ContainElement(route.Message{
					Routes: route.Routes{
						RegisteredRoutes: []string{"bar.example.com"},
					},
					InstanceID: pods[0].Name,
					Name:       odinLRP.GUID,
					Address:    pods[0].Status.PodIP,
					Port:       9090,
					TLSPort:    0,
				}))
			})
		})
	})

	Context("URIChangeInformer", func() {
		var (
			stopChan         chan struct{}
			informerWG       sync.WaitGroup
			fakeRouteEmitter *routefakes.FakeEmitter
		)

		BeforeEach(func() {
			odinLRP.TargetInstances = 2
			Expect(desirer.Desire(odinLRP)).To(Succeed())
			Eventually(func() bool {
				pods := listPods(odinLRP.LRPIdentifier)
				if len(pods) < 2 {
					return false
				}
				return podReady(pods[0]) && podReady(pods[1])
			}, timeout).Should(BeTrue())

			stopChan = make(chan struct{})
			informerWG = sync.WaitGroup{}
			informerWG.Add(1)

			logger := lagertest.NewTestLogger("instance-informer-test")
			fakeRouteEmitter = new(routefakes.FakeEmitter)
			updateEventHandler := event.PodUpdateHandler{
				Client:       clientset.AppsV1().StatefulSets(namespace),
				Logger:       logger,
				RouteEmitter: fakeRouteEmitter,
			}
			informer := &informerroute.InstanceChangeInformer{
				Client:        clientset,
				Cancel:        stopChan,
				Namespace:     namespace,
				UpdateHandler: updateEventHandler,
			}
			go func() {
				informer.Start()
				informerWG.Done()
			}()
		})

		AfterEach(func() {
			close(stopChan)
			informerWG.Wait()
		})

		When("the app is scaled down", func() {
			It("sends unregister routes message", func() {
				pods := listPods(odinLRP.LRPIdentifier)
				odinLRP.TargetInstances = 1
				Expect(desirer.Update(odinLRP)).To(Succeed())

				Eventually(fakeRouteEmitter.EmitCallCount, timeout).Should(Equal(1))
				Expect(fakeRouteEmitter.EmitArgsForCall(0)).To(Equal(route.Message{
					Routes: route.Routes{
						UnregisteredRoutes: []string{"foo.example.com"},
					},
					InstanceID: pods[1].Name,
					Name:       odinLRP.GUID,
					Address:    pods[1].Status.PodIP,
					Port:       8080,
					TLSPort:    0,
				}))
			})
		})

		When("an app instance is stopped", func() {
			It("sends unregister routes message", func() {
				pods := listPods(odinLRP.LRPIdentifier)
				Expect(desirer.StopInstance(odinLRP.LRPIdentifier, 0)).To(Succeed())

				Eventually(fakeRouteEmitter.EmitCallCount, timeout).Should(Equal(1))
				Expect(fakeRouteEmitter.EmitArgsForCall(0)).To(Equal(route.Message{
					Routes: route.Routes{
						UnregisteredRoutes: []string{"foo.example.com"},
					},
					InstanceID: pods[0].Name,
					Name:       odinLRP.GUID,
					Address:    pods[0].Status.PodIP,
					Port:       8080,
					TLSPort:    0,
				}))
			})
		})
	})
})
