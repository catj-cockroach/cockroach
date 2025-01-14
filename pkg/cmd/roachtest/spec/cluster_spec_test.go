// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package spec

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cockroachdb/datadriven"
	"github.com/stretchr/testify/require"
)

func TestClusterSpec_Args(t *testing.T) {
	const nodeCount = 12
	const instanceType = ""
	var opts []Option
	opts = append(opts, Geo())
	s := MakeClusterSpec(AWS, instanceType, nodeCount, opts...)
	t.Logf("%#v", s)
	// Regression test against a bug in which we would request an SSD machine type
	// together with the --local-ssd=false option.
	args, err := s.Args()
	require.NoError(t, err)
	act := fmt.Sprint(args)
	require.NotContains(t, act, "--aws-machine-type-ssd")
	require.Contains(t, act, "--local-ssd=false")
}

func TestClusterSpec_Args_DataDriven(t *testing.T) {
	// The following list of specs was obtained by collecting all of the specs
	// registered by tests (at the time of writing) both under `--cloud=aws` and
	// `--cloud=gce`.
	specs := []ClusterSpec{
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      3,
			CPUs:           2,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           2,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-east1-b,us-east1-b,us-west1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-central1-a,us-central1-b,us-central1-c",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           32,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyTagged{Tag: "offset-injector"},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      32,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyTagged{Tag: "jepsen"},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      2,
			CPUs:           32,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      5,
			CPUs:           96,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     2500,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-central1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      8,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "europe-west2-b,europe-west4-b,asia-northeast1-b,us-west1-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      33,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      3,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      10,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      5,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      11,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           1,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           32,
			SSDs:           2,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      10,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           96,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           4,
			SSDs:           1,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      2,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      8,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      5,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      13,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      3,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      12,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      21,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      2,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      10,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyNone{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      21,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyNone{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      8,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "europe-west2-b,europe-west4-b,asia-northeast1-b,us-west1-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      33,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      10,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           32,
			SSDs:           2,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      2,
			CPUs:           32,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      5,
			CPUs:           96,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      2,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           1,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      11,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      5,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-central1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      1,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyTagged{Tag: "offset-injector"},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      3,
			CPUs:           2,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      32,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyTagged{Tag: "jepsen"},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      2,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      13,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-central1-a,us-central1-b,us-central1-c",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      8,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      5,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      3,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           96,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           2,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-east1-b,us-east1-b,us-west1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           32,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      12,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east1-b,us-west1-b,europe-west2-b",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      6,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     2500,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      3,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east-2b,us-west-1a,eu-west-1a",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      10,
			CPUs:           4,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      10,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "us-east-2b,us-west-1a,eu-west-1a",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      9,
			CPUs:           8,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           4,
			SSDs:           1,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
	}
	n := len(specs)
	// Append a copy of the slice in which PreferLocalSSD is always false
	// (it's always true above).
	for i := 0; i < n; i++ {
		s := specs[i]
		s.PreferLocalSSD = false
		specs = append(specs, s)
	}
	// Ditto but cloud is now local.
	for i := 0; i < n; i++ {
		s := specs[i]
		if s.Cloud != GCE {
			continue
		}
		s.Cloud = Local
		specs = append(specs, s)
	}

	zfsSpecs := []ClusterSpec{
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
			FileSystem:     Zfs,
		},
		// The following spec should error out, cause we
		// don't support node creation with zfs on aws, yet.
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
			FileSystem:     Zfs,
		},
		// The following spec should error out, cause we
		// don't support node creation with zfs on azure, yet.
		{
			Cloud:          Azure,
			InstanceType:   "",
			NodeCount:      7,
			CPUs:           16,
			SSDs:           0,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            true,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
			FileSystem:     Zfs,
		},
	}

	specs = append(specs, zfsSpecs...)

	// RAID 0 is off by default, unless explicitly asked for in a roachtest.
	raidSpecs := []ClusterSpec{
		{
			Cloud:          GCE,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           4,
			SSDs:           1,
			RAID0:          true,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
		// AWS roachtests do not currently support local SSDs, so this config
		// should result in an error.
		{
			Cloud:          AWS,
			InstanceType:   "",
			NodeCount:      4,
			CPUs:           4,
			SSDs:           1,
			RAID0:          true,
			VolumeSize:     0,
			PreferLocalSSD: true,
			Zones:          "",
			Geo:            false,
			Lifetime:       43200000000000,
			ReusePolicy:    ReusePolicyAny{},
		},
	}
	specs = append(specs, raidSpecs...)

	path := filepath.Join("testdata", "collected_specs.txt")
	datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
		if td.Cmd != "print-args-for-all" {
			t.Fatalf("unsupported command %s", td.Cmd)
		}
		var buf strings.Builder
		for i, s := range specs {
			pos := i + 1
			args, err := s.Args()
			fmt.Fprintf(&buf, "%d: %#v\n  ", pos, s)
			if err != nil {
				fmt.Fprintln(&buf, err)
			} else {
				fmt.Fprintln(&buf, strings.Join(args, " "))
			}
		}
		return buf.String()
	})
}
