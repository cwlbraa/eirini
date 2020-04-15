package util

import (
	"bufio"
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/gomega" //nolint:golint,stylecheck
)

type Fixture struct {
	Clientset      kubernetes.Interface
	Namespace      string
	PspName        string
	KubeConfigPath string
	Writer         io.Writer
}

func NewFixture(writer io.Writer) *Fixture {
	kubeConfigPath := GetKubeconfig()

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	Expect(err).NotTo(HaveOccurred(), "failed to build config from flags")

	clientset, err := kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred(), "failed to create clientset")

	return &Fixture{
		KubeConfigPath: kubeConfigPath,
		Clientset:      clientset,
		Writer:         writer,
	}
}

func (f *Fixture) SetUp() {
	f.Namespace = CreateRandomNamespace(f.Clientset)
	f.PspName = fmt.Sprintf("%s-psp", f.Namespace)
	Expect(CreatePodCreationPSP(f.Namespace, f.PspName, f.Clientset)).To(Succeed(), "failed to create pod creation PSP")
}

func (f *Fixture) TearDown() {
	var errs *multierror.Error
	errs = multierror.Append(errs, f.printDebugInfo())

	errs = multierror.Append(errs, DeleteNamespace(f.Namespace, f.Clientset))
	errs = multierror.Append(errs, DeletePSP(f.PspName, f.Clientset))

	Expect(errs.ErrorOrNil()).NotTo(HaveOccurred())
}

//nolint:gocyclo
func (f *Fixture) printDebugInfo() error {
	if _, err := f.Writer.Write([]byte("Jobs:\n")); err != nil {
		return err
	}
	jobs, _ := f.Clientset.BatchV1().Jobs(f.Namespace).List(v1.ListOptions{})
	for _, job := range jobs.Items {
		fmt.Fprintf(f.Writer, "Job: %s status is: %#v\n", job.Name, job.Status)
		if _, err := f.Writer.Write([]byte("-----------\n")); err != nil {
			return err
		}
	}

	statefulsets, _ := f.Clientset.AppsV1().StatefulSets(f.Namespace).List(v1.ListOptions{})
	if _, err := f.Writer.Write([]byte("StatefulSets:\n")); err != nil {
		return err
	}
	for _, s := range statefulsets.Items {
		fmt.Fprintf(f.Writer, "StatefulSet: %s status is: %#v\n", s.Name, s.Status)
		if _, err := f.Writer.Write([]byte("-----------\n")); err != nil {
			return err
		}
	}

	pods, _ := f.Clientset.CoreV1().Pods(f.Namespace).List(v1.ListOptions{})
	if _, err := f.Writer.Write([]byte("Pods:\n")); err != nil {
		return err
	}
	for _, p := range pods.Items {
		fmt.Fprintf(f.Writer, "Pod: %s status is: %#v\n", p.Name, p.Status)
		if _, err := f.Writer.Write([]byte("-----------\n")); err != nil {
			return err
		}
		fmt.Fprintf(f.Writer, "Pod: %s logs are: \n", p.Name)
		logsReq := f.Clientset.CoreV1().Pods(f.Namespace).GetLogs(p.Name, &corev1.PodLogOptions{})
		if err := consumeRequest(logsReq, f.Writer); err != nil {
			fmt.Fprintf(f.Writer, "Failed to get logs for Pod: %s becase: %v \n", p.Name, err)
		}
	}

	return nil
}

func consumeRequest(request rest.ResponseWrapper, out io.Writer) error {
	readCloser, err := request.Stream()
	if err != nil {
		return err
	}
	defer readCloser.Close()

	r := bufio.NewReader(readCloser)
	for {
		bytes, err := r.ReadBytes('\n')
		if _, writeErr := out.Write(bytes); writeErr != nil {
			return writeErr
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}