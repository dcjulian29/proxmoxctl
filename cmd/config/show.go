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
package config

import (
	"fmt"

	"github.com/dcjulian29/proxmoxctl/internal/settings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show current configuration (token is masked)",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := viper.GetString(settings.KeyServerName)
			serverURL := viper.GetString(settings.KeyServerURL)
			token := viper.GetString(settings.KeyAPIToken)

			if serverURL == "" {
				err := fmt.Errorf("no configuration found. Run `proxmoxctl config set` to get started")
				return err
			}

			masked := "****"
			if len(token) > 8 {
				masked = token[:4] + "****" + token[len(token)-4:]
			}

			fmt.Printf("Server Name : %s\n", serverName)
			fmt.Printf("Server URL  : %s\n", serverURL)
			fmt.Printf("API Token   : %s\n", masked)

			return nil
		},
	}
}
