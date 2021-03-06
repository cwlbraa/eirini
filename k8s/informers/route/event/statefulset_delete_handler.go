package event

import (
	"code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/route"
	eiriniroute "code.cloudfoundry.org/eirini/route"
	"code.cloudfoundry.org/lager"
	set "github.com/deckarep/golang-set"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type StatefulSetDeleteHandler struct {
	Pods         typedv1.PodInterface
	Logger       lager.Logger
	RouteEmitter eiriniroute.Emitter
}

func (h StatefulSetDeleteHandler) Handle(deletedStatefulSet *appsv1.StatefulSet) {
	loggerSession := h.Logger.Session("statefulset-delete", lager.Data{"guid": deletedStatefulSet.Annotations[k8s.AnnotationProcessGUID]})

	routeSet, err := decodeRoutesAsSet(deletedStatefulSet)
	if err != nil {
		loggerSession.Error("failed-to-decode-deleted-user-defined-routes", err)
		return
	}

	routeGroups := groupRoutesByPort(routeSet, set.NewSet())
	routes := h.createRoutesOnDelete(
		loggerSession,
		deletedStatefulSet,
		routeGroups,
	)
	for _, route := range routes {
		h.RouteEmitter.Emit(*route)
	}
}

func (h StatefulSetDeleteHandler) createRoutesOnDelete(loggerSession lager.Logger, statefulset *appsv1.StatefulSet, grouped portGroup) []*route.Message {
	pods, err := getChildrenPods(h.Pods, statefulset)
	if err != nil {
		loggerSession.Error("failed-to-get-child-pods", err)
		return []*route.Message{}
	}

	resultRoutes := []*route.Message{}
	for _, pod := range pods {
		resultRoutes = append(resultRoutes, createRouteMessages(loggerSession, pod, grouped)...)
	}
	return resultRoutes
}

func getChildrenPods(podClient typedv1.PodInterface, st *appsv1.StatefulSet) ([]corev1.Pod, error) {
	set := labels.Set(st.Spec.Selector.MatchLabels)
	opts := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	podlist, err := podClient.List(opts)
	if err != nil {
		return []corev1.Pod{}, err
	}
	return podlist.Items, nil
}
