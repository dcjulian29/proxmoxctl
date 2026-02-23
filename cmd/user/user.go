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
package user

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage Proxmox users",
		Long: `Manage Proxmox users via the /access/users API.

User IDs in Proxmox are always in the form USER@REALM (e.g. alice@pam, bob@pve).

Examples:
  proxmoxctl user list
  proxmoxctl user show alice@pam
  proxmoxctl user create alice@pam --password secret --firstname Alice --lastname Smith
  proxmoxctl user modify alice@pam --email alice@example.com --groups devs,ops
  proxmoxctl user delete alice@pam
  proxmoxctl user passwd alice@pam`,
	}

	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(groupCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(modifyCmd())
	cmd.AddCommand(passwordCmd())
	cmd.AddCommand(showCmd())

	return cmd
}

func formatBool(v interface{}) string {
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

	return fmt.Sprintf("%v", v)
}

func formatExpire(v interface{}) string {
	if v == nil {
		return "never"
	}

	switch val := v.(type) {
	case float64:
		if val == 0 {
			return "never"
		}

		return fmt.Sprintf("%d", int64(val))
	}

	return "never"
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}

	return fmt.Sprintf("%v", v)
}
