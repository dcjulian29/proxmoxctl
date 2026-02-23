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
	"strings"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func resourcesCmd() *cobra.Command {
	var resourceType string

	cmd := &cobra.Command{
		Use:   "resources",
		Short: "List all cluster resources (VMs, containers, storage, nodes)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			path := "/cluster/resources"
			if resourceType != "" {
				path += "?type=" + resourceType
			}

			var resp struct {
				Data []map[string]interface{} `json:"data"`
			}

			if err := client.Get(path, &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			grouped := map[string][]map[string]interface{}{}
			order := []string{}

			for _, r := range resp.Data {
				t := toString(r["type"])
				if _, exists := grouped[t]; !exists {
					order = append(order, t)
				}

				grouped[t] = append(grouped[t], r)
			}

			for _, t := range order {
				items := grouped[t]
				fmt.Println()
				printSectionHeader(strings.ToUpper(t) + "S")

				switch t {
				case "qemu", "lxc":
					headers := []string{"ID", "NAME", "NODE", "STATUS", "CPU", "MEM", "DISK", "UPTIME"}
					rows := make([][]string, 0, len(items))

					for _, r := range items {
						rows = append(rows, []string{
							fmt.Sprintf("%.0f", toFloat(r["vmid"])),
							toString(r["name"]),
							toString(r["node"]),
							toString(r["status"]),
							formatPct(toFloat(r["cpu"])),
							formatBytes(toFloat(r["mem"])),
							formatBytes(toFloat(r["disk"])),
							formatUptime(toFloat(r["uptime"])),
						})
					}

					output.Table(headers, rows)

				case "storage":
					headers := []string{"NAME", "NODE", "STATUS", "USED", "TOTAL", "USAGE"}
					rows := make([][]string, 0, len(items))

					for _, r := range items {
						rows = append(rows, []string{
							toString(r["storage"]),
							toString(r["node"]),
							toString(r["status"]),
							formatBytes(toFloat(r["disk"])),
							formatBytes(toFloat(r["maxdisk"])),
							formatBar(toFloat(r["disk"]), toFloat(r["maxdisk"]), 16),
						})
					}

					output.Table(headers, rows)

				case "node":
					headers := []string{"NODE", "STATUS", "CPU", "MEM USED", "MEM TOTAL", "UPTIME"}
					rows := make([][]string, 0, len(items))

					for _, r := range items {
						rows = append(rows, []string{
							toString(r["node"]),
							toString(r["status"]),
							formatPct(toFloat(r["cpu"])),
							formatBytes(toFloat(r["mem"])),
							formatBytes(toFloat(r["maxmem"])),
							formatUptime(toFloat(r["uptime"])),
						})
					}

					output.Table(headers, rows)

				default:
					headers := []string{"ID", "TYPE", "STATUS", "NODE"}
					rows := make([][]string, 0, len(items))

					for _, r := range items {
						rows = append(rows, []string{
							toString(r["id"]),
							toString(r["type"]),
							toString(r["status"]),
							toString(r["node"]),
						})
					}

					output.Table(headers, rows)
				}
			}

			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVar(&resourceType, "type", "", "Filter by type: vm, lxc, storage, or node")

	return cmd
}
