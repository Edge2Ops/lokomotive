// Copyright 2020 The Lokomotive Authors
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

package dns

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"sort"

	"github.com/kinvolk/lokomotive/pkg/terraform"
	"github.com/pkg/errors"
)

const (
	// DNSManual represents manual DNS configuration.
	DNSManual = "manual"
	// DNSRoute53 represents DNS managed in Route 53.
	DNSRoute53 = "route53"
	// DNSCloudflare represents DNS managed in Cloudflare.
	DNSCloudflare = "cloudflare"
)

type dnsEntry struct {
	Name      string   `json:"name"`
	Ttl       int      `json:"ttl"`
	EntryType string   `json:"type"`
	Records   []string `json:"records"`
}

// Validate ensures the DNS provider p is a valid provider.
func Validate(p string) error {
	switch p {
	case DNSManual:
		return nil
	case DNSRoute53:
		return nil
	case DNSCloudflare:
		return nil
	}

	return fmt.Errorf("invalid DNS provider %q", p)
}

// AskToConfigure reads the required DNS entries from a Terraform output,
// asks the user to configure them and checks if the configuration is correct.
func AskToConfigure(ex *terraform.Executor, zone string) error {
	dnsEntries, err := readDNSEntries(ex)
	if err != nil {
		return err
	}

	fmt.Printf("Please configure the following DNS entries at the DNS provider which hosts %q:\n", zone)
	prettyPrintDNSEntries(dnsEntries)

	for {
		fmt.Printf("Press Enter to check the entries or type \"skip\" to continue the installation: ")

		var input string
		fmt.Scanln(&input)

		if input == "skip" {
			break
		} else if input != "" {
			continue
		}

		if checkDNSEntries(dnsEntries) {
			break
		}

		fmt.Println("Entries are not correctly configured, please verify.")
	}

	return nil
}

func readDNSEntries(ex *terraform.Executor) ([]dnsEntry, error) {
	output, err := ex.ExecuteSync("output", "-json", "dns_entries")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get DNS entries")
	}

	var entries []dnsEntry

	if err := json.Unmarshal(output, &entries); err != nil {
		return nil, errors.Wrap(err, "failed to parse DNS entries file")
	}

	return entries, nil
}

func prettyPrintDNSEntries(entries []dnsEntry) {
	fmt.Println("------------------------------------------------------------------------")

	for _, entry := range entries {
		fmt.Printf("Name: %s\n", entry.Name)
		fmt.Printf("Type: %s\n", entry.EntryType)
		fmt.Printf("Ttl: %d\n", entry.Ttl)
		fmt.Printf("Records:\n")
		for _, record := range entry.Records {
			fmt.Printf("- %s\n", record)
		}
		fmt.Println("------------------------------------------------------------------------")
	}
}

func checkDNSEntries(entries []dnsEntry) bool {
	for _, entry := range entries {
		ips, err := net.LookupIP(entry.Name)
		if err != nil {
			return false
		}

		var ipsString []string
		for _, ip := range ips {
			ipsString = append(ipsString, ip.String())
		}

		sort.Strings(ipsString)
		sort.Strings(entry.Records)
		if !reflect.DeepEqual(ipsString, entry.Records) {
			return false
		}
	}

	return true
}
