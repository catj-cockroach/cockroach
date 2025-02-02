// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package kvprober

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
	"github.com/stretchr/testify/require"
)

func TestReadProbe(t *testing.T) {
	ctx := context.Background()

	t.Run("disabled by default", func(t *testing.T) {
		m := &mock{
			t:      t,
			noPlan: true,
			noGet:  true,
		}
		p := initTestProber(m)

		p.readProbeImpl(ctx, m, m)

		require.Zero(t, p.Metrics().ProbePlanAttempts.Count())
		require.Zero(t, p.Metrics().ReadProbeAttempts.Count())
		require.Zero(t, p.Metrics().ProbePlanFailures.Count())
		require.Zero(t, p.Metrics().ReadProbeFailures.Count())
	})

	t.Run("happy path", func(t *testing.T) {
		m := &mock{t: t}
		p := initTestProber(m)
		readEnabled.Override(ctx, &p.settings.SV, true)

		p.readProbeImpl(ctx, m, m)

		require.Equal(t, int64(1), p.Metrics().ProbePlanAttempts.Count())
		require.Equal(t, int64(1), p.Metrics().ReadProbeAttempts.Count())
		require.Zero(t, p.Metrics().ProbePlanFailures.Count())
		require.Zero(t, p.Metrics().ReadProbeFailures.Count())
	})

	t.Run("planning fails", func(t *testing.T) {
		m := &mock{
			t:       t,
			planErr: fmt.Errorf("inject plan failure"),
			noGet:   true,
		}
		p := initTestProber(m)
		readEnabled.Override(ctx, &p.settings.SV, true)

		p.readProbeImpl(ctx, m, m)

		require.Equal(t, int64(1), p.Metrics().ProbePlanAttempts.Count())
		require.Zero(t, p.Metrics().ReadProbeAttempts.Count())
		require.Equal(t, int64(1), p.Metrics().ProbePlanFailures.Count())
		require.Zero(t, p.Metrics().ReadProbeFailures.Count())
	})

	t.Run("get fails", func(t *testing.T) {
		m := &mock{
			t:      t,
			getErr: fmt.Errorf("inject get failure"),
		}
		p := initTestProber(m)
		readEnabled.Override(ctx, &p.settings.SV, true)

		p.readProbeImpl(ctx, m, m)

		require.Equal(t, int64(1), p.Metrics().ProbePlanAttempts.Count())
		require.Equal(t, int64(1), p.Metrics().ReadProbeAttempts.Count())
		require.Zero(t, p.Metrics().ProbePlanFailures.Count())
		require.Equal(t, int64(1), p.Metrics().ReadProbeFailures.Count())
	})
}

func TestWriteProbe(t *testing.T) {
	ctx := context.Background()

	t.Run("disabled by default", func(t *testing.T) {
		m := &mock{
			t:      t,
			noPlan: true,
			noGet:  true,
		}
		p := initTestProber(m)

		p.writeProbeImpl(ctx, m, m)

		require.Zero(t, p.Metrics().ProbePlanAttempts.Count())
		require.Zero(t, p.Metrics().WriteProbeAttempts.Count())
		require.Zero(t, p.Metrics().ProbePlanFailures.Count())
		require.Zero(t, p.Metrics().WriteProbeFailures.Count())
	})

	t.Run("happy path", func(t *testing.T) {
		m := &mock{t: t}
		p := initTestProber(m)
		writeEnabled.Override(ctx, &p.settings.SV, true)

		p.writeProbeImpl(ctx, m, m)

		require.Equal(t, int64(1), p.Metrics().ProbePlanAttempts.Count())
		require.Equal(t, int64(1), p.Metrics().WriteProbeAttempts.Count())
		require.Zero(t, p.Metrics().ProbePlanFailures.Count())
		require.Zero(t, p.Metrics().WriteProbeFailures.Count())
	})

	t.Run("planning fails", func(t *testing.T) {
		m := &mock{
			t:       t,
			planErr: fmt.Errorf("inject plan failure"),
			noGet:   true,
		}
		p := initTestProber(m)
		writeEnabled.Override(ctx, &p.settings.SV, true)

		p.writeProbeImpl(ctx, m, m)

		require.Equal(t, int64(1), p.Metrics().ProbePlanAttempts.Count())
		require.Zero(t, p.Metrics().WriteProbeAttempts.Count())
		require.Equal(t, int64(1), p.Metrics().ProbePlanFailures.Count())
		require.Zero(t, p.Metrics().WriteProbeFailures.Count())
	})

	t.Run("open txn fails", func(t *testing.T) {
		m := &mock{
			t:      t,
			txnErr: fmt.Errorf("inject txn failure"),
		}
		p := initTestProber(m)
		writeEnabled.Override(ctx, &p.settings.SV, true)

		p.writeProbeImpl(ctx, m, m)

		require.Equal(t, int64(1), p.Metrics().ProbePlanAttempts.Count())
		require.Equal(t, int64(1), p.Metrics().WriteProbeAttempts.Count())
		require.Zero(t, p.Metrics().ProbePlanFailures.Count())
		require.Equal(t, int64(1), p.Metrics().WriteProbeFailures.Count())
	})
	// TODO(josh): Add cases for put & del failures.
}

func initTestProber(m *mock) *Prober {
	p := NewProber(Opts{
		Tracer:                  tracing.NewTracer(),
		HistogramWindowInterval: time.Minute, // actual value not important to test
		Settings:                cluster.MakeTestingClusterSettings(),
	})
	p.readPlanner = m
	return p
}

type mock struct {
	t *testing.T

	noPlan  bool
	planErr error

	noGet  bool
	getErr error
	txnErr error
}

func (m *mock) next(ctx context.Context) (Step, error) {
	if m.noPlan {
		m.t.Errorf("plan call made but not expected")
	}
	return Step{}, m.planErr
}

func (m *mock) Get(ctx context.Context, key interface{}) (kv.KeyValue, error) {
	if m.noGet {
		m.t.Errorf("get call made but not expected")
	}
	return kv.KeyValue{}, m.getErr
}

func (m *mock) Txn(ctx context.Context, f func(ctx context.Context, txn *kv.Txn) error) error {
	return m.txnErr
}
