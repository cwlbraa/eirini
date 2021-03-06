package route

import (
	"errors"

	"code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/route"
	eiriniroute "code.cloudfoundry.org/eirini/route"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

//go:generate counterfeiter . StatefulSetUpdateEventHandler
type StatefulSetUpdateEventHandler interface {
	Handle(oldObj, updatedObj *appsv1.StatefulSet)
}

//go:generate counterfeiter . StatefulSetDeleteEventHandler
type StatefulSetDeleteEventHandler interface {
	Handle(obj *appsv1.StatefulSet)
}

type URIChangeInformer struct {
	Cancel        <-chan struct{}
	Client        kubernetes.Interface
	UpdateHandler StatefulSetUpdateEventHandler
	DeleteHandler StatefulSetDeleteEventHandler
	Namespace     string
}

func NewURIChangeInformer(client kubernetes.Interface, namespace string, updateEventHandler StatefulSetUpdateEventHandler, deleteEventHandler StatefulSetDeleteEventHandler) route.Informer {
	return &URIChangeInformer{
		Client:        client,
		Namespace:     namespace,
		Cancel:        make(<-chan struct{}),
		UpdateHandler: updateEventHandler,
		DeleteHandler: deleteEventHandler,
	}
}

func (i *URIChangeInformer) Start() {
	factory := informers.NewSharedInformerFactoryWithOptions(i.Client,
		NoResync,
		informers.WithNamespace(i.Namespace))

	informer := factory.Apps().V1().StatefulSets().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, updatedObj interface{}) {
			oldStatefulSet := oldObj.(*appsv1.StatefulSet)
			updatedStatefulSet := updatedObj.(*appsv1.StatefulSet)
			i.UpdateHandler.Handle(oldStatefulSet, updatedStatefulSet)
		},
		DeleteFunc: func(obj interface{}) {
			statefulSet := obj.(*appsv1.StatefulSet)
			i.DeleteHandler.Handle(statefulSet)
		},
	})

	informer.Run(i.Cancel)
}

func NewRouteMessage(pod *corev1.Pod, port uint32, routes eiriniroute.Routes) (*eiriniroute.Message, error) {
	if len(pod.Status.PodIP) == 0 {
		return nil, errors.New("missing ip address")
	}

	message := &eiriniroute.Message{
		Routes: eiriniroute.Routes{
			UnregisteredRoutes: routes.UnregisteredRoutes,
		},
		Name:       pod.Labels[k8s.LabelGUID],
		InstanceID: pod.Name,
		Address:    pod.Status.PodIP,
		Port:       port,
		TLSPort:    0,
	}
	if isReady(pod.Status.Conditions) {
		message.RegisteredRoutes = routes.RegisteredRoutes
	}

	if len(message.RegisteredRoutes) == 0 && len(message.UnregisteredRoutes) == 0 {
		return nil, errors.New("no-routes-provided")
	}

	return message, nil
}

func isReady(conditions []v1.PodCondition) bool {
	for _, c := range conditions {
		if c.Type == v1.PodReady {
			return c.Status == v1.ConditionTrue
		}
	}
	return false
}
