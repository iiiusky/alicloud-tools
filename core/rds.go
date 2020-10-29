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
	"alicloud-tools/common"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"strconv"
	"strings"
)

// 获取指定区域的所有RDS实例列表
func GetRegionDBInstances(regionId string) (instances []rds.DBInstance) {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取指定区域的所有RDS实例列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return
	}

	for i := 1; i < 100; i++ {
		request := rds.CreateDescribeDBInstancesRequest()
		request.Scheme = "https"
		request.PageSize = requests.Integer("100")
		request.InstanceLevel = requests.NewInteger(1)
		request.PageNumber = requests.Integer(strconv.Itoa(i))

		response, err := client.DescribeDBInstances(request)
		if err != nil {
			common.Logger().Error(fmt.Sprintf("【获取指定区域的所有RDS实例列表】创建获取实例请求发生异常,异常信息为 %s", err.Error()))
			return
		}

		instances = append(instances, response.Items.DBInstance...)

		if len(response.Items.DBInstance) == 0 {
			break
		}
	}

	return instances
}

// 获取所有RDS实例
func GetAllDBInstances(regionId string, printInfo bool) (instances []rds.DBInstance) {
	for _, region := range common.Regions {
		if regionId != "" && regionId != region.RegionId {
			continue
		}

		if printInfo {
			fmt.Printf("正在扫描  %s 区域的 RDS\t", region.LocalName)
		}

		regionInstances := GetRegionDBInstances(region.RegionId)
		instances = append(instances, regionInstances...)

		if printInfo {
			fmt.Printf("扫描到 %d 台 RDS\n", len(regionInstances))
		}
	}

	return instances
}

// 获取指定rds实例的账号列表
func getRdsAccounts(regionId string, dbInstanceId string) (accounts []rds.DBInstanceAccount) {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取指定rds实例的账号列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return accounts
	}

	for i := 1; i < 100; i++ {
		request := rds.CreateDescribeAccountsRequest()
		request.Scheme = "https"
		request.DBInstanceId = dbInstanceId
		request.PageNumber = requests.NewInteger(i)

		response, err := client.DescribeAccounts(request)
		if err != nil {
			common.Logger().Error(fmt.Sprintf("【获取指定rds实例的账号列表】创建获取账号列表请求发生异常,异常信息为 %s", err.Error()))
			return accounts
		}

		accounts = append(accounts, response.Accounts.DBInstanceAccount...)

		if len(response.Accounts.DBInstanceAccount) == 0 {
			break
		}
	}

	return accounts
}

// 增加指定rds实例的账号列表
func AddRdsAccount(regionId, dbInstanceId, accountName, accountPassword, accountDescription string, isSuper bool) bool {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【增加指定rds实例的账号列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	request := rds.CreateCreateAccountRequest()
	request.Scheme = "https"
	request.DBInstanceId = dbInstanceId
	request.AccountName = accountName
	request.AccountPassword = accountPassword
	request.AccountDescription = accountDescription

	if isSuper {
		request.AccountType = "Super"
	} else {
		request.AccountType = "Normal"
	}

	response, err := client.CreateAccount(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【增加指定rds实例的账号列表】发送创建账号请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	return response.IsSuccess()

}

// 删除指定rds实例的账号列表
func DelRdsAccount(regionId, dbInstanceId, accountName string) bool {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【删除指定rds实例的账号列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	request := rds.CreateDeleteAccountRequest()
	request.Scheme = "https"
	request.DBInstanceId = dbInstanceId
	request.AccountName = accountName

	response, err := client.DeleteAccount(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【删除指定rds实例的账号列表】发送删除账号请求发生异常,异常信息为 %s", err.Error()))
		return false
	}

	return response.IsSuccess()
}

// 获取指定rds实例的备份列表
func GetRdsBackups(regionId string, dbInstanceId string) (backups []rds.Backup) {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取指定rds实例的备份列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return backups
	}

	for i := 1; i < 100; i++ {
		request := rds.CreateDescribeBackupsRequest()
		request.Scheme = "https"
		request.DBInstanceId = dbInstanceId
		request.PageSize = requests.NewInteger(100)
		request.PageNumber = requests.NewInteger(i)

		response, err := client.DescribeBackups(request)
		if err != nil {
			common.Logger().Error(fmt.Sprintf("【获取指定rds实例的备份列表】发送获取请求发生异常,异常信息为 %s", err.Error()))
			return backups
		}

		backups = append(backups, response.Items.Backup...)

		if len(response.Items.Backup) == 0 {
			break
		}
	}

	return backups
}

// 创建指定rds实例的备份
func CreateRdsBackups(regionId string, dbInstanceId string) (backups []rds.Backup) {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取指定rds实例的备份列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return backups
	}

	for i := 1; i < 100; i++ {
		request := rds.CreateDescribeBackupsRequest()
		request.Scheme = "https"
		request.DBInstanceId = dbInstanceId
		request.PageSize = requests.NewInteger(100)
		request.PageNumber = requests.NewInteger(i)

		response, err := client.DescribeBackups(request)
		if err != nil {
			common.Logger().Error(fmt.Sprintf("【获取指定rds实例的备份列表】发送获取请求发生异常,异常信息为 %s", err.Error()))
			return backups
		}

		backups = append(backups, response.Items.Backup...)

		if len(response.Items.Backup) == 0 {
			break
		}
	}

	return backups
}

// 删除指定rds实例的备份
func DelRdsBackups(regionId string, dbInstanceId string) (backups []rds.Backup) {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【获取指定rds实例的备份列表】创建客户端发生异常,异常信息为 %s", err.Error()))
		return backups
	}

	for i := 1; i < 100; i++ {
		request := rds.CreateDescribeBackupsRequest()
		request.Scheme = "https"
		request.DBInstanceId = dbInstanceId
		request.PageSize = requests.NewInteger(100)
		request.PageNumber = requests.NewInteger(i)

		response, err := client.DescribeBackups(request)
		if err != nil {
			common.Logger().Error(fmt.Sprintf("【获取指定rds实例的备份列表】发送获取请求发生异常,异常信息为 %s", err.Error()))
			return backups
		}

		backups = append(backups, response.Items.Backup...)

		if len(response.Items.Backup) == 0 {
			break
		}
	}

	return backups
}

// 查询单个DB实例详细信息
func queryDbInstanceAttribute(regionId, instanceId string) (dbInfo rds.DBInstanceAttribute) {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【查询单个DB实例详细信息】创建客户端发生异常,异常信息为 %s", err.Error()))
		return dbInfo
	}

	request := rds.CreateDescribeDBInstanceAttributeRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceId
	request.Expired = "False"

	response, err := client.DescribeDBInstanceAttribute(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【查询单个DB实例详细信息】创建获取实例请求发生异常,异常信息为 %s", err.Error()))
		return dbInfo
	}

	return response.Items.DBInstanceAttribute[0]
}

// 查询单个DB实例网络详细信息
func queryDbInstanceNetInfo(regionId, instanceId string) (dbInfo rds.DBInstanceNetInfo) {
	client, err := rds.NewClientWithAccessKey(regionId, common.AccessKey, common.SecretKey)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【查询单个DB实例网络详细信息】创建客户端发生异常,异常信息为 %s", err.Error()))
		return dbInfo
	}

	request := rds.CreateDescribeDBInstanceNetInfoRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceId

	response, err := client.DescribeDBInstanceNetInfo(request)
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【查询单个DB实例网络详细信息】创建获取实例请求发生异常,异常信息为 %s", err.Error()))
		return dbInfo
	}

	return response.DBInstanceNetInfos.DBInstanceNetInfo[0]
}

// 显示rds列表信息
func ShowDBInstancesInfo(instances []rds.DBInstance, isRunner bool) {
	for _, instance := range instances {
		if isRunner && instance.DBInstanceStatus != "Running" {
			continue
		}

		fmt.Println("++++++++++++++++++++++++++++++++++++++++")
		ShowDBInstanceInfo(instance)
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
}

// 显示rds信息
func ShowDBInstanceInfo(instance rds.DBInstance) {
	accounts := getRdsAccounts(instance.RegionId, instance.DBInstanceId)
	dbInfo := queryDbInstanceAttribute(instance.RegionId, instance.DBInstanceId)
	netInfo := queryDbInstanceNetInfo(instance.RegionId, instance.DBInstanceId)

	var DBInstanceNetType string

	if strings.ToLower(netInfo.IPType) == "inner" || strings.ToLower(netInfo.IPType) == "private" {
		DBInstanceNetType = "内网"
	} else if strings.ToLower(netInfo.IPType) == "public" {
		DBInstanceNetType = "外网"
	}

	fmt.Printf("当前 %s 实例下共有 %d 个账号,网络连接类型为 %s", instance.DBInstanceId, len(accounts), DBInstanceNetType)
	fmt.Printf("实例描述: %s \n", instance.DBInstanceDescription)
	fmt.Printf("实例规格: %s \n", instance.DBInstanceClass)
	fmt.Printf("实例状态: %s \n", instance.DBInstanceStatus)
	fmt.Printf("VPC ID: %v \n", instance.VpcId)
	fmt.Printf("数据库类型: %s \n", instance.Engine)
	fmt.Printf("数据库版本: %s \n", instance.EngineVersion)
	fmt.Printf("连接地址: %v \n", netInfo.ConnectionString)
	fmt.Printf("IP地址: %v \n", netInfo.IPAddress)
	fmt.Printf("连接端口: %s \n", netInfo.Port)
	fmt.Printf("实例创建时间: %s \n", instance.CreateTime)
	fmt.Printf("实例过期时间: %s \n", instance.ExpireTime)
	fmt.Printf("实例销毁时间: %s \n", instance.DestroyTime)
	fmt.Printf("实例内存: %v \n", dbInfo.DBInstanceMemory)
	fmt.Printf("实例CPU: %s \n", dbInfo.DBInstanceCPU)
	fmt.Printf("实例存储空间: %v M\n", dbInfo.DBInstanceStorage)
	fmt.Printf("实例地域ID: %s \n", instance.RegionId)

	if strings.ToLower(instance.PayType) == "postpaid" {
		fmt.Println("付费类型: 按量付费")
	} else {
		fmt.Println("付费类型: 包年包月")
	}

	switch strings.ToLower(instance.DBInstanceType) {
	case "primary":
		fmt.Println("实例类型: 主实例")
		break
	case "readonly":
		fmt.Println("实例类型: 只读实例")
		break
	case "guard":
		fmt.Println("实例类型: 灾备实例")
		break
	case "temp":
		fmt.Println("实例类型: 临时实例")
		break
	}

	if strings.ToLower(instance.InstanceNetworkType) == "classic" {
		fmt.Println("网络类型: 经典网络")
	} else {
		fmt.Println("网络类型: VPC网络")
	}

	if strings.ToLower(instance.ConnectionMode) == "standard" {
		fmt.Println("访问模式: 标准访问模式")
	} else if strings.ToLower(instance.ConnectionMode) == "safe" {
		fmt.Println("访问模式: 数据库代理模式")
	} else {
		fmt.Printf("访问模式: %s \n", instance.ConnectionMode)
	}

	if strings.ToLower(instance.Category) == "basic" {
		fmt.Println("实例系列: 基础版")
	} else if strings.ToLower(instance.Category) == "highavailability" {
		fmt.Println("实例系列: 高可用版")
	} else if strings.ToLower(instance.Category) == "finance" {
		fmt.Println("实例系列: 三节点企业版")
	}

	for _, securityIPGroup := range netInfo.SecurityIPGroups.SecurityIPGroup {
		fmt.Printf("白名单信息[%s]: %s \n", securityIPGroup.SecurityIPGroupName, securityIPGroup.SecurityIPs)
	}

	for i, account := range accounts {
		fmt.Printf("实例第%d个账号信息-账户名: %s \n", i, account.AccountName)
		fmt.Printf("实例第%d个账号信息-账号状态: %s \n", i, account.AccountStatus)
		fmt.Printf("实例第%d个账号信息-账号描述: %s \n", i, account.AccountDescription)
		switch strings.ToLower(account.AccountType) {
		case "normal":
			fmt.Printf("实例第%d个账号信息-账号类型: %s \n", i, "普通账号")
			break
		case "super":
			fmt.Printf("实例第%d个账号信息-账号类型: %s \n", i, "高权限账号")
			break
		case "sysadmin":
			fmt.Printf("实例第%d个账号信息-账号类型: %s \n", i, "具备超级权限（SA）的账号")
			break
		}
	}

}
