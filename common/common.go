/*
Copyright © 2020 iiusky sky@03sec.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/bndr/gotabulate"
)

var AccessKey string
var SecretKey string
var ECSRegions []ecs.Region

// 初始化区域信息表
func InitEcsRegions() bool {
	client, err := ecs.NewClientWithAccessKey("cn-hangzhou", AccessKey, SecretKey)
	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"

	if err != nil {
		Logger().Error(fmt.Sprintf("【初始化区域信息表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	response, err := client.DescribeRegions(request)
	if err != nil {
		Logger().Error(fmt.Sprintf("【初始化区域信息表】创建获取区域信息请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	ECSRegions = response.Regions.Region
	return true
}

// 显示地域信息
func ShowRegions() {
	var dates [][]string
	count := 0

	for _, region := range ECSRegions {
		data := []string{fmt.Sprintf("#%d", count+1), region.LocalName, region.RegionId}
		dates = append(dates, data)
		count = count + 1
	}

	t := gotabulate.Create(dates)
	t.SetHeaders([]string{"#", "名称", "区域ID"})
	t.SetAlign("center")

	fmt.Println(t.Render("grid"))
}
