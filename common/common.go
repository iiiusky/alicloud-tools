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
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/bndr/gotabulate"
)

var AccessKey string
var SecretKey string
var STSAccessKey string
var STSSecretKey string
var STSToken string
var UseSTS bool
var Verbose bool
var ECSRegions []ecs.Region
var APPVersion string

// InitEcsRegions 初始化区域信息表
func InitEcsRegions() bool {
	client, err := GetEcsClient("cn-hangzhou")
	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"

	if Verbose {
		requestByte, _ := json.Marshal(request)
		fmt.Println(fmt.Sprintf("\r\n InitEcsRegions request is: %s", string(requestByte)))
	}

	if err != nil {
		Logger().Error(fmt.Sprintf("【初始化区域信息表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	response, err := client.DescribeRegions(request)

	if err != nil {
		Logger().Error(fmt.Sprintf("【初始化区域信息表】创建获取区域信息请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	if Verbose {
		fmt.Println(fmt.Sprintf("\r\n InitEcsRegions response is: %s", response.String()))
	}

	ECSRegions = response.Regions.Region
	return true
}

// GetEcsClient 获取ECS 客户端
func GetEcsClient(regionId string) (*ecs.Client, error) {
	if UseSTS {
		return ecs.NewClientWithStsToken(regionId, STSAccessKey, STSSecretKey, STSToken)
	} else {
		return ecs.NewClientWithAccessKey(regionId, AccessKey, SecretKey)
	}
}

// ShowRegions 显示地域信息
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
