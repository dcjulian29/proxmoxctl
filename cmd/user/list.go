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
package user

import (
	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	var enabledOnly bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Proxmox users",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var resp struct {
				Data []map[string]any `json:"data"`
			}

			if err := client.Get("/access/users", &resp); err != nil {
				return err
			}

			data := resp.Data

			if enabledOnly {
				filtered := data[:0]

				for _, u := range data {
					if enabled, ok := u["enable"].(bool); ok && enabled {
						filtered = append(filtered, u)
					}
				}

				data = filtered
			}

			if output.IsJSON() {
				return output.JSON(data)
			}

			headers := []string{"USERID", "FIRSTNAME", "LASTNAME", "EMAIL", "ENABLED", "EXPIRE", "GROUPS"}
			rows := make([][]string, 0, len(data))

			for _, u := range data {
				rows = append(rows, []string{
					toString(u["userid"]),
					toString(u["firstname"]),
					toString(u["lastname"]),
					toString(u["email"]),
					formatBool(u["enable"]),
					formatExpire(u["expire"]),
					toString(u["groups"]),
				})
			}

			output.Table(headers, rows)

			return nil
		},
	}

	cmd.Flags().BoolVar(&enabledOnly, "enabled", false, "Only show enabled users")

	return cmd
}
