package e2e_test

import (
	"context"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/e2e-framework/klient/conf"
)

var (
	k8sClient client.Client
	ctx       context.Context
	scheme    = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	log.SetLogger(GinkgoLogr)
	SetDefaultEventuallyTimeout(3 * time.Minute)
	RunSpecs(t, "Policy Controller E2E Suite")

	format.MaxLength = 0
}

var _ = SynchronizedBeforeSuite(func() []byte {
	kubeconfig := conf.ResolveKubeConfigFile()
	data, err := os.ReadFile(kubeconfig)
	Expect(err).NotTo(HaveOccurred())
	return data
}, func(data []byte) {
	restCfg, err := clientcmd.RESTConfigFromKubeConfig(data)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(restCfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred())

	ctx = context.Background()
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	k8sClient = nil
})
