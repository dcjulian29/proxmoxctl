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

func createCmd() *cobra.Command {
	var node, gtype, snapname, description string
	var vmstate bool

	cmd := &cobra.Command{
		Use:   "create <vmid>",
		Short: "Create a snapshot of a VM or container",
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

			payload := map[string]any{
				"snapname": snapname,
			}

			if description != "" {
				payload["description"] = description
			}

			if vmstate && gtype == string(guestVM) {
				payload["vmstate"] = 1
			}

			var resp any

			path := fmt.Sprintf("/nodes/%s/%s/%s/snapshot", node, gtype, args[0])

			if err := client.Post(path, payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("Snapshot '%s' of %s %s creation task queued", snapname, gtype, args[0]))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name")
	cmd.Flags().StringVar(&gtype, "type", "qemu", "Guest type: qemu or lxc")
	cmd.Flags().StringVar(&snapname, "name", "", "Snapshot name (required, no spaces)")
	cmd.Flags().StringVar(&description, "desc", "", "Human-readable description")
	cmd.Flags().BoolVar(&vmstate, "vmstate", false, "Include RAM state in snapshot (VMs only, requires guest to be running)")

	_ = cmd.MarkFlagRequired("name")

	return cmd
}
