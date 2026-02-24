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
package clone

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func vmCmd() *cobra.Command {
	var (
		node     string
		newid    int
		name     string
		snapname string
		pool     string
		storage  string
		linked   bool
		full     bool
	)

	cmd := &cobra.Command{
		Use:   "vm <vmid>",
		Short: "Clone a KVM VM",
		Args:  cobra.ExactArgs(1),
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

			if newid == 0 {
				return fmt.Errorf("--newid is required")
			}

			payload := map[string]any{
				"newid": newid,
			}

			if name != "" {
				payload["name"] = name
			}

			if snapname != "" {
				payload["snapname"] = snapname
			}

			if pool != "" {
				payload["pool"] = pool
			}

			if storage != "" {
				payload["storage"] = storage
			}

			if linked {
				payload["full"] = 0
			} else {
				payload["full"] = 1
			}

			if full {
				payload["full"] = 1
			}

			var resp any

			path := fmt.Sprintf("/nodes/%s/qemu/%s/clone", node, args[0])

			if err := client.Post(path, payload, &resp); err != nil {
				return err
			}

			cloneType := "full"

			if linked {
				cloneType = "linked"
			}

			output.Success(fmt.Sprintf(
				"%s clone of VM %s → VM %d queued on node %s",
				cloneType, args[0], newid, node,
			))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node (auto-detected if not set)")
	cmd.Flags().IntVar(&newid, "newid", 0, "New VM ID for the clone (required)")
	cmd.Flags().StringVar(&name, "name", "", "Name for the cloned VM")
	cmd.Flags().StringVar(&snapname, "snapname", "", "Clone from this snapshot instead of current state")
	cmd.Flags().StringVar(&pool, "pool", "", "Resource pool to assign the clone to")
	cmd.Flags().StringVar(&storage, "storage", "", "Target storage for the cloned disks (e.g. local-lvm)")
	cmd.Flags().BoolVar(&linked, "linked", false, "Create a linked clone (shares base disk; source must be a template)")
	cmd.Flags().BoolVar(&full, "full", false, "Force a full independent clone (default behavior)")

	_ = cmd.MarkFlagRequired("newid")

	return cmd
}
