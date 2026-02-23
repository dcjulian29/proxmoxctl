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
package backup

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var (
		node    string
		storage string
		vmid    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List backup files on a storage",
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

			path := fmt.Sprintf("/nodes/%s/storage/%s/content?content=backup", node, storage)

			if vmid != "" {
				path += "&vmid=" + vmid
			}

			var resp struct {
				Data []map[string]any `json:"data"`
			}

			if err := client.Get(path, &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			if len(resp.Data) == 0 {
				fmt.Printf("No backups found on storage '%s'.\n", storage)
				return nil
			}

			headers := []string{"VOLID", "VMID", "FORMAT", "SIZE", "CREATED"}
			rows := make([][]string, 0, len(resp.Data))

			for _, b := range resp.Data {
				rows = append(rows, []string{
					toString(b["volid"]),
					fmt.Sprintf("%.0f", toFloat(b["vmid"])),
					toString(b["format"]),
					formatBytes(toFloat(b["size"])),
					formatEpoch(b["ctime"]),
				})
			}

			output.Table(headers, rows)

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node (auto-detected if not set)")
	cmd.Flags().StringVar(&storage, "storage", "", "Storage to list backups from (required)")
	cmd.Flags().StringVar(&vmid, "vmid", "", "Filter by VM/container ID")

	_ = cmd.MarkFlagRequired("storage")

	return cmd
}
