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
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone a KVM VM or LXC container",
		Long: `Clone a KVM VM or LXC container into a new guest.

Examples:
  # Full clone of VM 100 → new VM 200
  proxmoxctl clone vm 100 --newid 200 --name webserver-copy

  # Linked clone (faster, shares base disk) of VM 100 → 201
  proxmoxctl clone vm 100 --newid 201 --name webserver-linked --linked

  # Clone from a specific snapshot
  proxmoxctl clone vm 100 --newid 202 --name from-snap --snapname pre-upgrade

  # Clone LXC container 300 → 400
  proxmoxctl clone lxc 300 --newid 400 --hostname new-container`,
	}

	cmd.AddCommand(lxcCmd())
	cmd.AddCommand(vmCmd())

	return cmd
}
