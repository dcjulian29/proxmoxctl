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
package storage

import (
	"fmt"
	"strings"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func contentCmd() *cobra.Command {
	var (
		node        string
		contentType string
		vmid        string
	)

	cmd := &cobra.Command{
		Use:   "content <storage>",
		Short: "List the contents of a storage pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			if node == "" {
				var err error
				node, err = client.DefaultNode()
				if err != nil {
					return err
				}
			}

			path := fmt.Sprintf("/nodes/%s/storage/%s/content", node, args[0])
			params := []string{}

			if contentType != "" {
				params = append(params, "content="+contentType)
			}

			if vmid != "" {
				params = append(params, "vmid="+vmid)
			}

			if len(params) > 0 {
				path += "?" + strings.Join(params, "&")
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
				fmt.Printf("No content found in storage '%s'.\n", args[0])
				return nil
			}

			fmt.Printf("\n  CONTENT — %s\n", strings.ToUpper(args[0]))
			fmt.Println("  " + strings.Repeat("─", 72))

			headers := []string{"VOLUME ID", "TYPE", "FORMAT", "SIZE", "VMID", "NOTES"}
			rows := make([][]string, 0, len(resp.Data))
			for _, item := range resp.Data {
				rows = append(rows, []string{
					toString(item["volid"]),
					toString(item["content"]),
					toString(item["format"]),
					formatBytes(toFloat(item["size"])),
					formatVMID(item["vmid"]),
					toString(item["notes"]),
				})
			}

			output.Table(headers, rows)

			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Node to query (auto-detected if not set)")
	cmd.Flags().StringVar(&contentType, "type", "", "Filter by content type: backup, iso, vztmpl, images, rootdir, snippets")
	cmd.Flags().StringVar(&vmid, "vmid", "", "Filter by VM or container ID")

	return cmd
}
