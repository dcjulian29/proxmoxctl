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
package storage

import (
	"fmt"
	"math"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage",
		Short: "List and inspect Proxmox storage",
		Long: `List storage pools across the cluster or on a specific node,
and inspect the contents of individual storage pools.

Examples:
  # List all storage pools visible to the cluster
  proxmoxctl storage list

  # List storage on a specific node
  proxmoxctl storage list --node pve2

  # Filter to only enabled/active storage
  proxmoxctl storage list --active

  # Show detailed info about one storage pool
  proxmoxctl storage show local

  # List backup/ISO/template content on a storage
  proxmoxctl storage content local --type iso
  proxmoxctl storage content local --type backup
  proxmoxctl storage content local --type vztmpl`,
	}

	cmd.AddCommand(contentCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(showCmd())

	return cmd
}

func formatBar(used, total float64, width int) string {
	if total == 0 {
		return "n/a"
	}

	pct := used / total
	filled := min(int(math.Round(pct*float64(width))), width)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	return fmt.Sprintf("[%s] %.1f%%", bar, pct*100)
}

func formatBool(v any) string {
	switch val := v.(type) {
	case bool:
		if val {
			return "yes"
		}

		return "no"
	case float64:
		if val == 1 {
			return "yes"
		}

		return "no"
	}

	return "no"
}

func formatBytes(b float64) string {
	if b == 0 {
		return "0 B"
	}

	const unit = 1024.0
	if b < unit {
		return fmt.Sprintf("%.0f B", b)
	}

	div, exp := unit, 0

	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", b/div, "KMGTPE"[exp])
}

func formatContent(v any) string {
	if v == nil {
		return ""
	}

	return strings.ReplaceAll(toString(v), ",", ", ")
}

func formatEnabled(v any) string {
	switch val := v.(type) {
	case bool:
		if val {
			return "no"
		}

		return "yes"
	case float64:
		if val == 1 {
			return "no"
		}

		return "yes"
	}

	return "yes"
}

func formatVMID(v any) string {
	if v == nil {
		return "—"
	}

	if val, ok := v.(float64); ok {
		if val == 0 {
			return "—"
		}

		return fmt.Sprintf("%.0f", val)
	}

	return toString(v)
}

func toFloat(v any) float64 {
	if val, ok := v.(float64); ok {
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
