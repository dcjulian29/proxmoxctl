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
package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dcjulian29/proxmoxctl/internal/color"
	"github.com/spf13/viper"
)

const (
	KeyServerName = "server_name"
	KeyServerURL  = "server_url"
	KeyAPIToken   = "api_token"
)

func Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine config directory: %w", err)
	}

	dir := filepath.Join(home, ".config", "proxmoxctl")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	configPath := filepath.Join(dir, "config.yml")
	viper.SetConfigFile(configPath)

	if err := viper.WriteConfig(); err != nil {
		if os.IsNotExist(err) {
			return viper.WriteConfigAs(configPath)
		}

		return fmt.Errorf("could not write config: %w", err)
	}

	fmt.Println(color.Info("\n✓ config saved to %s\n", configPath))

	return nil
}

func RequireConfig() error {
	missing := []string{}
	for _, key := range []string{KeyServerURL, KeyAPIToken} {
		if viper.GetString(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config keys: %v — run `proxmoxctl config set` to configure", missing)
	}

	return nil
}
