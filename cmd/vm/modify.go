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

func modifyCmd() *cobra.Command {
	var node, memory, cores, name string

	cmd := &cobra.Command{
		Use:   "modify <vmid>",
		Short: "Modify an existing VM's configuration",
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

			payload := map[string]any{}

			if name != "" {
				payload["name"] = name
			}

			if memory != "" {
				payload["memory"] = memory
			}

			if cores != "" {
				payload["cores"] = cores
			}

			if len(payload) == 0 {
				return fmt.Errorf("no changes specified — use --name, --memory, or --cores")
			}

			if err := client.Put(fmt.Sprintf("/nodes/%s/qemu/%s/config", node, args[0]), payload, nil); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("VM %s updated successfully", args[0]))

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node name")
	cmd.Flags().StringVar(&name, "name", "", "New VM name")
	cmd.Flags().StringVar(&memory, "memory", "", "New memory in MB")
	cmd.Flags().StringVar(&cores, "cores", "", "New CPU core count")

	return cmd
}
