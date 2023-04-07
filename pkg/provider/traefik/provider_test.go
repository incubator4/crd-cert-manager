package traefik

import (
	"auto-cert/pkg/event"
	"auto-cert/pkg/provider"
	"context"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"testing"
	"time"
)

func TestTraefikProvider(t *testing.T) {
	p := Provider{Resync: 10 * time.Minute}

	testCases := []struct {
		Ing  v1alpha1.IngressRoute
		Cert v1.Certificate
	}{
		{
			Ing: v1alpha1.IngressRoute{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "traefik.containo.us/v1alpha1",
					Kind:       "IngressRoute",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Annotations: map[string]string{
						provider.IngressClusterIssuerNameAnnotationKey: "test-issuer",
					},
				},
				Spec: v1alpha1.IngressRouteSpec{
					Routes: []v1alpha1.Route{
						{
							Kind:  "Rule",
							Match: "Host(`test.domain.xyz`)",
							Services: []v1alpha1.Service{
								{LoadBalancerSpec: v1alpha1.LoadBalancerSpec{
									Name: "test-svc",
									Port: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
								}},
							},
						},
					},
					TLS: &v1alpha1.TLS{SecretName: "test-cert"},
				},
			},
			Cert: v1.Certificate{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Certificate",
					APIVersion: "cert-manager.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cert",
				},
				Spec: v1.CertificateSpec{
					DNSNames: []string{"test.domain.xyz"},
					IssuerRef: cmmeta.ObjectReference{
						Name: "test-issuer",
						Kind: "ClusterIssuer",
					},
					SecretName: "test-cert",
				},
				Status: v1.CertificateStatus{},
			},
		},
	}

	for _, testCase := range testCases {
		e := event.Event{
			Event: watch.Event{
				Type:   watch.Added,
				Object: &testCase.Ing,
			},
			ProviderName: "traefik",
		}

		cert, err := p.Provide(context.Background(), e)

		assert.Nil(t, err)
		assert.Equal(t, *cert, testCase.Cert)

	}
}
