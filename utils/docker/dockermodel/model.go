// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
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

package dockermodel

type DockerImage struct {
	// containers
	// Required: true
	Containers int64 `json:"Containers"`

	// created
	// Required: true
	Created int64 `json:"Created"`

	// Id
	// Required: true
	ID string `json:"Id"`

	// labels
	// Required: true
	Labels map[string]string `json:"Labels"`

	// parent Id
	// Required: true
	ParentID string `json:"ParentId"`

	// repo digests
	// Required: true
	RepoDigests []string `json:"RepoDigests"`

	// repo tags
	// Required: true
	RepoTags []string `json:"RepoTags"`

	RepoTag string `json:"RepoTag"`

	// shared size
	// Required: true
	SharedSize int64 `json:"SharedSize"`

	// size
	// Required: true
	Size int64 `json:"Size"`

	// virtual size
	// Required: true
	VirtualSize int64 `json:"VirtualSize"`
}

type DockerContainer struct {
	ImageID    string
	Image      string
	ID         string
	Ports      []string
	State      string
	HostConfig map[string]string
	Mounts     []string
	Names      []string
	Created    int64
}
