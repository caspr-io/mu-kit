package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type K8sConfig struct {
	Master               string
	CertificateAuthority []byte
	User                 string
	ClientCertificate    []byte
	ClientKey            []byte
}

type ObjectLoader struct {
	kubeConfig *K8sConfig
	Fallback   Loader
}

var _ Loader = (*ObjectLoader)(nil)

func (ol *ObjectLoader) Load() (*rest.Config, error) {
	c := clientcmdapi.NewConfig()

	cluster := clientcmdapi.NewCluster()
	cluster.Server = ol.kubeConfig.Master
	cluster.CertificateAuthorityData = ol.kubeConfig.CertificateAuthority

	auth := clientcmdapi.NewAuthInfo()
	auth.ClientCertificateData = ol.kubeConfig.ClientCertificate
	auth.Username = ol.kubeConfig.User
	auth.ClientKeyData = ol.kubeConfig.ClientKey

	context := clientcmdapi.NewContext()
	context.Cluster = "my-cluster"
	context.AuthInfo = "my-cluster-auth"

	c.AuthInfos["my-cluster-auth"] = auth
	c.Clusters["my-cluster"] = cluster
	c.Contexts["default"] = context
	c.CurrentContext = "default"

	dccfg := clientcmd.NewDefaultClientConfig(*c, nil)

	return dccfg.ClientConfig()
}
