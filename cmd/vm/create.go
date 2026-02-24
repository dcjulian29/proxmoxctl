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
package vm

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	var node string
	var vmid int
	var name, memory, cores, disk, iso string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new KVM VM",
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
				"vmid":   vmid,
				"name":   name,
				"memory": memory,
				"cores":  cores,
				"ide2":   iso + ",media=cdrom",
				"scsi0":  disk,
				"net0":   "virtio,bridge=vmbr0",
				"ostype": "l26",
				"scsihw": "virtio-scsi-pci",
				"boot":   "order=scsi0;ide2",
			}

			var resp any

			if err := client.Post(fmt.Sprintf("/nodes/%s/qemu", node), payload, &resp); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("VM %d (%s) creation task queued on node %s", vmid, name, node))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name")
	cmd.Flags().IntVar(&vmid, "vmid", 0, "VM ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "VM name (required)")
	cmd.Flags().StringVar(&memory, "memory", "2048", "Memory in MB")
	cmd.Flags().StringVar(&cores, "cores", "2", "Number of CPU cores")
	cmd.Flags().StringVar(&disk, "disk", "local-lvm:32", "Disk spec (e.g. local-lvm:32)")
	cmd.Flags().StringVar(&iso, "iso", "", "ISO path (e.g. local:iso/debian.iso)")

	_ = cmd.MarkFlagRequired("vmid")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}
