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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch ramAction {
		case "add":
			if core.CreateRamUser(username, password) {
				fmt.Printf("创建用户 %s 成功\n", username)
			} else {
				fmt.Printf("创建用户 %s 失败\n", username)
			}
			break
		case "del":
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
