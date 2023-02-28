/*
Copyright © 2023 iiusky sky@03sec.com

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
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/iiiusky/alicloud-tools/common"
)

// CreateRamUser 创建RAM用户
func CreateRamUser(username string, password string) bool {
	client, err := common.GetRamClient()
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【创建RAM用户】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}
	createUserRequest := ram.CreateCreateUserRequest()
	createUserRequest.UserName = username
	createUserRequest.DisplayName = username
	createUserRequest.Scheme = "https"
	createUserResponse, err := client.CreateUser(createUserRequest)
	if err != nil || createUserResponse.IsSuccess() == false {
		common.Logger().Error(fmt.Sprintf("【创建RAM用户】创建用户请求发生异常,异常信息为 %s", err.Error()))
		fmt.Println(createUserResponse)
		return false
	}

	createLoginProfileRequest := ram.CreateCreateLoginProfileRequest()
	createLoginProfileRequest.Scheme = "https"
	createLoginProfileRequest.UserName = username
	createLoginProfileRequest.Password = password
	createLoginProfileResponse, err := client.CreateLoginProfile(createLoginProfileRequest)
	if err != nil || createLoginProfileResponse.IsSuccess() == false {
		common.Logger().Error(fmt.Sprintf("【创建RAM用户】创建用户登录文件请求发生异常,异常信息为 %s", err.Error()))
		fmt.Println(createLoginProfileResponse)
		return false
	}

	createAttachPolicyToUserRequest := ram.CreateAttachPolicyToUserRequest()

	createAttachPolicyToUserRequest.Scheme = "https"
	createAttachPolicyToUserRequest.PolicyType = "System"
	createAttachPolicyToUserRequest.PolicyName = "AdministratorAccess"
	createAttachPolicyToUserRequest.UserName = username
	createAttachPolicyToUserResponse, err := client.AttachPolicyToUser(createAttachPolicyToUserRequest)
	if err != nil || createAttachPolicyToUserResponse.IsSuccess() == false {
		common.Logger().Error(fmt.Sprintf("【创建RAM用户】创建授权用户超管权限请求发生异常,异常信息为 %s", err.Error()))
		fmt.Println(createAttachPolicyToUserRequest)
		return false
	}

	createGetAccountAliasRequest := ram.CreateGetAccountAliasRequest()
	createGetAccountAliasRequest.Scheme = "https"
	createGetAccountAliasResponse, err := client.GetAccountAlias(createGetAccountAliasRequest)
	if err != nil || createGetAccountAliasResponse.IsSuccess() == false {
		common.Logger().Error(fmt.Sprintf("【创建RAM用户】获取账号别名请求发生异常,异常信息为 %s", err.Error()))
		fmt.Println(createGetAccountAliasResponse)
		return false
	}
	fmt.Println("+++++++++++++++++登录信息+++++++++++++++++++")
	fmt.Printf("用户名: %s\n", username)
	fmt.Printf("登录名: %s@%s.onaliyun.com\n", username, createGetAccountAliasResponse.AccountAlias)
	fmt.Printf("密码: %s\n", password)
	fmt.Printf("登录地址1: %s\n", "https://signin.aliyun.com/login.htm")
	fmt.Printf("登录地址2: https://signin.aliyun.com/login.htm?username=%s@%s.onaliyun.com\n", username, createGetAccountAliasResponse.AccountAlias)
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	return true
}

// ListRamUser 列出所有RAM用户
func ListRamUser() {
	client, err := common.GetRamClient()
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【删除RAM用户】创建客户端发生异常,异常信息为 %s", err.Error()))
		return
	}

	request := ram.CreateListUsersRequest()
	request.Scheme = "https"
	response, err := client.ListUsers(request)
	for _, user := range response.Users.User {
		fmt.Println("++++++++++++++++++++++++++++++++++++++++")
		showUserInfo(user)
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	}
}

// DeleteRamUser 删除指定RAM用户
func DeleteRamUser(username string) bool {
	client, err := common.GetRamClient()
	if err != nil {
		common.Logger().Error(fmt.Sprintf("【删除RAM用户】创建客户端发生异常,异常信息为 %s", err.Error()))
		return false
	}

	createDetachPolicyFromUserRequest := ram.CreateDetachPolicyFromUserRequest()
	createDetachPolicyFromUserRequest.Scheme = "https"
	createDetachPolicyFromUserRequest.PolicyType = "System"
	createDetachPolicyFromUserRequest.PolicyName = "AdministratorAccess"
	createDetachPolicyFromUserRequest.UserName = username
	createDetachPolicyFromUserResponse, err := client.DetachPolicyFromUser(createDetachPolicyFromUserRequest)
	if err != nil || createDetachPolicyFromUserResponse.IsSuccess() == false {
		common.Logger().Error(fmt.Sprintf("【删除RAM用户】创建删除策略请求发生异常,异常信息为 %s", err.Error()))
		fmt.Println(createDetachPolicyFromUserResponse)
		return false
	}

	createDeleteUserRequest := ram.CreateDeleteUserRequest()
	createDeleteUserRequest.Scheme = "https"
	createDeleteUserRequest.UserName = username
	createDeleteUserResponse, err := client.DeleteUser(createDeleteUserRequest)
	if err != nil || createDeleteUserResponse.IsSuccess() == false {
		common.Logger().Error(fmt.Sprintf("【删除RAM用户】创建删除策略请求发生异常,异常信息为 %s", err.Error()))
		fmt.Println(createDeleteUserResponse.String())
		return false
	}

	return true
}

func showUserInfo(user ram.UserInListUsers) {
	fmt.Printf("用户ID: %s \n", user.UserId)
	fmt.Printf("用户名: %s \n", user.UserName)
	fmt.Printf("用户别名: %s \n", user.DisplayName)
	fmt.Printf("用户备注: %s \n", user.Comments)
}
