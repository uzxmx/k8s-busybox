package eventgenerator

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"
	"k8s.io/klog"
	"time"
)

// Generator generates event.
type Generator struct {
	kind      string
	name      string
	namespace string

	eventType    string
	eventAction  string
	eventReason  string
	eventMessage string

	restClientGetter resource.RESTClientGetter
}

const (
	defaultNamespace = "default"
	defaultEventType = v1.EventTypeNormal
)

// NewGenerator creates a generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// AddFlags adds flags to command.
func (g *Generator) AddFlags(cmd *cobra.Command) {
	configFlags := genericclioptions.NewConfigFlags(true)
	configFlags.AddFlags(cmd.PersistentFlags())
	g.restClientGetter = configFlags

	flags := cmd.Flags()
	flags.StringVar(&g.kind, "kind", "", "Resource kind to get.")
	flags.StringVar(&g.name, "name", "", "Resource name to get.")
	flags.StringVar(&g.namespace, "namespace", defaultNamespace, "Resource namespace to get.")
	flags.StringVar(&g.eventType, "type", defaultEventType, "Event type.")
	flags.StringVar(&g.eventAction, "action", "", "Event action.")
	flags.StringVar(&g.eventReason, "reason", "", "Event reason.")
	flags.StringVar(&g.eventMessage, "message", "", "Event message.")
}

// Run generates event.
func (g *Generator) Run() error {
	r := resource.NewBuilder(g.restClientGetter).
		Unstructured().
		NamespaceParam(g.namespace).
		ResourceTypeOrNameArgs(true, g.kind, g.name).
		Do()
	if err := r.Err(); err != nil {
		return err
	}

	infos, err := r.Infos()
	if err != nil {
		return err
	}

	ref, err := reference.GetReference(scheme.Scheme, infos[0].Object)
	if err != nil {
		return err
	}

	restConfig, err := g.restClientGetter.ToRESTConfig()
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	if len(g.eventAction) == 0 {
		g.eventAction = g.eventReason
	}

	now := time.Now()
	event, err := client.CoreV1().Events("").CreateWithEventNamespace(&v1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v.%x", g.name, now.UnixNano()),
			Namespace: g.namespace,
		},
		FirstTimestamp:      metav1.NewTime(now),
		LastTimestamp:       metav1.NewTime(now),
		EventTime:           metav1.NewMicroTime(now),
		ReportingController: "eventgenerator",
		ReportingInstance:   "eventgenerator",
		Action:              g.eventAction,
		InvolvedObject:      *ref,
		Reason:              g.eventReason,
		Type:                g.eventType,
		Message:             g.eventMessage,
	})

	if err == nil {
		klog.Infof("Event generated successfully: %v", event)
	}

	return err
}
