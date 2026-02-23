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

func createCmd() *cobra.Command {
	var (
		node        string
		storage     string
		mode        string
		compress    string
		mailTo      string
		notes       string
		removeOlder int
	)

	cmd := &cobra.Command{
		Use:   "create <vmid>[,vmid,...]",
		Short: "Create an on-demand backup of one or more guests",
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

			vmids := args[0]

			payload := map[string]any{
				"storage":  storage,
				"mode":     mode,
				"compress": compress,
			}

			if vmids == "all" {
				payload["all"] = 1
			} else {
				payload["vmid"] = vmids
			}

			if mailTo != "" {
				payload["mailto"] = mailTo
			}

			if notes != "" {
				payload["notes-template"] = notes
			}

			if removeOlder > 0 {
				payload["remove"] = removeOlder
			}

			var resp any

			if err := client.Post(fmt.Sprintf("/nodes/%s/vzdump", node), payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf(
				"Backup task queued for guest(s) %s on node %s → storage: %s (mode: %s, compress: %s)",
				vmids, node, storage, mode, compress,
			))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node (auto-detected if not set)")
	cmd.Flags().StringVar(&storage, "storage", "", "Target storage for the backup file (required)")
	cmd.Flags().StringVar(&mode, "mode", "snapshot", "Backup mode: snapshot, suspend, or stop")
	cmd.Flags().StringVar(&compress, "compress", "zstd", "Compression: zstd, lzo, gzip, or 0 (none)")
	cmd.Flags().StringVar(&mailTo, "mailto", "", "Email address(es) for job notifications")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes template stored with the backup")
	cmd.Flags().IntVar(&removeOlder, "remove-older", 0, "Remove backups older than N days (0 = keep all)")

	_ = cmd.MarkFlagRequired("storage")

	return cmd
}
