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
package lxc

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	var node string
	cmd := &cobra.Command{
		Use:   "status <vmid>",
		Short: "Show status of an LXC container",
		Args:  cobra.ExactArgs(1),
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

			var resp struct {
				Data map[string]any `json:"data"`
			}

			if err := client.Get(fmt.Sprintf("/nodes/%s/lxc/%s/status/current", node, args[0]), &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			headers := []string{"FIELD", "VALUE"}
			rows := [][]string{
				{"VMID", toString(resp.Data["vmid"])},
				{"Name", toString(resp.Data["name"])},
				{"Status", toString(resp.Data["status"])},
				{"CPU Usage", fmt.Sprintf("%.2f%%", toFloat(resp.Data["cpu"])*100)},
				{"Memory", fmt.Sprintf("%.0f / %.0f MB", toFloat(resp.Data["mem"])/1024/1024, toFloat(resp.Data["maxmem"])/1024/1024)},
				{"Uptime", fmt.Sprintf("%.0fs", toFloat(resp.Data["uptime"]))},
			}

			output.Table(headers, rows)

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name")

	return cmd
}
