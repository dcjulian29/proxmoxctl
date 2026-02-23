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
package group

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func modifyCmd() *cobra.Command {
	var comment string
	cmd := &cobra.Command{
		Use:   "modify <groupid>",
		Short: "Modify a user group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			payload := map[string]any{}

			if comment != "" {
				payload["comment"] = comment
			}

			if len(payload) == 0 {
				return fmt.Errorf("no changes specified — use --comment")
			}

			if err := client.Put(fmt.Sprintf("/access/groups/%s", args[0]), payload, nil); err != nil {
				return err
			}

			output.Success(fmt.Sprintf("Group '%s' updated", args[0]))

			return nil
		},
	}

	cmd.Flags().StringVar(&comment, "comment", "", "New description for the group")

	return cmd
}
