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
package snapshot

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func rollbackCmd() *cobra.Command {
	var node, gtype, snapname string
	var force bool

	cmd := &cobra.Command{
		Use:   "rollback <vmid>",
		Short: "Roll back a VM or container to a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			node, err = resolveNode(client, node)
			if err != nil {
				return err
			}

			if !force {
				var confirm string

				fmt.Printf("Roll back %s %s to snapshot '%s'? This cannot be undone. [y/N]: ",
					gtype, args[0], snapname)
				_, _ = fmt.Scanln(&confirm)

				if confirm != "y" && confirm != "Y" {
					output.Aborted("Rollback aborted.")
					return nil
				}
			}

			var resp any

			path := fmt.Sprintf("/nodes/%s/%s/%s/snapshot/%s/rollback", node, gtype, args[0], snapname)

			if err := client.Post(path, nil, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("Rollback of %s %s to snapshot '%s' task queued", gtype, args[0], snapname))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name")
	cmd.Flags().StringVar(&gtype, "type", "qemu", "Guest type: qemu or lxc")
	cmd.Flags().StringVar(&snapname, "name", "", "Snapshot name to roll back to (required)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")

	_ = cmd.MarkFlagRequired("name")

	return cmd
}
