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
package backup

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func deleteCmd() *cobra.Command {
	var (
		node    string
		storage string
		file    string
		force   bool
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a backup file from storage",
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

			volid := fmt.Sprintf("%s:%s", storage, file)

			if !force {
				var confirm string

				fmt.Printf("Delete backup '%s'? This cannot be undone. [y/N]: ", volid)
				_, _ = fmt.Scanln(&confirm)

				if confirm != "y" && confirm != "Y" {
					output.Aborted("Aborted.")
					return nil
				}
			}

			if err := client.Delete(
				fmt.Sprintf("/nodes/%s/storage/%s/content/%s", node, storage, file),
			); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("Backup '%s' deleted", volid))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node (auto-detected if not set)")
	cmd.Flags().StringVar(&storage, "storage", "", "Storage containing the backup (required)")
	cmd.Flags().StringVar(&file, "file", "", "Backup filename to delete (required)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("storage")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
