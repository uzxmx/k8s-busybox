package main

import (
	"github.com/spf13/cobra"
	"github.com/uzxmx/k8s-busybox/pkg/tls"
	"k8s.io/klog"
)

func main() {
	c := tls.NewController()

	cmd := &cobra.Command{
		Use:   "tlsinfo",
		Short: "Show secret tls info in kubernetes cluster",
		Long: `
tlsinfo is an utility that can help you get the secret tls information in kubernetes
cluster, e.g. certificate common name, expiration.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.Run(); err != nil {
				klog.Error(err)
			}
		},
	}
	c.AddFlags(cmd)

	if err := cmd.Execute(); err != nil {
		klog.Error(err)
	}
}
