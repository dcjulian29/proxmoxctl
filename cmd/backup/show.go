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
package backup

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	var (
		node    string
		storage string
		file    string
	)

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show details of a specific backup file",
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

			volid := fmt.Sprintf("%s:%s", storage, file)

			var resp struct {
				Data map[string]any `json:"data"`
			}

			if err := client.Get(
				fmt.Sprintf("/nodes/%s/storage/%s/content/%s", node, storage, file),
				&resp,
			); err != nil {
				return err
			}

			if output.IsJSON() {
				return output.JSON(resp.Data)
			}

			d := resp.Data

			headers := []string{"FIELD", "VALUE"}
			rows := [][]string{
				{"Volume ID", volid},
				{"Format", toString(d["format"])},
				{"Size", formatBytes(toFloat(d["size"]))},
				{"Created", formatEpoch(d["ctime"])},
				{"Notes", toString(d["notes"])},
				{"Protected", formatBool(d["protected"])},
				{"Encrypted", formatBool(d["encrypted"])},
			}

			output.Table(headers, rows)

			return nil
		},
	}

	cmd.Flags().StringVar(&node, "node", "", "Proxmox node (auto-detected if not set)")
	cmd.Flags().StringVar(&storage, "storage", "", "Storage containing the backup (required)")
	cmd.Flags().StringVar(&file, "file", "", "Backup filename (required)")

	_ = cmd.MarkFlagRequired("storage")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
