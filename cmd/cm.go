package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	metric "github.com/neoseele/kgen/pkg"
)

// cmCmd represents the cm command
var cmCmd = &cobra.Command{
	Use:   "cm",
	Short: "Generate a custom metric scraping and exporting pipeline",
	Long: `Generate a custom metric scraping and exporting pipeline.

The build-in Prometheus scrapes The following components:

- Kubernetes api
- Nodes and Pods that have the following annotation(s)
  - cm.example.com/scrape=true
  - cm.example.com/port=xxxx (applicable to pod)

The "--metrics" flag tells "stackdriver-prometheus-sidecar" which metric(s) should be exported to Stackdriver.

Example:

kgen cm \
  --name=cm
  --namespace=default \
  --project=YOUR_PROJET \
  --cluster=YOUR_CLUSTER \
  --location=us-central1-a \
  --metrics="cilium.*" \
  --metrics="container_network.*"
`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			fmt.Println("Resource name prefix cannot be empty.")
			os.Exit(1)
		}

		namespace, _ := cmd.Flags().GetString("namespace")
		if namespace == "" {
			fmt.Println("Namespace cannot be empty.")
			os.Exit(1)
		}

		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			fmt.Println("Project cannot be empty.")
			os.Exit(1)
		}

		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			fmt.Println("Cluster cannot be empty.")
			os.Exit(1)
		}

		location, _ := cmd.Flags().GetString("location")
		if location == "" {
			fmt.Println("Location cannot be empty.")
			os.Exit(1)
		}

		metrics, _ := cmd.Flags().GetStringSlice("metrics")
		if len(metrics) == 0 {
			metrics = append(metrics, "kubelet_pleg.*")
		}

		metric.Gen(name, namespace, project, cluster, location, metrics)
	},
}

func init() {
	rootCmd.AddCommand(cmCmd)

	cmCmd.Flags().String("name", "cm", "Resource Name Prefix (default: cm)")
	cmCmd.Flags().String("namespace", "default", "Namespace Name (default: default)")
	cmCmd.Flags().StringP("project", "p", "", "Project Name")
	cmCmd.Flags().StringP("cluster", "c", "", "Cluster Name")
	cmCmd.Flags().StringP("location", "l", "", "Cluster Location (i.e. australia-southeast1-a)")
	cmCmd.Flags().StringSliceP("metrics", "m", []string{}, "Custom metrics to be exported")
}
