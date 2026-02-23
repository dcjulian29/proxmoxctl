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
package job

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all scheduled backup jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var resp struct {
				Data []map[string]any `json:"data"`
			}

			if err := client.Get("/cluster/backup", &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			if len(resp.Data) == 0 {
				fmt.Println("No scheduled backup jobs found.")
				return nil
			}

			headers := []string{"JOB ID", "VMIDS", "STORAGE", "SCHEDULE", "MODE", "ENABLED"}
			rows := make([][]string, 0, len(resp.Data))

			for _, j := range resp.Data {
				rows = append(rows, []string{
					toString(j["id"]),
					toString(j["vmid"]),
					toString(j["storage"]),
					toString(j["schedule"]),
					toString(j["mode"]),
					formatBool(j["enabled"]),
				})
			}

			output.Table(headers, rows)

			return nil
		},
	}
}
