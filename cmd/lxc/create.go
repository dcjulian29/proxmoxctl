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
package lxc

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var node string
	var vmid int
	var hostname, memory, cores, disk, ostemplate, password string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new LXC container",
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
				"vmid":       vmid,
				"hostname":   hostname,
				"memory":     memory,
				"cores":      cores,
				"rootfs":     disk,
				"ostemplate": ostemplate,
				"password":   password,
				"net0":       "name=eth0,bridge=vmbr0,ip=dhcp",
			}

			var resp any

			if err := client.Post(fmt.Sprintf("/nodes/%s/lxc", node), payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("LXC container %d (%s) creation task queued on node %s", vmid, hostname, node))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name")
	cmd.Flags().IntVar(&vmid, "vmid", 0, "Container ID (required)")
	cmd.Flags().StringVar(&hostname, "hostname", "", "Container hostname (required)")
	cmd.Flags().StringVar(&memory, "memory", "512", "Memory in MB")
	cmd.Flags().StringVar(&cores, "cores", "1", "CPU cores")
	cmd.Flags().StringVar(&disk, "disk", "local-lvm:8", "Root FS spec (e.g. local-lvm:8)")
	cmd.Flags().StringVar(&ostemplate, "template", "", "OS template path (e.g. local:vztmpl/debian-12.tar.zst) (required)")
	cmd.Flags().StringVar(&password, "password", "", "Root password")

	_ = cmd.MarkFlagRequired("vmid")
	_ = cmd.MarkFlagRequired("hostname")
	_ = cmd.MarkFlagRequired("template")

	return cmd
}
