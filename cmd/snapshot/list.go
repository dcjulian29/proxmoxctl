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
package snapshot

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var node, gtype string

	cmd := &cobra.Command{
		Use:   "list <vmid>",
		Short: "List snapshots for a VM or container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			node, err = resolveNode(client, node)
			if err != nil {
				return err
			}

			var resp struct {
				Data []map[string]any `json:"data"`
			}

			path := fmt.Sprintf("/nodes/%s/%s/%s/snapshot", node, gtype, args[0])

			if err := client.Get(path, &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			headers := []string{"NAME", "DESCRIPTION", "VMSTATE", "CREATED"}
			rows := make([][]string, 0, len(resp.Data))

			for _, s := range resp.Data {
				// Skip the synthetic "current" entry Proxmox always returns
				if toString(s["name"]) == "current" {
					continue
				}

				rows = append(rows, []string{
					toString(s["name"]),
					toString(s["description"]),
					fmt.Sprintf("%v", s["vmstate"] == true || s["vmstate"] == 1.0),
					formatEpoch(s["snaptime"]),
				})
			}

			if len(rows) == 0 {
				fmt.Printf("No snapshots found for %s %s.\n", gtype, args[0])
				return nil
			}

			output.Table(headers, rows)

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name (auto-detected if not set)")
	cmd.Flags().StringVar(&gtype, "type", "qemu", "Guest type: qemu (VM) or lxc (container)")

	return cmd
}
