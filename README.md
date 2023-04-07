Crd Cert-Manager
---------
Provide Cert-Manager automatic registration for Ingress Controller CRD

## Concept

Due to some Ingress Controller such as `Traefik` has CRD resources, cert-manager could not easily create cert by annotations.  
This project mimics the configuration of Cert-Manager `Ingress` annotations
to provide automatic discovery support for Ingress Controller CRD

## Ingress Controller Supported
- [x] Traefik (IngressRoute)
- [ ] Apisix Ingress Controller (ApisixTls)
- [ ] Contour (HTTPProxy)

Welcome to provide more Ingress Controller with CRD.

## 