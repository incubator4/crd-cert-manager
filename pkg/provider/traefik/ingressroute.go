package traefik

import (
	"auto-cert/pkg/provider"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"regexp"
)

type IngressRoute struct {
	provider.DefaultCert
	v1alpha1.IngressRoute
}

func (i *IngressRoute) GetName() string {
	return i.IngressRoute.Name
}

func (i *IngressRoute) GetSecretName() string {
	tls := i.Spec.TLS
	if tls != nil {
		return tls.SecretName
	}
	return ""
}

func (i *IngressRoute) GetHosts() []string {
	hostMap := make(map[string]bool)
	re := regexp.MustCompile(`Host\(([^\)]+)\)`)

	for _, route := range i.Spec.Routes {
		matches := re.FindAllStringSubmatch(route.Match, -1)
		for _, match := range matches {
			hostMap[match[1]] = true
		}
	}

	var hosts []string
	for host := range hostMap {
		hosts = append(hosts, host)
	}

	return hosts
}

var _ provider.Cert = (*IngressRoute)(nil)
