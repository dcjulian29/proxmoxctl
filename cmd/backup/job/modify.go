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

func modifyCmd() *cobra.Command {
	var (
		vmids    string
		storage  string
		schedule string
		mode     string
		compress string
		mailTo   string
		maxFiles int
		enabled  string
	)

	cmd := &cobra.Command{
		Use:   "modify <jobid>",
		Short: "Modify an existing scheduled backup job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			payload := map[string]any{}

			if vmids != "" {
				payload["vmid"] = vmids
			}

			if storage != "" {
				payload["storage"] = storage
			}

			if schedule != "" {
				payload["schedule"] = schedule
			}

			if mode != "" {
				payload["mode"] = mode
			}

			if compress != "" {
				payload["compress"] = compress
			}

			if mailTo != "" {
				payload["mailto"] = mailTo
			}

			if maxFiles > 0 {
				payload["maxfiles"] = maxFiles
			}

			switch enabled {
			case "true":
				payload["enabled"] = true
			case "false":
				payload["enabled"] = false
			}

			if len(payload) == 0 {
				return fmt.Errorf("no changes specified")
			}

			if err := client.Put(fmt.Sprintf("/cluster/backup/%s", args[0]), payload, nil); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("Backup job '%s' updated", args[0]))

			return nil
		},
	}

	cmd.Flags().StringVar(&vmids, "vmids", "", "New comma-separated VM/container IDs")
	cmd.Flags().StringVar(&storage, "storage", "", "New destination storage")
	cmd.Flags().StringVar(&schedule, "schedule", "", "New cron-style schedule")
	cmd.Flags().StringVar(&mode, "mode", "", "New backup mode")
	cmd.Flags().StringVar(&compress, "compress", "", "New compression type")
	cmd.Flags().StringVar(&mailTo, "mailto", "", "New notification email")
	cmd.Flags().IntVar(&maxFiles, "max-files", 0, "New max backup count per guest")
	cmd.Flags().StringVar(&enabled, "enabled", "", "Enable or disable: true or false")

	return cmd
}
