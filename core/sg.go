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
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/bndr/gotabulate"
	"github.com/iiiusky/alicloud-tools/common"
)

// GetEcsSecurityGroupInfo 获取安全组信息
func GetEcsSecurityGroupInfo(regionId, securityGroupId string) ecs.DescribeSecurityGroupAttributeResponse {
	client, err := common.GetEcsClient(regionId)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取安全组信息】创建客户端发生异常,异常信息为 %s", err.Error()))
		return ecs.DescribeSecurityGroupAttributeResponse{}
	}

	request := ecs.CreateDescribeSecurityGroupAttributeRequest()
	request.Scheme = "https"
	request.SecurityGroupId = securityGroupId

	response, err := client.DescribeSecurityGroupAttribute(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取安全组信息】创建获取安全组信息请求发生异常,异常信息为 %s", err.Error()))
		return ecs.DescribeSecurityGroupAttributeResponse{}
	} else {
		return *response
	}
}

// AddSecurityGroupPolicy 添加指定安全组ID的端口策略
func AddSecurityGroupPolicy(regionId, securityGroupId, ipProtocol, portRange, cidrIp string) bool {
	client, err := common.GetEcsClient(regionId)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【添加指定安全组ID的端口策略】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	// 入方向
	request := ecs.CreateAuthorizeSecurityGroupRequest()
	request.Scheme = "https"
	request.SourceCidrIp = cidrIp
	request.IpProtocol = ipProtocol
	request.PortRange = portRange
	request.SecurityGroupId = securityGroupId
	request.Priority = "1"

	response, err := client.AuthorizeSecurityGroup(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【添加指定安全组ID的端口策略】创建添加端口策略【入方向】请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	// 出方向
	requestEg := ecs.CreateAuthorizeSecurityGroupEgressRequest()
	requestEg.Scheme = "https"
	requestEg.DestCidrIp = cidrIp
	requestEg.IpProtocol = ipProtocol
	requestEg.PortRange = portRange
	requestEg.SecurityGroupId = securityGroupId

	responseEg, err := client.AuthorizeSecurityGroupEgress(requestEg)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【添加指定安全组ID的端口策略】创建添加端口策略【出方向】请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	return responseEg.IsSuccess() && response.IsSuccess()
}

// RemoveSecurityGroupPolicy 删除指定安全组ID的端口
func RemoveSecurityGroupPolicy(regionId, securityGroupId, ipProtocol, portRange, cidrIp string) bool {
	client, err := common.GetEcsClient(regionId)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【删除指定安全组ID的端口】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	// 入方向
	request := ecs.CreateRevokeSecurityGroupRequest()
	request.Scheme = "https"
	request.SourceCidrIp = cidrIp
	request.IpProtocol = ipProtocol
	request.PortRange = portRange
	request.SecurityGroupId = securityGroupId

	response, err := client.RevokeSecurityGroup(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【删除指定安全组ID的端口】删除添加端口策略【入方向】请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	// 出方向
	requestEg := ecs.CreateRevokeSecurityGroupEgressRequest()
	requestEg.Scheme = "https"
	requestEg.DestCidrIp = cidrIp
	requestEg.IpProtocol = ipProtocol
	requestEg.PortRange = portRange
	requestEg.SecurityGroupId = securityGroupId

	responseEg, err := client.RevokeSecurityGroupEgress(requestEg)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【删除指定安全组ID的端口】删除添加端口策略【出方向】请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	return responseEg.IsSuccess() && response.IsSuccess()
}

// ShowEcsSecurityGroupInfo 显示ecs安全组的信息
func ShowEcsSecurityGroupInfo(securityGroup ecs.DescribeSecurityGroupAttributeResponse) {
	fmt.Printf("安全组ID: %s \t 安全组名称: %s \t安全组描述: %s \t 策略条数:%d\n", securityGroup.SecurityGroupId,
		securityGroup.SecurityGroupName, securityGroup.Description, len(securityGroup.Permissions.Permission))
	var dates [][]string
	var accessPolicy string

	innerAccessPolicy := securityGroup.InnerAccessPolicy
	if innerAccessPolicy == "Accept" {
		accessPolicy = "内网互通"
	} else {
		accessPolicy = "内网隔离"
	}
	count := 0

	for _, permission := range securityGroup.Permissions.Permission {
		var direction string

		if permission.Direction == "ingress" {
			direction = "入口"
		} else if permission.Direction == "egress" {
			direction = "出口"
		} else {
			direction = "不区分方向"
		}

		data := []string{permission.Priority, direction, accessPolicy, permission.PortRange, permission.IpProtocol,
			permission.SourceCidrIp, permission.DestCidrIp, permission.CreateTime, permission.Description}
		dates = append(dates, data)
		count = count + 1
	}

	t := gotabulate.Create(dates)
	t.SetHeaders([]string{"优先级", "方向", "安全组策略", "端口信息", "协议类型", "源IP地址段", "目标IP地址段", "授权时间", "描述"})
	t.SetAlign("center")

	fmt.Println(t.Render("grid"))
}
