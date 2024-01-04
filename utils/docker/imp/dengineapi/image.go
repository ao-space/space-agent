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

package dengineapi

import (
	"agent/utils/docker/dockermodel"
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func PullImageByUsernameAndPasswd(cli *client.Client, imageName string, username, password string) error {
	authStr, err := GenAuthString(username, password)
	if err != nil {
		return err
	}

	return PullImageByAuthStr(cli, imageName, authStr)
}

func PullImageWithoutAuth(cli *client.Client, imageName string) error {
	return PullImageByAuthStr(cli, imageName, "")
}

func PullImageByAuthStr(cli *client.Client, imageName string, authStr string) error {

	fmt.Printf(">>>> imageName: %v, authStr=%v \n", imageName, authStr)

	events, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{
		All:           false,
		RegistryAuth:  authStr,
		PrivilegeFunc: nil,
	})

	if err != nil {
		return err
	}
	defer events.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(events)
	s := buf.String()
	fmt.Printf("#### ImagePull return: %v\n", s)
	// fmt.Println("image pull success")
	return nil
}

func ListImages(cli *client.Client) ([]*dockermodel.DockerImage, error) {
	var err error
	if cli == nil {
		cli, err = NewClient()
		if err != nil {
			return nil, fmt.Errorf("failed NewClient, err:%v", err)
		}
		defer cli.Close()
	}

	imags, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed ImageList, err:%v", err)
	}

	ret := []*dockermodel.DockerImage{}
	for _, img := range imags {
		ret = append(ret, &dockermodel.DockerImage{
			Containers:  img.Containers,
			Created:     img.Created,
			Labels:      img.Labels,
			ParentID:    img.ParentID,
			RepoDigests: img.RepoDigests,
			RepoTags:    img.RepoTags,
			RepoTag:     img.RepoTags[0],
			SharedSize:  img.SharedSize,
			ID:          img.ID,
			Size:        img.Size,
			VirtualSize: img.VirtualSize})
	}

	return ret, nil
}

func RemoveImage(cli *client.Client, imageId string, removeOpt types.ImageRemoveOptions) error {
	var err error
	if cli == nil {
		cli, err = NewClient()
		if err != nil {
			return fmt.Errorf("failed NewClient, err:%v", err)
		}
		defer cli.Close()
	}
	_, err = cli.ImageRemove(context.Background(), imageId, removeOpt)
	if err != nil {
		return err
	}
	return nil
}
