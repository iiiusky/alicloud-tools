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
	"strings"
)

var securityGroupId string
var ipProtocol string
var portRange string
var cidrIp string
var action string

var sgCmd = &cobra.Command{
	Use:   "sg",
	Short: "安全组操作,当前命令支持地域ID设置.",
	Long:  `该命令主要针对安全组的策略进行操作,比如增加组策略中开放的端口、删除组策略中指定的端口开放信息,当前命令支持地域ID设置.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ipProtocol = strings.ToLower(ipProtocol)

		if regionId == "" {
			regionId = "cn-beijing"
		}

		if securityGroupId == "" {
			return errors.New("安全组ID不允许为空")
		}

		switch strings.ToLower(action) {
		case "add":
			if core.AddSecurityGroupPolicy(regionId, securityGroupId, ipProtocol, portRange, cidrIp) {
				fmt.Printf("向安全组 %s 添加新的端口策略成功", securityGroupId)
			} else {
				fmt.Printf("向安全组 %s 添加新的端口策略失败", securityGroupId)
			}
			break
		case "del":
			if core.RemoveSecurityGroupPolicy(regionId, securityGroupId, ipProtocol, portRange, cidrIp) {
				fmt.Printf("安全组 %s 删除端口范围为 %s 原地址为 %s 协议为 %s 的端口策略成功", securityGroupId, portRange, cidrIp, ipProtocol)
			} else {
				fmt.Printf("安全组 %s 删除端口范围为 %s 原地址为 %s 协议为 %s 的端口策略失败", securityGroupId, portRange, cidrIp, ipProtocol)
			}
			break
		default:
			info := core.GetEcsSecurityGroupInfo(regionId, securityGroupId)
			core.ShowEcsSecurityGroupInfo(info)
			break
		}

		return nil
	},
}

func init() {
	sgCmd.Flags().StringVar(&securityGroupId, "sid", "", "* 安全组ID")
	sgCmd.Flags().StringVar(&action, "action", "", "操作类型,add-增加,del-删除")
	sgCmd.Flags().StringVar(&ipProtocol, "protocol", "tcp", "传输层协议,默认为tcp协议,当前支持tcp/udp/icmp/gre/all")
	sgCmd.Flags().StringVar(&portRange, "port", "", "端口范围,TCP/UDP协议:取值范围为1~65535,使用斜线（/）隔开起始端口和终止端口,例如：1/200\nICMP协议：-1/-1\nGRE协议：-1/-1\nIpProtocol取值为all：-1/-1")
	sgCmd.Flags().StringVar(&cidrIp, "ip", "", "源端IPv4 CIDR地址段,支持CIDR格式和IPv4格式的IP地址范围。")

	rootCmd.AddCommand(sgCmd)
}
