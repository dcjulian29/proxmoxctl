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
package vm

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var node string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all KVM VMs",
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
				Data []map[string]any `json:"data"`
			}

			if err := client.Get(fmt.Sprintf("/nodes/%s/qemu", node), &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			headers := []string{"VMID", "NAME", "STATUS", "MEM(MB)", "CPUS"}
			rows := make([][]string, 0, len(resp.Data))

			for _, v := range resp.Data {
				rows = append(rows, []string{
					fmt.Sprintf("%.0f", toFloat(v["vmid"])),
					toString(v["name"]),
					toString(v["status"]),
					fmt.Sprintf("%.0f", toFloat(v["maxmem"])/1024/1024),
					fmt.Sprintf("%.0f", toFloat(v["cpus"])),
				})
			}

			output.Table(headers, rows)

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name (auto-detected if not set)")

	return cmd
}
