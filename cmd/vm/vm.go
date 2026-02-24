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
	"strconv"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "Manage KVM virtual machines",
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(modifyCmd())
	cmd.AddCommand(startCmd())
	cmd.AddCommand(statusCmd())
	cmd.AddCommand(stopCmd())

	return cmd
}

func toFloat(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}

	return 0
}

func toString(v any) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}

func vmPowerAction(node, vmid, action string) error {
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

	var resp any

	if err := client.Post(fmt.Sprintf("/nodes/%s/qemu/%s/status/%s", node, vmid, action), nil, &resp); err != nil {
		return err
	}

	output.Success(fmt.Sprintf("VM %s %s task queued", vmid, action))

	return nil
}
