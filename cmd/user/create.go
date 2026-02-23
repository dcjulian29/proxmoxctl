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

func createCmd() *cobra.Command {
	var (
		password  string
		firstname string
		lastname  string
		email     string
		comment   string
		groups    string
		expire    int64
		enabled   bool
	)

	cmd := &cobra.Command{
		Use:   "create <userid>",
		Short: "Create a new Proxmox user (e.g. alice@pam)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !strings.Contains(args[0], "@") {
				return fmt.Errorf("userid must be in USER@REALM format (e.g. alice@pam)")
			}

			client, err := api.New()
			if err != nil {
				return err
			}

			payload := map[string]any{
				"userid": args[0],
				"enable": enabled,
			}

			if password != "" {
				payload["password"] = password
			}

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

			if expire > 0 {
				payload["expire"] = expire
			}

			var resp any
			if err := client.Post("/access/users", payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("User '%s' created", args[0]))

			return nil
		},
	}

	cmd.Flags().StringVar(&password, "password", "", "Initial password (for @pam/@pve realms)")
	cmd.Flags().StringVar(&firstname, "firstname", "", "First name")
	cmd.Flags().StringVar(&lastname, "lastname", "", "Last name")
	cmd.Flags().StringVar(&email, "email", "", "Email address")
	cmd.Flags().StringVar(&comment, "comment", "", "Optional comment")
	cmd.Flags().StringVar(&groups, "groups", "", "Comma-separated list of group IDs to assign")
	cmd.Flags().Int64Var(&expire, "expire", 0, "Expiry as Unix timestamp (0 = never)")
	cmd.Flags().BoolVar(&enabled, "enabled", true, "Whether the account is enabled")

	return cmd
}
