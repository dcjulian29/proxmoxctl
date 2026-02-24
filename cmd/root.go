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
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dcjulian29/proxmoxctl/cmd/backup"
	"github.com/dcjulian29/proxmoxctl/cmd/config"
	"github.com/dcjulian29/proxmoxctl/cmd/group"
	"github.com/dcjulian29/proxmoxctl/cmd/snapshot"
	"github.com/dcjulian29/proxmoxctl/cmd/status"
	"github.com/dcjulian29/proxmoxctl/cmd/storage"
	"github.com/dcjulian29/proxmoxctl/cmd/user"
	"github.com/dcjulian29/proxmoxctl/internal/api"
	"github.com/dcjulian29/proxmoxctl/internal/color"
	"github.com/dcjulian29/proxmoxctl/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.szostok.io/version/extension"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "proxmoxctl",
	Short: "A full-featured CLI for managing Proxmox Virtual Environment (PVE) infrastructure",
	Long: `proxmoxctl is a full-featured command-line interface for managing Proxmox
Virtual Environment (PVE) infrastructure via the Proxmox REST API.

RESOURCES
  backup      On-demand and scheduled backups (vzdump), restore, and backup file management
  group       Create and manage Proxmox user groups
  snapshot    Create, list, rollback, and delete snapshots for VMs and containers
  storage     List storage pools, inspect configuration, and browse storage contents
  user        Create and manage Proxmox users, passwords, and group membership

OBSERVABILITY
  status      Cluster health, per-node resource usage, resource inventory, and task history

CONFIGURATION
  config      Set and display connection settings (server URL, username, API token)

All commands support --output table (default) or --output json (-o json) for
scripting and piping. Destructive operations prompt for confirmation unless
--force is passed. The --node flag is optional on all node-scoped commands —
the first available cluster node is used when omitted.

Run 'proxmoxctl config set' to get started.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	rootCmd.AddCommand(
		extension.NewVersionCobraCmd(
			extension.WithUpgradeNotice("dcjulian29", "proxmoxctl"),
		),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "\n"+color.Fatal(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "specify configuration file")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format (table or json)")
	rootCmd.PersistentFlags().Bool("insecure", false, "disable TLS certificate verification (not recommended)")

	cobra.OnInitialize(initConfig)

	if err := viper.BindPFlag(output.KeyOutputFormat, rootCmd.PersistentFlags().Lookup("output")); err != nil {
		fmt.Fprintln(os.Stderr, color.Fatal(err))
		os.Exit(1)
	}

	if err := viper.BindPFlag(api.KeyInsecureSkipVerify, rootCmd.PersistentFlags().Lookup("insecure")); err != nil {
		fmt.Fprintln(os.Stderr, color.Fatal(err))
		os.Exit(1)
	}

	rootCmd.AddCommand(backup.NewCommand())
	rootCmd.AddCommand(config.NewCommand())
	rootCmd.AddCommand(group.NewCommand())
	rootCmd.AddCommand(status.NewCommand())
	rootCmd.AddCommand(snapshot.NewCommand())
	rootCmd.AddCommand(storage.NewCommand())
	rootCmd.AddCommand(user.NewCommand())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, _ := os.UserHomeDir()
		dir := filepath.Join(home, ".config", "proxmoxctl")

		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		viper.AddConfigPath(dir)
	}

	viper.SetEnvPrefix("PROXMOX")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, color.Fatal("error reading config:", err))
			os.Exit(1)
		}
	}
}
