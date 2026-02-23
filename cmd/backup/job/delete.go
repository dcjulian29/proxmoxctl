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

func deleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <jobid>",
		Short: "Delete a scheduled backup job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			if !force {
				var confirm string

				fmt.Printf("Delete backup job '%s'? [y/N]: ", args[0])
				_, _ = fmt.Scanln(&confirm)

				if confirm != "y" && confirm != "Y" {
					output.Aborted("Aborted.")
					return nil
				}
			}

			if err := client.Delete(fmt.Sprintf("/cluster/backup/%s", args[0])); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("Backup job '%s' deleted", args[0]))

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	return cmd
}
