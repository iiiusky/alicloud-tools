/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"github.com/iiiusky/alicloud-tools/core"
	"github.com/spf13/cobra"
)

var username string
var password string
var ramAction string

// ramCmd represents the ram command
var ramCmd = &cobra.Command{
	Use:   "ram",
	Short: "RAM账号增删查操作",
	Long:  `支持新增、删除、查看RAM账号功能`,
	Run: func(cmd *cobra.Command, args []string) {
		switch ramAction {
		case "add":
			if username == "" || password == "" {
				fmt.Printf("请使用--username 以及--password 指定用户名及密码\n")
				break
			}
			if core.CreateRamUser(username, password) {
				fmt.Printf("创建用户 %s 成功\n", username)
			} else {
				fmt.Printf("创建用户 %s 失败\n", username)
			}
			break
		case "del":
			if username == "" {
				fmt.Printf("请使用--username 指定用户名\n")
				break
			}
			if core.DeleteRamUser(username) {
				fmt.Printf("删除用户 %s 成功\n", username)
			} else {
				fmt.Printf("删除用户 %s 失败\n", username)
			}
			break
		case "list":
			core.ListRamUser()
			break
		default:
			fmt.Println("未知操作")
		}
	},
}

func init() {
	ramCmd.Flags().StringVar(&username, "username", "", "* 安全组ID")
	ramCmd.Flags().StringVar(&password, "password", "", "RAM密码")
	ramCmd.Flags().StringVar(&ramAction, "action", "", "操作类型,add-增加,del-删除,list-列出所有用户")

	rootCmd.AddCommand(ramCmd)
}
