// Copyright 2021 TiKV Project Authors.
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

package statistics

import (
	"time"

	. "github.com/pingcap/check"
	"github.com/pingcap/kvproto/pkg/metapb"
	"github.com/pingcap/kvproto/pkg/pdpb"
	"github.com/qiaohao9/pd/server/core"
)

var _ = Suite(&testStoreSuite{})

type testStoreSuite struct{}

func (s *testStoreSuite) TestFilterUnhealtyStore(c *C) {
	stats := NewStoresStats()
	cluster := core.NewBasicCluster()
	for i := uint64(1); i <= 5; i++ {
		cluster.PutStore(core.NewStoreInfo(&metapb.Store{Id: i}, core.SetLastHeartbeatTS(time.Now())))
		stats.Observe(i, &pdpb.StoreStats{})
	}
	c.Assert(stats.GetStoresLoads(), HasLen, 5)

	cluster.PutStore(cluster.GetStore(1).Clone(core.SetLastHeartbeatTS(time.Now().Add(-24 * time.Hour))))
	cluster.PutStore(cluster.GetStore(2).Clone(core.TombstoneStore()))
	cluster.DeleteStore(cluster.GetStore(3))

	stats.FilterUnhealthyStore(cluster)
	loads := stats.GetStoresLoads()
	c.Assert(loads, HasLen, 2)
	c.Assert(loads[4], NotNil)
	c.Assert(loads[5], NotNil)
}
