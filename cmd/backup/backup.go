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
	"time"

	"github.com/dcjulian29/proxmoxctl/cmd/backup/job"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Manage backups for VMs and LXC containers",
		Long: `Manage Proxmox backups via the vzdump API.

Covers on-demand backups, listing backup files in storage,
restoring from a backup, deleting backup files, and managing
scheduled backup jobs.

Examples:
  # Immediately back up VM 100
  proxmoxctl backup create 100 --storage local --mode snapshot

  # Back up multiple guests at once
  proxmoxctl backup create 100,101,102 --storage backup-nfs

  # List backups on a storage
  proxmoxctl backup list --storage local

  # List backups for a specific VM
  proxmoxctl backup list --storage local --vmid 100

  # Restore VM 100 from a backup file
  proxmoxctl backup restore --vmid 100 --storage local --file vzdump-qemu-100-2024_01_15-03_00_01.vma.zst

  # Show scheduled backup jobs
  proxmoxctl backup jobs list

  # Create a scheduled backup job
  proxmoxctl backup jobs create --vmids 100,101 --storage backup-nfs --schedule "0 2 * * *"`,
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(showCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(restoreCmd())
	cmd.AddCommand(job.NewCommand())

	return cmd
}

func formatBool(v any) string {
	if v == nil {
		return "false"
	}

	switch val := v.(type) {
	case bool:
		if val {
			return "true"
		}

		return "false"
	case float64:
		if val == 1 {
			return "true"
		}

		return "false"
	}

	return "false"
}

func formatBytes(b float64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%.0f B", b)
	}

	div, exp := float64(unit), 0

	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", b/div, "KMGTPE"[exp])
}

func formatEpoch(v any) string {
	if v == nil {
		return ""
	}

	if val, ok := v.(float64); ok && val > 0 {
		return time.Unix(int64(val), 0).Format("2006-01-02 15:04:05")
	}

	return ""
}

func toFloat(v any) float64 {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case float64:
		return val
	}

	return 0
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}
