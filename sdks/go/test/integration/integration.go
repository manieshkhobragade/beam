// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package integration provides functionality that needs to be shared between all
// integration tests.
package integration

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/apache/beam/sdks/go/pkg/beam/testing/ptest"
)

// Filters for temporarily skipping integration tests. All filters are regex
// matchers that must match the full name of a test at the point where
// CheckFilters is called. Multiple tests can be skipped by using regex
// wildcards. (ex. "TestXLang_.*" filters all tests starting with TestXLang_)
//
// It is strongly recommended to include, TODOs, Jira issues, or just comments
// describing why tests are being skipped.

// sickbay filters tests that fail due to Go SDK errors. These tests will not
// execute on any runners.
var sickbay = []string{
	// TODO(BEAM-11418): These tests are currently failing with the Go SDK.
	"TestXLang_CoGroupBy",
	"TestXLang_Multi",
	"TestXLang_Partition",
}

// Runner-specific test filters, for features that are not yet supported on
// specific runners.

var directFilters = []string{
	// The direct runner does not yet support cross-language.
	"TestXLang.*",
}

var portableFilters = []string{
	// TODO(BEAM-11499): XLang tests timing out.
	"TestXLang.*",
}

var flinkFilters = []string{
	// TODO(BEAM-11500): Flink tests timing out on reads.
	"TestXLang_Combine.*",
}

var sparkFilters = []string{
	// TODO(BEAM-11498): XLang tests broken with Spark runner.
	"TestXLang.*",
	"TestParDoSideInput",
	"TestParDoKVSideInput",
}

// CheckFilters checks if an integration test is filtered to be skipped, either
// because the intended runner does not support it, or the test is sickbayed.
// This method should be called at the beginning of any integration test. If
// t.Run is used, CheckFilters should be called within the t.Run callback, so
// that sub-tests can be skipped individually.
func CheckFilters(t *testing.T) {
	// Check for sickbaying first.
	n := t.Name()
	for _, f := range sickbay {
		// Add start and end of string regexp matchers so only a full match is
		// counted.
		f = fmt.Sprintf("^%v$", f)
		match, err := regexp.MatchString(f, n)
		if err != nil {
			t.Errorf("Matching of regex '%v' with test '%v' failed: %v", f, n, err)
		}
		if match {
			t.Skipf("Test %v is currently sickbayed on all runners", n)
		}
	}

	// Test for runner-specific skipping second.
	var filters []string
	runner := *ptest.Runner
	if runner == "" {
		runner = ptest.DefaultRunner()
	}
	switch runner {
	case "direct", "DirectRunner":
		filters = directFilters
	case "portable", "PortableRunner":
		filters = portableFilters
	case "flink", "FlinkRunner":
		filters = flinkFilters
	case "spark", "SparkRunner":
		filters = sparkFilters
	default:
		return
	}

	for _, f := range filters {
		// Add start and end of string regexp matchers so only a full match is
		// counted.
		f = fmt.Sprintf("^%v$", f)
		match, err := regexp.MatchString(f, n)
		if err != nil {
			t.Errorf("Matching of regex '%v' with test '%v' failed: %v", f, n, err)
		}
		if match {
			t.Skipf("Test %v is currently filtered for runner %v", n, runner)
		}
	}
}
