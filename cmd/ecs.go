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
package cmd

import (
	"fmt"
	"github.com/iiiusky/alicloud-tools/core"
	"github.com/spf13/cobra"
)

var listEcs bool
var isRunner bool
var instancesId string

var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "ECS 操作(查询/执行命令),当前命令支持地域ID设置.",
	Long:  `该命令主要可以查看所有ecs列表、查询单个ecs信息、执行命令,当前命令支持地域ID设置.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if instancesId != "" {
			fmt.Printf("查找地域名称为 %s 且 InstanceId为 %s 的实例\n", regionId, instancesId)
			fmt.Println("++++++++++++++++++++++++++++++++++++++++")
			core.ShowInstanceInfo(core.QuerySingleInstance(regionId, instancesId))
			fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
		}
		if listEcs {
			core.ShowInstancesInfo(core.GetAllInstances(regionId, true), isRunner)
		}
		return nil
	},
}

func init() {
	ecsCmd.Flags().BoolVar(&listEcs, "list", false, "显示ecs列表")
	ecsCmd.Flags().BoolVar(&isRunner, "runner", true, "是否只显示正在运行的实例信息,默认为true")
	ecsCmd.Flags().StringVar(&instancesId, "eid", "", "实例ID")

	rootCmd.AddCommand(ecsCmd)
}
