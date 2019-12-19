package crash_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/eirini/integration/util"
	"code.cloudfoundry.org/eirini/k8s"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crashes", func() {

	var (
		desirer    opi.Desirer
		odinLRP    *opi.LRP
		configFile *os.File
		capiServer *ghttp.Server
	)

	BeforeEach(func() {
		var err error
		capiServer = createTestServer(
			util.PathToTestFixture("cert"),
			util.PathToTestFixture("cert"),
			util.PathToTestFixture("key"),
		)
		config := defaultEventReporterConfig()
		config.CcInternalAPI = capiServer.URL()
		configFile, err = createEventReporterConfigFile(config)
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command(pathToCrashEmitter, "-c", configFile.Name()) // #nosec G204
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		desirer = k8s.NewStatefulSetDesirer(
			clientset,
			namespace,
			"registry-secret",
			"rootfsversion",
			logger,
		)
	})

	AfterEach(func() {
		capiServer.Stop()
	})

	Context("When an app crashes", func() {

		BeforeEach(func() {
			odinLRP = createCrashingLRP("ödin")
			Expect(desirer.Desire(odinLRP)).To(Succeed())
		})

		It("generates crash report for the initial error", func() {
			// "ProcessGUID": Equal(fmt.Sprintf("%s-%s", odinLRP.GUID, odinLRP.Version)),
			// "AppCrashedRequest": MatchFields(IgnoreExtras, Fields{
			// 	"Instance":        ContainSubstring(odinLRP.GUID),
			// 	"Index":           Equal(0),
			// 	"Reason":          Equal("Error"),
			// 	"ExitStatus":      Equal(1),
			// 	"ExitDescription": Equal("Error"),
			// }),
		})

		It("generates crash report when the app goes into CrashLoopBackOff", func() {
			// "ProcessGUID": Equal(fmt.Sprintf("%s-%s", odinLRP.GUID, odinLRP.Version)),
			// "AppCrashedRequest": MatchFields(IgnoreExtras, Fields{
			// 	"Instance":        ContainSubstring(odinLRP.GUID),
			// 	"Index":           Equal(0),
			// 	"Reason":          Equal("CrashLoopBackOff"),
			// 	"ExitStatus":      Equal(1),
			// 	"ExitDescription": Equal("Error"),
			// }),
		})
	})
})

func createCrashingLRP(name string) *opi.LRP {
	guid := util.RandomString()
	routes, err := json.Marshal([]cf.Route{{Hostname: "foo.example.com", Port: 8080}})
	Expect(err).ToNot(HaveOccurred())
	return &opi.LRP{
		Command: []string{
			"/bin/sh",
			"-c",
			"exit 1",
		},
		AppName:         name,
		SpaceName:       "space-foo",
		TargetInstances: 1,
		Image:           "alpine",
		AppURIs:         string(routes),
		LRPIdentifier:   opi.LRPIdentifier{GUID: guid, Version: "version_" + guid},
		LRP:             "metadata",
		DiskMB:          2047,
	}
}

func createTestServer(certName, keyName, caCertName string) *ghttp.Server {
	certPath := filepath.Join(certsPath, certName)
	keyPath := filepath.Join(certsPath, keyName)
	caCertPath := filepath.Join(certsPath, caCertName)

	tlsConf, tlsErr := tlsconfig.Build(
		tlsconfig.WithInternalServiceDefaults(),
		tlsconfig.WithIdentityFromFile(certPath, keyPath),
	).Server(
		tlsconfig.WithClientAuthenticationFromFile(caCertPath),
	)
	Expect(tlsErr).NotTo(HaveOccurred())

	testServer := ghttp.NewUnstartedServer()
	testServer.HTTPTestServer.TLS = tlsConf

	return testServer
}
