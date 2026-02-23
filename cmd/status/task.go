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

func tasksCmd() *cobra.Command {
	var (
		node   string
		limit  int
		errors bool
	)

	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Show recent tasks on a node or across the cluster",
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

			path := fmt.Sprintf("/nodes/%s/tasks?limit=%d", node, limit)
			if errors {
				path += "&errors=1"
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
				fmt.Println("No tasks found.")
				return nil
			}

			fmt.Println()
			printSectionHeader(fmt.Sprintf("RECENT TASKS — NODE: %s", strings.ToUpper(node)))
			headers := []string{"UPID", "TYPE", "USER", "STATUS", "STARTED", "ENDED"}
			rows := make([][]string, 0, len(resp.Data))

			for _, t := range resp.Data {
				rows = append(rows, []string{
					truncate(toString(t["upid"]), 40),
					toString(t["type"]),
					toString(t["user"]),
					taskStatus(t),
					formatEpoch(t["starttime"]),
					formatEpoch(t["endtime"]),
				})
			}

			output.Table(headers, rows)
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Node name (auto-detected if not set)")
	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum number of tasks to show")
	cmd.Flags().BoolVar(&errors, "errors", false, "Show only failed tasks")

	return cmd
}
