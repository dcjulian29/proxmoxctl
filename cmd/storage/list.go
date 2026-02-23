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

func listCmd() *cobra.Command {
	var (
		node    string
		active  bool
		content string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List storage pools",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var path string

			if node != "" {
				path = fmt.Sprintf("/nodes/%s/storage", node)
			} else {
				path = "/storage"
			}

			if active {
				path += "?enabled=1"
			}

			if content != "" {
				sep := "?"
				if strings.Contains(path, "?") {
					sep = "&"
				}

				path += sep + "content=" + content
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
				fmt.Println("No storage pools found.")
				return nil
			}

			if node != "" {
				fmt.Printf("\n  STORAGE — NODE: %s\n", strings.ToUpper(node))
				fmt.Println("  " + strings.Repeat("─", 70))

				headers := []string{"NAME", "TYPE", "STATUS", "USED", "AVAIL", "TOTAL", "USAGE", "CONTENT"}
				rows := make([][]string, 0, len(resp.Data))

				for _, s := range resp.Data {
					rows = append(rows, []string{
						toString(s["storage"]),
						toString(s["type"]),
						toString(s["status"]),
						formatBytes(toFloat(s["used"])),
						formatBytes(toFloat(s["avail"])),
						formatBytes(toFloat(s["total"])),
						formatBar(toFloat(s["used"]), toFloat(s["total"]), 14),
						formatContent(s["content"]),
					})
				}

				output.Table(headers, rows)
			} else {
				fmt.Println("\n  STORAGE — CLUSTER CONFIG")
				fmt.Println("  " + strings.Repeat("─", 60))

				headers := []string{"NAME", "TYPE", "SHARED", "ENABLED", "CONTENT", "PATH / SERVER"}
				rows := make([][]string, 0, len(resp.Data))

				for _, s := range resp.Data {
					location := toString(s["path"])
					if location == "" {
						location = toString(s["server"])
					}

					rows = append(rows, []string{
						toString(s["storage"]),
						toString(s["type"]),
						formatBool(s["shared"]),
						formatEnabled(s["disable"]),
						formatContent(s["content"]),
						location,
					})
				}

				output.Table(headers, rows)
			}

			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Show live usage stats for a specific node")
	cmd.Flags().BoolVar(&active, "active", false, "Only show enabled/active storage pools")
	cmd.Flags().StringVar(&content, "content", "", "Filter by content type (e.g. backup, iso, images)")

	return cmd
}
