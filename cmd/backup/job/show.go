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
package job

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <jobid>",
		Short: "Show details of a scheduled backup job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.New()
			if err != nil {
				return err
			}

			var resp struct {
				Data map[string]any `json:"data"`
			}

			if err := client.Get(fmt.Sprintf("/cluster/backup/%s", args[0]), &resp); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			d := resp.Data

			headers := []string{"FIELD", "VALUE"}
			rows := [][]string{
				{"Job ID", toString(d["id"])},
				{"VM IDs", toString(d["vmid"])},
				{"Storage", toString(d["storage"])},
				{"Schedule", toString(d["schedule"])},
				{"Mode", toString(d["mode"])},
				{"Compression", toString(d["compress"])},
				{"Enabled", formatBool(d["enabled"])},
				{"Mail To", toString(d["mailto"])},
				{"Max Files", toString(d["maxfiles"])},
				{"Notes Template", toString(d["notes-template"])},
			}

			output.Table(headers, rows)

			return nil
		},
	}
}
