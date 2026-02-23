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

func showCmd() *cobra.Command {
	var node string

	cmd := &cobra.Command{
		Use:   "show <storage>",
		Short: "Show detailed configuration of a storage pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var resp struct {
				Data map[string]any `json:"data"`
			}

			var path string
			if node != "" {
				// Node-level detail (includes usage)
				path = fmt.Sprintf("/nodes/%s/storage/%s/status", node, args[0])
			} else {
				// Cluster config detail
				path = fmt.Sprintf("/storage/%s", args[0])
			}

			if err := client.Get(path, &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			d := resp.Data

			fmt.Printf("\n  STORAGE: %s\n", strings.ToUpper(args[0]))
			fmt.Println("  " + strings.Repeat("─", 40))

			rows := [][]string{
				{"Name", toString(d["storage"])},
				{"Type", toString(d["type"])},
				{"Content", formatContent(d["content"])},
				{"Shared", formatBool(d["shared"])},
				{"Enabled", formatEnabled(d["disable"])},
			}

			if toString(d["path"]) != "" {
				rows = append(rows, []string{"Path", toString(d["path"])})
			}

			if toString(d["server"]) != "" {
				rows = append(rows, []string{"Server", toString(d["server"])})
			}

			if toString(d["export"]) != "" {
				rows = append(rows, []string{"Export", toString(d["export"])})
			}

			if toString(d["pool"]) != "" {
				rows = append(rows, []string{"Pool", toString(d["pool"])})
			}

			if toString(d["datastore"]) != "" {
				rows = append(rows, []string{"Datastore", toString(d["datastore"])})
			}

			if toFloat(d["total"]) > 0 {
				rows = append(rows, []string{"Used", formatBytes(toFloat(d["used"]))})
				rows = append(rows, []string{"Avail", formatBytes(toFloat(d["avail"]))})
				rows = append(rows, []string{"Total", formatBytes(toFloat(d["total"]))})
				rows = append(rows, []string{"Usage", formatBar(toFloat(d["used"]), toFloat(d["total"]), 20)})
			}

			output.Table([]string{"FIELD", "VALUE"}, rows)

			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Node to query for live usage stats")

	return cmd
}
