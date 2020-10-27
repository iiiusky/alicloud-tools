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
	"alicloud-tools/core"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

var windowsDefaultScriptType string
var commandContent string
var instanceIds []string

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "执行命令,当前命令支持地域ID设置.",
	Long:  `执行指定的命令,当前命令支持地域ID设置.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if commandContent == "" || len(instanceIds) == 0 {
			return errors.New("执行的命令以及实例ID不允许为空")
		}
		fmt.Printf("开始获取指定实例 %v 的基本信息,耗时可能较长,请耐心等待...\n", instanceIds)
		instances := core.GetAllInstances(regionId, false)

		for _, instance := range instances {
			for _, instanceId := range instanceIds {
				if instance.InstanceId == instanceId {
					cloudAssistantStatus := core.CheckCloudAssistantStatus(instance.RegionId, instanceId)
					fmt.Printf("检测当前实例 %s 是否安装云助手. 检测结果: %v \n", instanceId, cloudAssistantStatus)
					if !cloudAssistantStatus {
						fmt.Printf("当前实例 %s 未安装云助手,无法执行命令. \n", instanceId)
						continue
					}

					fmt.Printf("当前实例 %s 的操作系统类型为 %s,", instance.InstanceId, instance.OsType)
					switch instance.OsType {
					case "windows":
						if windowsDefaultScriptType == "RunBatScript" || windowsDefaultScriptType == "RunPowerShellScript" {
							fmt.Printf("正在以 %s 的方式执行命令,命令内容: %s ", windowsDefaultScriptType, commandContent)
							status := core.EcsRunCommand(instance.RegionId, windowsDefaultScriptType, commandContent, instanceId)
							fmt.Printf("执行命令结果状态为: %v \n", status)
						} else {
							fmt.Println("执行命令失败,未知的执行命令类型.")
						}
						break
					default:
						fmt.Printf("正在以 %s 的方式执行命令,命令内容: %s ", "RunShellScript", commandContent)
						status := core.EcsRunCommand(instance.RegionId, "RunShellScript", commandContent, instanceId)
						fmt.Printf("执行命令结果状态为: %v \n", status)
						break
					}
				}
			}
		}

		return nil
	},
}

func init() {
	execCmd.Flags().StringVarP(&windowsDefaultScriptType, "scriptType", "t", "RunBatScript", "Windows类型的实例脚本执行类型,Linux类型实例全部为Shell,默认为RunBatScript,类型说明:[RunBatScript：适用于Windows实例的Bat脚本。\nRunPowerShellScript：适用于Windows实例的PowerShell脚本。]")
	execCmd.Flags().StringVarP(&commandContent, "command", "c", "", "执行的命令")
	execCmd.Flags().StringSliceVarP(&instanceIds, "instanceIds", "I", nil, "要执行命令的实例ID,多个使用英文逗号分隔")

	ecsCmd.AddCommand(execCmd)
}
