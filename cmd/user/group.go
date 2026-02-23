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
	"fmt"
	"strings"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func groupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "groups <userid>",
		Short: "Show which groups a user belongs to",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var resp struct {
				Data map[string]any `json:"data"`
			}

			if err := client.Get(fmt.Sprintf("/access/users/%s", args[0]), &resp); err != nil {
				return err
			}

			groups := toString(resp.Data["groups"])

			if output.IsJSON() {
				return output.JSON(map[string]any{
					"userid": args[0],
					"groups": strings.Split(groups, ","),
				})
			}

			if groups == "" {
				fmt.Printf("User '%s' is not a member of any groups.\n", args[0])
				return nil
			}

			headers := []string{"GROUP"}
			rows := [][]string{}

			for _, g := range strings.Split(groups, ",") {
				g = strings.TrimSpace(g)
				if g != "" {
					rows = append(rows, []string{g})
				}
			}

			output.Table(headers, rows)

			return nil
		},
	}
}
