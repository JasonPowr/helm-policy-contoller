package e2e_test

import (
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ns             = "policy-controller-operator"
	bundlePath     = "charts/common_install_cr.yaml"
	deploymentName = "policycontroller-sample-policy-controller-webhook"
)

var (
	absBundle = ""
	err       error
)

var _ = Describe("policy-controller-operator installation", Ordered, func() {

	var err error
	absBundle, err = filepath.Abs(bundlePath)
	Expect(err).ToNot(HaveOccurred())

	BeforeAll(func() {
		By("ensuring the policy-controller-operator namespace exists")
		Expect(k8sClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: ns},
		})).To(SatisfyAny(Succeed(), MatchError(ContainSubstring("already exists"))))

		By("applying the operator bundle: " + absBundle)
		cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", absBundle, "-n", ns)
		Expect(cmd.Run()).To(Succeed())
	})

	AfterAll(func() {
		By("removing the operator bundle: " + absBundle)
		cmd := exec.CommandContext(ctx, "kubectl", "delete", "-f", absBundle, "-n", ns)
		Expect(cmd.Run()).To(Succeed())
	})

	It("creates and becomes ready the policy-controller Deployment", func() {
		dep := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: deploymentName}, dep)
		}).Should(Succeed(), "timed out waiting for Deployment %q to exist", deploymentName)

		desired := *dep.Spec.Replicas
		Eventually(func() int32 {
			k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: deploymentName}, dep)
			Expect(err).ToNot(HaveOccurred())
			return dep.Status.ReadyReplicas
		}).Should(Equal(desired), "timed out waiting for %d pods to be Ready in Deployment %q", desired, deploymentName)
	})

	It("creates the ValidatingWebhookConfiguration", func() {
		validatingWebhookName := "policy.rhtas.com"
		vwc := &admissionregistrationv1.ValidatingWebhookConfiguration{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Name: validatingWebhookName}, vwc)
		}).Should(Succeed(), "timed out waiting for ValidatingWebhookConfiguration %q to exist", validatingWebhookName)
	})

	// MutatingWebhookConfiguration
	It("creates the MutatingWebhookConfiguration", func() {
		mutatingWebhookName := "policy.rhtas.com"
		mwc := &admissionregistrationv1.MutatingWebhookConfiguration{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Name: mutatingWebhookName}, mwc)
		}).Should(Succeed(), "timed out waiting for MutatingWebhookConfiguration %q to exist", mutatingWebhookName)
	})

	It("creates the ValidatingWebhookConfiguration", func() {
		validatingWebhookName := "validating.clusterimagepolicy.rhtas.com"
		vwc := &admissionregistrationv1.ValidatingWebhookConfiguration{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Name: validatingWebhookName}, vwc)
		}).Should(Succeed(), "timed out waiting for ValidatingWebhookConfiguration %q to exist", validatingWebhookName)
	})

	// MutatingWebhookConfiguration
	It("creates the MutatingWebhookConfiguration", func() {
		mutatingWebhookName := "defaulting.clusterimagepolicy.rhtas.com"
		mwc := &admissionregistrationv1.MutatingWebhookConfiguration{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Name: mutatingWebhookName}, mwc)
		}).Should(Succeed(), "timed out waiting for MutatingWebhookConfiguration %q to exist", mutatingWebhookName)
	})

	// webhook Service
	It("creates the webhook Service", func() {
		svcName := "webhook"
		svc := &corev1.Service{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: svcName}, svc)
		}).Should(Succeed(), "timed out waiting for Service %q to exist", svcName)
	})

	// webhook-metrics Service
	It("creates the webhook-metrics Service", func() {
		metricsSvcName := "policycontroller-sample-policy-controller-webhook-metrics"
		ms := &corev1.Service{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: metricsSvcName}, ms)
		}).Should(Succeed(), "timed out waiting for Service %q to exist", metricsSvcName)
	})

	// webhook-certs Secret
	It("creates the webhook-certs Secret", func() {
		secretName := "webhook-certs"
		sec := &corev1.Secret{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: secretName}, sec)
		}).Should(Succeed(), "timed out waiting for Secret %q to exist", secretName)
	})

	// ConfigMaps
	It("creates all required ConfigMaps", func() {
		configMapNames := []string{
			"config-policy-controller",
			"config-image-policies",
			"config-sigstore-keys",
			"policycontroller-sample-policy-controller-webhook-logging",
		}
		for _, cmName := range configMapNames {
			cm := &corev1.ConfigMap{}
			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: cmName}, cm)
			}).Should(Succeed(), "timed out waiting for ConfigMap %q to exist", cmName)
		}
	})
})
