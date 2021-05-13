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
	"errors"
	"fmt"
	"github.com/iiiusky/alicloud-tools/common"
	"github.com/spf13/cobra"
	"os"
)

var accessKey string
var secretKey string
var regionId string
var stsAccessKey string
var stsSecretKey string
var stsToken string
var verbose bool
var useSTS bool
var showRegions bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "AliCloud-Tools",
	Short: "阿里云API利用工具",
	Long:  `该工具主要是方便快速使用阿里云api执行一些操作`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use == "version" {
			return nil
		}
		common.Verbose = verbose

		if (accessKey == "" || secretKey == "") && useSTS == false {
			return errors.New("请设置ak以及sk的值")
		}

		if useSTS && (stsAccessKey == "" || stsSecretKey == "" || stsToken == "") {
			return errors.New("请设置stsAccessKey、stsSecretKey、stsToken的值")
		}

		if useSTS {
			common.STSAccessKey = stsAccessKey
			common.STSSecretKey = stsSecretKey
			common.STSToken = stsToken
			common.UseSTS = useSTS
		} else {
			common.AccessKey = accessKey
			common.SecretKey = secretKey
		}

		if !common.InitEcsRegions() {
			return errors.New("ak、sk验证失败.")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if showRegions {
			common.ShowRegions()
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&showRegions, "regions", false, "显示所有地域信息")
	rootCmd.PersistentFlags().StringVarP(&accessKey, "ak", "a", "", "阿里云 AccessKey")
	rootCmd.PersistentFlags().StringVarP(&secretKey, "sk", "s", "", "阿里云 SecretKey")

	rootCmd.PersistentFlags().StringVar(&stsAccessKey, "sak", "", "阿里云 STS AccessKey")
	rootCmd.PersistentFlags().StringVar(&stsSecretKey, "ssk", "", "阿里云 STS SecretKey")
	rootCmd.PersistentFlags().StringVar(&stsToken, "token", "", "阿里云 STS Session Token")
	rootCmd.PersistentFlags().BoolVar(&useSTS, "sts", false, "启用STSToken模式")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细的执行过程")
	rootCmd.PersistentFlags().StringVarP(&regionId, "rid", "r", "", "阿里云 地域ID,在其他支持rid的子命令中,如果设置了地域ID,则只显示指定区域的信息,否则为全部.")
}
