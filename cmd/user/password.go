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

func passwordCmd() *cobra.Command {
	var password string

	cmd := &cobra.Command{
		Use:   "passwd <userid>",
		Short: "Change a user's password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			if password == "" {
				var p string

				fmt.Printf("New password for %s: ", args[0])
				_, _ = fmt.Scanln(&p)

				password = strings.TrimSpace(p)
				if password == "" {
					return fmt.Errorf("password cannot be empty")
				}
			}

			payload := map[string]any{
				"userid":   args[0],
				"password": password,
			}

			if err := client.Put("/access/password", payload, nil); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("Password changed for user '%s'", args[0]))

			return nil
		},
	}

	cmd.Flags().StringVar(&password, "password", "", "New password (prompted if not provided)")

	return cmd
}
