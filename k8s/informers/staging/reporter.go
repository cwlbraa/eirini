package staging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/runtimeschema/cc_messages"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

const PodInitializing = "PodInitializing"

type CCFailedStagingReporter struct {
	Client http.Client
}

func (r CCFailedStagingReporter) Report(pod *v1.Pod) error {
	stagingGUID := pod.Labels[k8s.LabelStagingGUID]
	status := getFailedContainerStatusIfAny(pod.Status.ContainerStatuses)
	if status == nil {
		return nil
	}
	reason := status.State.Waiting.Reason
	completionCallback, _ := getEnvVarValue("COMPLETION_CALLBACK", pod.Spec.Containers[0].Env)
	eiriniAddr, _ := getEnvVarValue("EIRINI_ADDRESS", pod.Spec.Containers[0].Env)
	response := r.createFailureResponse(reason, stagingGUID, completionCallback)
	r.sendResponse(eiriniAddr, response)
	return nil
}

func getEnvVarValue(key string, vars []v1.EnvVar) (string, error) {
	for _, envVar := range vars {
		if envVar.Name == key {
			return envVar.Value, nil
		}
	}
	return "", errors.New("failed to find env var")
}

func getFailedContainerStatusIfAny(statuses []v1.ContainerStatus) *v1.ContainerStatus {
	for _, status := range statuses {
		waiting := status.State.Waiting
		if waiting != nil && waiting.Reason != PodInitializing {
			return &status
		}
	}
	return nil
}

func (r CCFailedStagingReporter) createFailureResponse(failure string, stagingGUID, completionCallback string) *models.TaskCallbackResponse {
	annotation := cc_messages.StagingTaskAnnotation{
		CompletionCallback: completionCallback,
	}

	annotationJSON, err := json.Marshal(annotation)
	if err != nil {
		panic(err)
	}

	return &models.TaskCallbackResponse{
		TaskGuid:      stagingGUID,
		Failed:        true,
		FailureReason: failure,
		Annotation:    string(annotationJSON),
	}
}

func (r CCFailedStagingReporter) sendResponse(eiriniAddr string, response *models.TaskCallbackResponse) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	uri := fmt.Sprintf("%s/stage/%s/completed", eiriniAddr, response.TaskGuid)

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(responseJSON))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		var message string
		if err == nil {
			message = string(body)
		}
		return fmt.Errorf("request not successful: status=%d taskGuid=%s %s", resp.StatusCode, response.TaskGuid, message)
	}

	return nil
}
