package tls

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Controller manages a tls record.
type Controller struct {
	kubeconfig     string
	namespace      string
	export         bool
	exportName     string
	fromSecretName string
	fromPemFile    string
	fromSecretFile string

	writer io.Writer
}

const (
	defaultNamespace  = "default"
	defaultExportName = "tls"
)

// NewController creates a new controller.
func NewController() *Controller {
	return &Controller{
		writer: os.Stdout,
	}
}

// AddFlags adds flags to command.
func (c *Controller) AddFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&c.kubeconfig, "kubeconfig", "", "Path to kubeconfig.")
	flags.StringVarP(&c.namespace, "namespace", "n", defaultNamespace, "Resource namespace.")
	flags.StringVar(&c.fromSecretName, "from-secret-name", "", "Secret name.")
	flags.StringVar(&c.fromPemFile, "from-pem-file", "", "Path to certificate pem file.")
	flags.StringVar(&c.fromSecretFile, "from-secret-file", "", "Path to secret file.")
	flags.BoolVarP(&c.export, "export", "E", false, "Whether to export certificate and private key.")
	flags.StringVar(&c.exportName, "export-name", defaultExportName, "Export name.")
}

// Run executes command.
func (c *Controller) Run() error {
	var crt, key []byte
	var secret *v1.Secret
	var ok bool
	var err error

	if len(c.fromSecretName) != 0 {
		kubeCfg, err := clientcmd.BuildConfigFromFlags("", c.kubeconfig)
		if err != nil {
			return err
		}
		client, err := kubernetes.NewForConfig(kubeCfg)
		if err != nil {
			return err
		}
		secret, err = client.CoreV1().Secrets(c.namespace).Get(c.fromSecretName, metav1.GetOptions{})
		if err != nil {
			return err
		}
	} else if len(c.fromPemFile) != 0 {
		crt, err = ioutil.ReadFile(c.fromPemFile)
		if err != nil {
			return err
		}
	} else if len(c.fromSecretFile) != 0 {
		reader, err := os.Open(c.fromSecretFile)
		if err != nil {
			return err
		}
		secret = &v1.Secret{}
		if err = yaml.NewYAMLToJSONDecoder(reader).Decode(secret); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("One of --from-secret-name, --from-pem-file, --from-secret-file is required.")
	}

	if secret != nil {
		crt, ok = secret.Data["tls.crt"]
		if !ok {
			return fmt.Errorf("tls.crt not exist")
		}
		key, ok = secret.Data["tls.key"]
		if !ok {
			return fmt.Errorf("tls.key not exist")
		}
	}

	if c.export && key != nil {
		if err = c.exportTls(crt, key); err != nil {
			return err
		}
	}

	pemBlock, _ := pem.Decode(crt)
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.writer, "Subject: %s\n", cert.Subject)
	fmt.Fprintf(c.writer, "Issur CommonName: %s\n", cert.Issuer.CommonName)
	fmt.Fprintf(c.writer, "Subject CommonName: %s\n", cert.Subject.CommonName)
	fmt.Fprintf(c.writer, "DNSNames: %v\n", cert.DNSNames)
	fmt.Fprintf(c.writer, "EmailAddresses: %v\n", cert.EmailAddresses)
	fmt.Fprintf(c.writer, "IPAddresses: %v\n", cert.IPAddresses)
	fmt.Fprintf(c.writer, "NotBefore: %v\n", cert.NotBefore)
	fmt.Fprintf(c.writer, "NotAfter: %v\n", cert.NotAfter)
	return nil
}

func (c *Controller) exportTls(cert, key []byte) error {
	crtFileName := fmt.Sprintf("%s-cert.pem", c.exportName)
	keyFileName := fmt.Sprintf("%s-privkey.pem", c.exportName)

	if _, err := os.Stat(crtFileName); err == nil {
		return fmt.Errorf("File %s exists", crtFileName)
	}
	if _, err := os.Stat(keyFileName); err == nil {
		return fmt.Errorf("File %s exists", keyFileName)
	}

	crtFile, err := os.Create(crtFileName)
	if err != nil {
		return err
	}
	keyFile, err := os.Create(keyFileName)
	if err != nil {
		return err
	}

	_, err = crtFile.Write(cert)
	if err != nil {
		return err
	}
	_, err = keyFile.Write(key)
	if err != nil {
		return err
	}
	return nil
}
