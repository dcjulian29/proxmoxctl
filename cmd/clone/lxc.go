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

func lxcCmd() *cobra.Command {
	var (
		node     string
		newid    int
		hostname string
		snapname string
		pool     string
		storage  string
	)

	cmd := &cobra.Command{
		Use:   "lxc <vmid>",
		Short: "Clone an LXC container",
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

			payload := map[string]any{
				"newid": newid,
			}

			if hostname != "" {
				payload["hostname"] = hostname
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

			var resp any

			path := fmt.Sprintf("/nodes/%s/lxc/%s/clone", node, args[0])

			if err := client.Post(path, payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf(
				"Clone of LXC container %s → container %d queued on node %s",
				args[0], newid, node,
			))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node (auto-detected if not set)")
	cmd.Flags().IntVar(&newid, "newid", 0, "New container ID for the clone (required)")
	cmd.Flags().StringVar(&hostname, "hostname", "", "Hostname for the cloned container")
	cmd.Flags().StringVar(&snapname, "snapname", "", "Clone from this snapshot instead of current state")
	cmd.Flags().StringVar(&pool, "pool", "", "Resource pool to assign the clone to")
	cmd.Flags().StringVar(&storage, "storage", "", "Target storage for cloned rootfs (e.g. local-lvm)")

	_ = cmd.MarkFlagRequired("newid")

	return cmd
}
