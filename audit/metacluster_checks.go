// Copyright 2024 The NATS Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package audit

import (
	"errors"
	"fmt"
	"github.com/nats-io/jsm.go/api"
	"github.com/nats-io/jsm.go/audit/archive"
	"github.com/nats-io/nats-server/v2/server"
)

func RegisterMetaChecks(collection *CheckCollection) error {
	return collection.Register(
		Check{
			Code:        "META_001",
			Suite:       "meta",
			Name:        "Meta cluster offline replicas",
			Description: "All nodes part of the meta group are online",
			Handler:     checkMetaClusterOfflineReplicas,
		},
		Check{
			Code:        "META_002",
			Suite:       "meta",
			Name:        "Meta cluster leader",
			Description: "All nodes part of the meta group agree on the meta cluster leader",
			Handler:     checkMetaClusterLeader,
		},
	)
}

// checkMetaClusterLeader verify that all server agree on the same meta group leader in each known cluster
func checkMetaClusterLeader(_ *Check, r *archive.Reader, examples *ExamplesCollection, log api.Logger) (Outcome, error) {
	jsTag := archive.TagServerJetStream()

	for _, clusterName := range r.ClusterNames() {
		clusterTag := archive.TagCluster(clusterName)
		leaderFollowers := make(map[string][]string)

		for _, serverName := range r.ClusterServerNames(clusterName) {
			serverTag := archive.TagServer(serverName)

			var resp server.ServerAPIJszResponse
			var jetStreamInfo *server.JSInfo
			err := r.Load(&resp, clusterTag, serverTag, jsTag)
			if errors.Is(err, archive.ErrNoMatches) {
				log.Warnf("Artifact 'JSZ' is missing for server %s cluster %s", serverName, clusterName)
				continue
			} else if err != nil {
				return Skipped, fmt.Errorf("failed to load JSZ for server %s: %w", serverName, err)
			}
			jetStreamInfo = resp.Data

			if jetStreamInfo.Disabled {
				continue
			}

			if jetStreamInfo.Meta == nil {
				log.Warnf("%s / %s does not have meta group info", clusterName, serverName)
				continue
			}

			leader := jetStreamInfo.Meta.Leader
			if leader == "" {
				leader = "NO_LEADER"
			}

			_, present := leaderFollowers[leader]
			if !present {
				leaderFollowers[leader] = make([]string, 0)
			}
			leaderFollowers[leader] = append(leaderFollowers[leader], serverName)
		}

		if len(leaderFollowers) > 1 {
			examples.Add("Members of %s disagree on meta leader (%v)", clusterName, leaderFollowers)
		}
	}

	if examples.Count() > 0 {
		log.Errorf("Found %d instance of replicas disagreeing on meta-cluster leader", examples.Count())
		return Fail, nil
	}

	return Pass, nil
}

// checkMetaClusterOfflineReplicas verify that all meta-cluster replicas are online for each known cluster
func checkMetaClusterOfflineReplicas(_ *Check, r *archive.Reader, examples *ExamplesCollection, log api.Logger) (Outcome, error) {
	_, err := r.EachClusterServerJsz(func(clusterTag *archive.Tag, serverTag *archive.Tag, err error, jsz *server.ServerAPIJszResponse) error {
		if errors.Is(err, archive.ErrNoMatches) {
			log.Warnf("Artifact 'JSZ' is missing for server %s", serverTag)
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to load variables for server %s: %w", serverTag, err)
		}

		if jsz.Data.Disabled {
			return nil
		}

		if jsz.Data.Meta == nil {
			log.Warnf("%s / %s does not have meta group info", clusterTag, serverTag)
			return nil
		}

		for _, peerInfo := range jsz.Data.Meta.Replicas {
			if peerInfo.Offline {
				examples.Add("%s - %s reports peer %s as offline", clusterTag, serverTag, peerInfo.Name)
			}
		}

		return nil
	})
	if err != nil {
		return Skipped, err
	}

	if examples.Count() > 0 {
		log.Errorf("Found %d instance of replicas marked offline by peers", examples.Count())
		return Fail, nil
	}

	return Pass, nil
}
