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

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <userid>",
		Short: "Show details of a Proxmox user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var resp struct {
				Data map[string]any `json:"data"`
			}

			if err := client.Get(fmt.Sprintf("/access/users/%s", args[0]), &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			d := resp.Data

			headers := []string{"FIELD", "VALUE"}
			rows := [][]string{
				{"User ID", toString(d["userid"])},
				{"First Name", toString(d["firstname"])},
				{"Last Name", toString(d["lastname"])},
				{"Email", toString(d["email"])},
				{"Comment", toString(d["comment"])},
				{"Enabled", formatBool(d["enable"])},
				{"Expire", formatExpire(d["expire"])},
				{"Groups", toString(d["groups"])},
				{"Keys", toString(d["keys"])},
			}

			output.Table(headers, rows)

			return nil
		},
	}
}
