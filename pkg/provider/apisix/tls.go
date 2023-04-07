package apisix

import (
	"auto-cert/pkg/provider"
	"github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2beta3"
)

type ApisixTls struct {
	provider.DefaultCert
	v2beta3.ApisixTls
}
