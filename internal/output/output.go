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
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dcjulian29/proxmoxctl/internal/color"
	"github.com/spf13/viper"
)

const KeyOutputFormat = "output_format"

func Format() string {
	f := viper.GetString(KeyOutputFormat)

	if f == "" {
		return "table"
	}

	return strings.ToLower(f)
}

func IsJSON() bool {
	return Format() == "json"
}

func JSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")

	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	fmt.Println(string(b))

	return nil
}

func Table(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(w, strings.Join(headers, "\t"))
	_, _ = fmt.Fprintln(w, strings.Repeat("-", 60))

	for _, row := range rows {
		_, _ = fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	_ = w.Flush()
}

func Success(msg string) {
	if IsJSON() {
		b, _ := json.Marshal(map[string]string{"status": "ok", "message": msg})
		fmt.Println(string(b))
		return
	}

	fmt.Println(color.Green("✓ ") + msg)
}

func Aborted(reason string) {
	if IsJSON() {
		b, _ := json.Marshal(map[string]string{"status": "aborted", "message": reason})
		fmt.Println(string(b))
		return
	}

	fmt.Println(color.Red("⨯ ") + reason)
}
