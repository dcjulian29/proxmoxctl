/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package status

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func clusterStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster",
		Short: "Show cluster-wide health and node summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var resp struct {
				Data []map[string]any `json:"data"`
			}

			if err := client.Get("/cluster/status", &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			var clusterInfo map[string]any
			var nodes []map[string]any

			for _, item := range resp.Data {
				switch toString(item["type"]) {
				case "cluster":
					clusterInfo = item
				case "node":
					nodes = append(nodes, item)
				}
			}

			if clusterInfo != nil {
				fmt.Println()
				printSectionHeader("CLUSTER")
				headers := []string{"FIELD", "VALUE"}
				rows := [][]string{
					{"Name", toString(clusterInfo["name"])},
					{"Quorum", formatBool(clusterInfo["quorate"])},
					{"Nodes (total)", fmt.Sprintf("%.0f", toFloat(clusterInfo["nodes"]))},
					{"Version", toString(clusterInfo["version"])},
				}

				output.Table(headers, rows)
			}

			if len(nodes) > 0 {
				fmt.Println()
				printSectionHeader("NODES")
				headers := []string{"NODE", "STATUS", "ONLINE", "CPU", "MEM USED", "MEM TOTAL", "UPTIME"}
				rows := make([][]string, 0, len(nodes))
				for _, n := range nodes {
					rows = append(rows, []string{
						toString(n["name"]),
						nodeStatus(n),
						formatBool(n["online"]),
						formatPct(toFloat(n["cpu"])),
						formatBytes(toFloat(n["mem"])),
						formatBytes(toFloat(n["maxmem"])),
						formatUptime(toFloat(n["uptime"])),
					})
				}

				output.Table(headers, rows)
			}

			return nil
		},
	}
}
