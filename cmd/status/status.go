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
package status

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show status of a node or the cluster",
		Long: `Show health, resource usage, and summary information for a
Proxmox node or the entire cluster.

Examples:
  # Cluster-wide overview
  proxmoxctl status cluster

  # Status of a specific node (auto-detected if omitted)
  proxmoxctl status node
  proxmoxctl status node --node pve2

  # Resource summary across the cluster (VMs, containers, storage)
  proxmoxctl status resources

  # Running tasks across the cluster
  proxmoxctl status tasks`,
	}

	cmd.AddCommand(clusterStatusCmd())
	cmd.AddCommand(nodeStatusCmd())
	cmd.AddCommand(resourcesCmd())
	cmd.AddCommand(tasksCmd())

	return cmd
}

func printSectionHeader(title string) {
	fmt.Printf("  %s\n", title)
	fmt.Println("  " + strings.Repeat("─", len(title)+2))
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}

func toFloat(v any) float64 {
	if v == nil {
		return 0
	}

	if val, ok := v.(float64); ok {
		return val
	}

	return 0
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

func formatPct(v float64) string {
	return fmt.Sprintf("%.1f%%", v*100)
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

func formatBar(used, total float64, width int) string {
	if total == 0 {
		return "n/a"
	}

	pct := used / total
	filled := min(int(math.Round(pct*float64(width))), width)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	return fmt.Sprintf("[%s] %.1f%%", bar, pct*100)
}

func formatUptime(seconds float64) string {
	if seconds == 0 {
		return "—"
	}

	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}

	return fmt.Sprintf("%dm", minutes)
}

func formatEpoch(v any) string {
	if v == nil {
		return ""
	}

	if val, ok := v.(float64); ok && val > 0 {
		return time.Unix(int64(val), 0).Format("2006-01-02 15:04")
	}

	return ""
}

func formatLoadAvg(v any) string {
	arr, ok := v.([]any)
	if !ok || len(arr) < 3 {
		return "n/a"
	}

	return fmt.Sprintf("%.2f / %.2f / %.2f",
		toFloat(arr[0]), toFloat(arr[1]), toFloat(arr[2]))
}

func nodeStatus(n map[string]any) string {
	if toFloat(n["online"]) == 1 || n["online"] == true {
		return "online"
	}

	return "offline"
}

func taskStatus(t map[string]any) string {
	status := toString(t["status"])
	if status == "" {
		if toString(t["endtime"]) == "" {
			return "running"
		}

		return "unknown"
	}

	return status
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-1] + "…"
}
