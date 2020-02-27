package main

import (
	"github.com/spf13/cobra"
	"github.com/uzxmx/k8s-busybox/pkg/eventgenerator"
	"k8s.io/klog"
)

func main() {
	gen := eventgenerator.NewGenerator()

	cmd := &cobra.Command{
		Use:   "eventgenerator",
		Short: "Generate fake events in kubernetes cluster",
		Long: `
eventgenerator is an utility that can help you generate fake events in kubernetes
cluster, especially useful when you use event-based tools like brigade.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := gen.Run(); err != nil {
				klog.Error(err)
			}
		},
	}
	gen.AddFlags(cmd)

	if err := cmd.Execute(); err != nil {
		klog.Error(err)
	}
}
