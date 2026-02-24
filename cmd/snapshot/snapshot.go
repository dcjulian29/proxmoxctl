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
	"time"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage snapshots for VMs and LXC containers",
		Long: `Manage snapshots for both KVM VMs and LXC containers.

Use --type to target a VM (qemu, default) or an LXC container (lxc).

Examples:
  proxmoxctl snapshot list 100
  proxmoxctl snapshot list 101 --type lxc
  proxmoxctl snapshot create 100 --name before-update --desc "Pre-upgrade snapshot"
  proxmoxctl snapshot rollback 100 --name before-update
  proxmoxctl snapshot delete 100 --name before-update`,
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(rollbackCmd())
	cmd.AddCommand(showCmd())

	return cmd
}

func formatEpoch(v any) string {
	if v == nil {
		return ""
	}

	if val, ok := v.(float64); ok {
		return time.Unix(int64(val), 0).Format("2006-01-02 15:04:05")
	}

	return fmt.Sprintf("%v", v)
}

func resolveNode(client *api.Client, node string) (string, error) {
	if node != "" {
		return node, nil
	}

	return client.DefaultNode()
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}
