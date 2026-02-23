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

func nodeStatusCmd() *cobra.Command {
	var node string

	cmd := &cobra.Command{
		Use:   "node",
		Short: "Show detailed status of a single node",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			if node == "" {
				node, err = client.DefaultNode()
				if err != nil {
					return err
				}
			}

			var statusResp struct {
				Data map[string]any `json:"data"`
			}

			if err := client.Get(fmt.Sprintf("/nodes/%s/status", node), &statusResp); err != nil {
				return err
			}

			var versionResp struct {
				Data map[string]any `json:"data"`
			}

			_ = client.Get(fmt.Sprintf("/nodes/%s/version", node), &versionResp)

			if output.IsJSON() {
				combined := map[string]any{
					"status":  statusResp.Data,
					"version": versionResp.Data,
				}

				return output.JSON(combined)
			}

			d := statusResp.Data
			v := versionResp.Data

			fmt.Println()
			printSectionHeader(fmt.Sprintf("NODE: %s", strings.ToUpper(node)))

			headers := []string{"FIELD", "VALUE"}
			overviewRows := [][]string{
				{"Hostname", toString(d["pveversion"])},
				{"Kernel", toString(d["kversion"])},
				{"PVE Version", toString(v["version"])},
				{"Uptime", formatUptime(toFloat(d["uptime"]))},
				{"Timezone", toString(d["timezone"])},
			}

			output.Table(headers, overviewRows)

			cpuInfo, _ := d["cpuinfo"].(map[string]any)
			fmt.Println()
			printSectionHeader("CPU")
			cpuRows := [][]string{
				{"Model", toString(cpuInfo["model"])},
				{"Sockets", fmt.Sprintf("%.0f", toFloat(cpuInfo["sockets"]))},
				{"Cores per Socket", fmt.Sprintf("%.0f", toFloat(cpuInfo["cores"]))},
				{"Threads (total)", fmt.Sprintf("%.0f", toFloat(cpuInfo["cpus"]))},
				{"Current Usage", formatPct(toFloat(d["cpu"]))},
				{"Load Avg (1/5/15m)", formatLoadAvg(d["loadavg"])},
			}

			output.Table(headers, cpuRows)

			memInfo, _ := d["memory"].(map[string]any)
			fmt.Println()
			printSectionHeader("MEMORY")
			memRows := [][]string{
				{"Used", formatBytes(toFloat(memInfo["used"]))},
				{"Free", formatBytes(toFloat(memInfo["free"]))},
				{"Total", formatBytes(toFloat(memInfo["total"]))},
				{"Usage", formatBar(toFloat(memInfo["used"]), toFloat(memInfo["total"]), 20)},
			}

			output.Table(headers, memRows)

			swapInfo, _ := d["swap"].(map[string]any)
			fmt.Println()
			printSectionHeader("SWAP")
			swapRows := [][]string{
				{"Used", formatBytes(toFloat(swapInfo["used"]))},
				{"Free", formatBytes(toFloat(swapInfo["free"]))},
				{"Total", formatBytes(toFloat(swapInfo["total"]))},
				{"Usage", formatBar(toFloat(swapInfo["used"]), toFloat(swapInfo["total"]), 20)},
			}

			output.Table(headers, swapRows)

			rootfsInfo, _ := d["rootfs"].(map[string]any)
			fmt.Println()
			printSectionHeader("ROOT FILESYSTEM")
			rootfsRows := [][]string{
				{"Used", formatBytes(toFloat(rootfsInfo["used"]))},
				{"Free", formatBytes(toFloat(rootfsInfo["free"]))},
				{"Total", formatBytes(toFloat(rootfsInfo["total"]))},
				{"Usage", formatBar(toFloat(rootfsInfo["used"]), toFloat(rootfsInfo["total"]), 20)},
			}

			output.Table(headers, rootfsRows)

			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Node name (auto-detected if not set)")

	return cmd
}
