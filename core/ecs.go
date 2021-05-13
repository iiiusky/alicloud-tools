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

package core

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/iiiusky/alicloud-tools/common"
	"strconv"
	"strings"
)

// 获取指定区域的所有实例列表
func GetRegionInstances(regionId string) (instances []ecs.Instance) {
	client, err := ecs.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取指定区域的所有实例列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return
	}

	for i := 1; i < 100; i++ {
		request := ecs.CreateDescribeInstancesRequest()
		request.Scheme = "https"
		request.PageSize = requests.Integer("100")
		request.PageNumber = requests.Integer(strconv.Itoa(i))

		response, err := client.DescribeInstances(request)
		if err != nil {
			common.Logger().Error(fmt.Sprintf("【获取指定区域的所有实例列表】创建获取实例请求发生异常,异常信息为 %s", err.Error()))
			return
		}

		instances = append(instances, response.Instances.Instance...)

		if len(response.Instances.Instance) == 0 {
			break
		}
	}

	return instances
}

// 获取所有实例
func GetAllInstances(regionId string, printInfo bool) (instances []ecs.Instance) {
	for _, region := range common.ECSRegions {
		if regionId != "" && regionId != region.RegionId {
			continue
		}

		if printInfo {
			fmt.Printf("正在扫描  %s 区域的 ECS\t", region.LocalName)
		}

		regionInstances := GetRegionInstances(region.RegionId)
		instances = append(instances, regionInstances...)

		if printInfo {
			fmt.Printf("扫描到 %d 台 ECS\n", len(regionInstances))
		}
	}

	return instances
}

// 查询单个实例
func QuerySingleInstance(regionId string, instanceId string) (instances ecs.Instance) {
	if regionId == "" {
		for _, region := range common.ECSRegions {
			regionInstances := GetRegionInstances(region.RegionId)
			for _, instance := range regionInstances {
				if instance.InstanceId == instanceId {
					return instance
				}
			}
		}
	} else {
		regionInstances := GetRegionInstances(regionId)
		for _, instance := range regionInstances {
			if instance.InstanceId == instanceId {
				return instance
			}
		}
	}
	return instances
}

// 执行命令
func EcsRunCommand(regionId, scriptType, commandContent string, instanceId string) bool {
	client, err := ecs.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【执行命令】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	request := ecs.CreateRunCommandRequest()
	request.Scheme = "https"
	request.Type = scriptType
	request.CommandContent = commandContent
	request.InstanceId = &[]string{instanceId}

	response, err := client.RunCommand(request)

	if err != nil {
		common.Logger().Error(fmt.Sprintf("【执行命令】创建执行命令请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	return response.IsSuccess()
}

// 检测云助手安装情况
func CheckCloudAssistantStatus(regionId, instanceId string) bool {
	client, err := ecs.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【检测云助手安装情况】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	request := ecs.CreateDescribeCloudAssistantStatusRequest()
	request.Scheme = "https"
	request.InstanceId = &[]string{instanceId}

	response, err := client.DescribeCloudAssistantStatus(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【检测云助手安装情况】创建检测云助手安装情况请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	if strings.ToLower(response.InstanceCloudAssistantStatusSet.InstanceCloudAssistantStatus[0].CloudAssistantStatus) != "true" {
		return false
	}

	return true
}

// 显示传入的实例列表具体信息
func ShowInstancesInfo(instances []ecs.Instance, isRunner bool) {
	for _, instance := range instances {
		if isRunner && instance.Status != "Running" {
			continue
		}

		fmt.Println("++++++++++++++++++++++++++++++++++++++++")
		ShowInstanceInfo(instance)
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
}

// 显示实例具体信息
func ShowInstanceInfo(instance ecs.Instance) {
	fmt.Printf("实例ID: %s \n", instance.InstanceId)
	fmt.Printf("云助手安装情况(未安装的不可以执行命令): %v \n", CheckCloudAssistantStatus(instance.RegionId, instance.InstanceId))
	fmt.Printf("实例名称: %s \n", instance.InstanceName)
	fmt.Printf("实例描述: %s \n", instance.Description)
	fmt.Printf("实例规格: %s \n", instance.InstanceType)
	fmt.Printf("实例状态: %s \n", instance.Status)
	fmt.Printf("实例主机名: %s \n", instance.HostName)
	fmt.Printf("实例VPC ID: %v \n", instance.VpcAttributes.VpcId)
	fmt.Printf("实例地域ID: %s \n", instance.RegionId)
	fmt.Printf("CPU信息: %d 核 \n", instance.Cpu)
	fmt.Printf("内存信息: %d M\n", instance.Memory)
	fmt.Printf("实例创建时间: %s \n", instance.CreationTime)
	fmt.Printf("实例过期时间: %s \n", instance.ExpiredTime)
	fmt.Printf("实例网卡列表: %v \n", instance.NetworkInterfaces.NetworkInterface)
	fmt.Printf("实例公网IP列表: %v \n", instance.PublicIpAddress.IpAddress)
	fmt.Printf("实例弹性公网信息: %v \n", instance.EipAddress.IpAddress)
	fmt.Printf("实例操作系统类型: %s \n", instance.OSType)
	fmt.Printf("实例操作系统名称: %s \n", instance.OSName)

	instanceChargeType := instance.InstanceChargeType
	if instanceChargeType == "PostPaid" {
		fmt.Printf("实例计费方式:按量付费 \n")
	} else {
		fmt.Printf("实例计费方式:包年包月 \n")
	}

	fmt.Printf("实例网络类型: %s \n", instance.InstanceNetworkType)

	fmt.Printf("实例所属安全组列表: %v \n", instance.SecurityGroupIds.SecurityGroupId)
	for _, s := range instance.SecurityGroupIds.SecurityGroupId {
		groupInfo := GetEcsSecurityGroupInfo(instance.RegionId, s)
		var all []string
		var ingress []string
		var egress []string

		for _, permission := range groupInfo.Permissions.Permission {
			if permission.Direction == "egress" {
				egress = append(egress, permission.PortRange)
			} else if permission.Direction == "ingress" {
				ingress = append(ingress, permission.PortRange)
			} else {
				all = append(all, permission.PortRange)
			}
		}
		fmt.Printf("实例所属安全组[%s]端口信息[入方向]: %s \n", groupInfo.SecurityGroupId, ingress)
		fmt.Printf("实例所属安全组[%s]端口信息[出方向]: %s \n", groupInfo.SecurityGroupId, egress)
		fmt.Printf("实例所属安全组[%s]端口信息[不区分方向]: %s \n", groupInfo.SecurityGroupId, all)
	}
}
