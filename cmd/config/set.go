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
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/dcjulian29/proxmoxctl/internal/color"
	"github.com/dcjulian29/proxmoxctl/internal/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func setCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "interactively set server connection settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print(color.Green("Server name (friendly label): "))
			serverName, _ := reader.ReadString('\n')
			serverName = strings.TrimSpace(serverName)

			fmt.Print(color.Green("Server URL (e.g. https://192.168.1.10:8006): "))
			serverURL, _ := reader.ReadString('\n')
			serverURL = strings.TrimSpace(serverURL)

			fmt.Print(color.Green("API Token (USER@REALM!TOKENID=SECRET): "))
			tokenBytes, err := term.ReadPassword(int(syscall.Stdin))

			if err != nil {
				// Fallback for non-TTY environments
				apiToken, _ := reader.ReadString('\n')
				tokenBytes = []byte(strings.TrimSpace(apiToken))
			}

			fmt.Println()

			apiToken := strings.TrimSpace(string(tokenBytes))

			viper.Set(settings.KeyServerName, serverName)
			viper.Set(settings.KeyServerURL, serverURL)
			viper.Set(settings.KeyAPIToken, apiToken)

			if err := settings.Save(); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			return nil
		},
	}
}
