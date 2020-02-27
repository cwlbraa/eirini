package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/eirini"
	cmdcommons "code.cloudfoundry.org/eirini/cmd"
	"code.cloudfoundry.org/eirini/k8s/informers/staging"
	"code.cloudfoundry.org/lager"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
)

type options struct {
	ConfigFile string `short:"c" long:"config" description:"Config for running event-reporter"`
}

func main() {
	//var opts options
	//_, err := flags.ParseArgs(&opts, os.Args)
	//cmdcommons.ExitIfError(err)

	// cfg, err := readConfigFile(opts.ConfigFile)
	// cmdcommons.ExitIfError(err)

	clientset := cmdcommons.CreateKubeClient("")

	launchStagingReporter(
		clientset,
		"",
		"",
		"",
		"",
		"eirini",
	)
}

func launchStagingReporter(clientset kubernetes.Interface, uri, ca, cert, key, namespace string) {
	stagingLogger := lager.NewLogger("staging-informer")
	stagingLogger.RegisterSink(lager.NewPrettySink(os.Stdout, lager.DEBUG))
	stagingInformer := staging.NewInformer(clientset, 0, namespace, make(chan struct{}), stagingLogger)

	stagingInformer.Start()
}

func readConfigFile(path string) (*eirini.EventReporterConfig, error) {
	fileBytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	var conf eirini.EventReporterConfig
	err = yaml.Unmarshal(fileBytes, &conf)
	return &conf, errors.Wrap(err, "failed to unmarshal yaml")
}
