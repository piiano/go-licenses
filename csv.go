// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/piiano/go-licenses/licenses"
)

const (
	OutputFormatTable         = "table"
	OutputFormatMarkdownTable = "md-table"
)

var (
	csvHelp = "Prints all licenses that apply to one or more Go packages and their dependencies."
	csvCmd  = &cobra.Command{
		Use:   "csv <package> [package...]",
		Short: csvHelp,
		Long:  csvHelp + packageHelp,
		Args:  cobra.MinimumNArgs(1),
		RunE:  csvMain,
	}

	gitRemotes    []string
	tableFormat   bool
	mdTableFormat bool
)

func init() {
	csvCmd.Flags().StringArrayVar(&gitRemotes, "git_remote", []string{"origin", "upstream"}, "Remote Git repositories to try")
	csvCmd.Flags().BoolVar(&tableFormat, OutputFormatTable, false, "Whether the output format should be table")
	csvCmd.Flags().BoolVar(&mdTableFormat, OutputFormatMarkdownTable, false, "Whether the output format should be table")

	csvCmd.MarkFlagsMutuallyExclusive(OutputFormatTable, OutputFormatMarkdownTable)

	rootCmd.AddCommand(csvCmd)
}

func csvMain(_ *cobra.Command, args []string) error {
	header := []string{"name", "license_url", "license_name"}

	var writer writer = NewCSVWriter(os.Stdout)
	if tableFormat {
		writer = NewTableWriter(os.Stdout, header)
	} else if mdTableFormat {
		writer = NewMarkdownTableWriter(os.Stdout, header)
	}

	classifier, err := licenses.NewClassifier(confidenceThreshold)
	if err != nil {
		return err
	}

	libs, err := licenses.Libraries(context.Background(), classifier, ignore, args...)
	if err != nil {
		return err
	}
	for _, lib := range libs {
		licenseURL := "Unknown"
		licenseName := "Unknown"
		if lib.LicensePath != "" {
			name, _, err := classifier.Identify(lib.LicensePath)
			if err == nil {
				licenseName = name
			} else {
				glog.Errorf("Error identifying license in %q: %v", lib.LicensePath, err)
			}
			url, err := lib.FileURL(context.Background(), lib.LicensePath)
			if err == nil {
				licenseURL = url
			} else {
				glog.Warningf("Error discovering license URL: %s", err)
			}
		}
		if err := writer.Write(lib.Name(), licenseURL, licenseName); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
