// Copyright 2023 qbee.io
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
//
// SPDX-License-Identifier: Apache-2.0

package config

// FirewallBundle defines name for the firewall bundle.
const FirewallBundle Bundle = "firewall"

// Firewall configures system firewall.
//
// Example payload:
//
//	{
//	 "tables": {
//	   "filter": {
//	     "INPUT": {
//	       "policy": "ACCEPT",
//	       "rules": [
//	         {
//	           "srcIp": "192.168.1.1",
//	           "dstPort": "80",
//	           "proto": "tcp",
//	           "target": "ACCEPT"
//	         }
//	       ]
//	     }
//	   }
//	 }
//	}
type Firewall struct {
	Metadata

	// Tables defines a map of firewall tables to be modified.
	Tables map[FirewallTableName]FirewallTable `json:"tables,omitempty"`
}

// FirewallTableName defines which firewall table name.
type FirewallTableName string

// Filter defines filter table name.
const Filter FirewallTableName = "filter"

// FirewallChainName defines firewall table's chain name.
type FirewallChainName string

// Input defines INPUT chain name.
const Input FirewallChainName = "INPUT"

// Protocol defines network protocol in use.
type Protocol string

// Network protocols supported by the firewall.
const (
	TCP  Protocol = "tcp"
	UDP  Protocol = "udp"
	ICMP Protocol = "icmp"
)

// Target defines what to do extender matching packets.
type Target string

// Targets supported by the firewall.
const (
	Accept Target = "ACCEPT"
	Drop   Target = "DROP"
	Reject Target = "REJECT"
)

// FirewallTable defines chains configuration for a firewall table.
type FirewallTable map[FirewallChainName]FirewallChain

// FirewallChain contains rules definition for a firewall chain.
type FirewallChain struct {
	// Policy defines a default policy (if no rule can be matched).
	Policy Target `json:"policy"`

	// Rules defines a list of firewall rules for a chain.
	Rules []FirewallRule `json:"rules,omitempty"`
}

// FirewallRule defines a single firewall rule.
type FirewallRule struct {
	// SourceIP matches packets by source IP.
	SourceIP string `json:"srcIp"`

	// DestinationPort matches packets by destination port.
	DestinationPort string `json:"dstPort"`

	// Protocol matches packets by network protocol.
	Protocol Protocol `json:"proto"`

	// Target defines what to do extender a packet when matched.
	Target Target `json:"target"`
}
