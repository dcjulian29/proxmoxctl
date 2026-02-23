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

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func modifyCmd() *cobra.Command {
	var (
		firstname string
		lastname  string
		email     string
		comment   string
		groups    string
		expire    int64
		enabled   string // "true" | "false" | "" (unset = don't change)
	)

	cmd := &cobra.Command{
		Use:   "modify <userid>",
		Short: "Modify an existing Proxmox user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			payload := map[string]any{}

			if firstname != "" {
				payload["firstname"] = firstname
			}

			if lastname != "" {
				payload["lastname"] = lastname
			}

			if email != "" {
				payload["email"] = email
			}

			if comment != "" {
				payload["comment"] = comment
			}

			if groups != "" {
				payload["groups"] = groups
			}

			if expire >= 0 && cmd.Flags().Changed("expire") {
				payload["expire"] = expire
			}

			switch enabled {
			case "true":
				payload["enable"] = true
			case "false":
				payload["enable"] = false
			}

			if len(payload) == 0 {
				return fmt.Errorf("no changes specified — use --firstname, --lastname, --email, --groups, --enabled, or --expire")
			}

			if err := client.Put(fmt.Sprintf("/access/users/%s", args[0]), payload, nil); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("User '%s' updated", args[0]))

			return nil
		},
	}

	cmd.Flags().StringVar(&firstname, "firstname", "", "New first name")
	cmd.Flags().StringVar(&lastname, "lastname", "", "New last name")
	cmd.Flags().StringVar(&email, "email", "", "New email address")
	cmd.Flags().StringVar(&comment, "comment", "", "New comment")
	cmd.Flags().StringVar(&groups, "groups", "", "Comma-separated group IDs (replaces current groups)")
	cmd.Flags().Int64Var(&expire, "expire", -1, "New expiry Unix timestamp (0 = never)")
	cmd.Flags().StringVar(&enabled, "enabled", "", "Enable or disable account: true or false")

	return cmd
}
