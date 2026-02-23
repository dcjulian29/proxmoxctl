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
	"strings"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func restoreCmd() *cobra.Command {
	var (
		node    string
		storage string
		file    string
		vmid    int
		gtype   string
		target  string
		force   bool
		start   bool
	)

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore a VM or container from a backup",
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

			archive := fmt.Sprintf("%s:%s", storage, file)

			if !force {
				var confirm string

				fmt.Printf("Restore %s %d from '%s'? [y/N]: ", gtype, vmid, archive)
				_, _ = fmt.Scanln(&confirm)

				if confirm != "y" && confirm != "Y" {
					output.Aborted("Aborted.")
					return nil
				}
			}

			payload := map[string]any{
				"vmid":    vmid,
				"archive": archive,
				"force":   1,
			}

			if target != "" {
				payload["storage"] = target
			}

			if start {
				payload["start"] = 1
			}

			var apiPath string

			switch strings.ToLower(gtype) {
			case "lxc", "ct":
				apiPath = fmt.Sprintf("/nodes/%s/lxc", node)
			default:
				apiPath = fmt.Sprintf("/nodes/%s/qemu", node)
			}

			var resp any

			if err := client.Post(apiPath, payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf(
				"Restore of %s %d from '%s' task queued on node %s",
				gtype, vmid, archive, node,
			))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node (auto-detected if not set)")
	cmd.Flags().StringVar(&storage, "storage", "", "Storage containing the backup (required)")
	cmd.Flags().StringVar(&file, "file", "", "Backup filename to restore from (required)")
	cmd.Flags().IntVar(&vmid, "vmid", 0, "VM/container ID to restore into (required)")
	cmd.Flags().StringVar(&gtype, "type", "qemu", "Guest type to restore: qemu or lxc")
	cmd.Flags().StringVar(&target, "target-storage", "", "Storage for restored disks (defaults to original)")
	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	cmd.Flags().BoolVar(&start, "start", false, "Start the guest immediately after restore")

	_ = cmd.MarkFlagRequired("storage")
	_ = cmd.MarkFlagRequired("file")
	_ = cmd.MarkFlagRequired("vmid")

	return cmd
}
