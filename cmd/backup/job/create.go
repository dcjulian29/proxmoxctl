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

func createCmd() *cobra.Command {
	var (
		vmids     string
		storage   string
		schedule  string
		mode      string
		compress  string
		mailTo    string
		maxFiles  int
		notes     string
		enabled   bool
		allGuests bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a scheduled backup job",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			payload := map[string]any{
				"storage":  storage,
				"schedule": schedule,
				"mode":     mode,
				"compress": compress,
				"enabled":  enabled,
			}

			if allGuests {
				payload["all"] = 1
			} else if vmids != "" {
				payload["vmid"] = vmids
			} else {
				return fmt.Errorf("specify --vmids or --all")
			}

			if mailTo != "" {
				payload["mailto"] = mailTo
			}

			if maxFiles > 0 {
				payload["maxfiles"] = maxFiles
			}

			if notes != "" {
				payload["notes-template"] = notes
			}

			var resp any

			if err := client.Post("/cluster/backup", payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf(
				"Scheduled backup job created (storage: %s, schedule: %s)",
				storage, schedule,
			))

			return nil
		},
	}

	cmd.Flags().StringVar(&vmids, "vmids", "", "Comma-separated VM/container IDs to back up")
	cmd.Flags().BoolVar(&allGuests, "all", false, "Back up all guests on the node")
	cmd.Flags().StringVar(&storage, "storage", "", "Destination storage (required)")
	cmd.Flags().StringVar(&schedule, "schedule", "0 2 * * *", "Cron-style schedule (default: 2am daily)")
	cmd.Flags().StringVar(&mode, "mode", "snapshot", "Backup mode: snapshot, suspend, or stop")
	cmd.Flags().StringVar(&compress, "compress", "zstd", "Compression: zstd, lzo, gzip, or 0")
	cmd.Flags().StringVar(&mailTo, "mailto", "", "Email for notifications")
	cmd.Flags().IntVar(&maxFiles, "max-files", 0, "Max number of backups to keep per guest (0 = unlimited)")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes template for created backups")
	cmd.Flags().BoolVar(&enabled, "enabled", true, "Enable the job immediately")

	_ = cmd.MarkFlagRequired("storage")

	return cmd
}
