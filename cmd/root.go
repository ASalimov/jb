/*
Copyright © 2020 asalimov

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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/ASalimov/bar"
	"github.com/ttacon/chalk"
	"strings"

	//"github.com/superhawk610/bar"
	//"github.com/ttacon/chalk"
	"net/url"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/chzyer/readline"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jb",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jb.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".jb" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".jb")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}



func run_job(jh *JHelper, name string, args []string){
	data := map[string]string{}
	fmt.Println("args: ",args, "name ", name)
	time.Sleep(time.Second*3)
	ji:=jh.getJobInfo(name)
	if len(ji.getParameterDefinitions())>len(args){
		fmt.Println("Please specify this build requires parameters:")
	}
	for i, pd:=range ji.getParameterDefinitions(){
		if len(args)>i{
			data[pd.Name]=args[i]
		}else{
			rl, err := readline.New(pd.Name+": ")
			if err != nil {
				panic(err)
			}
			defer rl.Close()
			line, err := rl.Readline()
			if err != nil { // io.EOF
				break
			}
			data[pd.Name]=line
		}

	}

	dataArr := []map[string]string{}
	for key, val := range data {
		dataArr = append(dataArr, map[string]string{"name": key, "value": val})
	}
	queryJson := map[string][]map[string]string{
		"parameter": dataArr,
	}
	queryBytes, _ := json.Marshal(queryJson)
	urlquery := url.Values{}
	urlquery.Add("json", string(queryBytes))
	lastBuild, err:= jh.getLastSuccessfulBuildInfo(name)
	if err != nil {
		panic("Failed to build "+ err.Error())
	}

	jh.build(name, []byte(urlquery.Encode()))

	fmt.Println(lastBuild.Duration)
	b := bar.NewWithOpts(
		bar.WithDimensions(100, 100),
		bar.WithLines(5),
		bar.WithFormat(
			fmt.Sprintf(
				" %sbuilding...%s :percent :bar %s:eta %s ",
				chalk.Blue,
				chalk.Reset,
				chalk.Green,
				chalk.Reset)))
	b.Tick()
	ticks:=1
	t := 0
	time.Sleep(3*time.Second)
	cursor:="0"
	stime:=getTime()
	for t > -1 {
		if t%5 == 0 && t > 1 {
			curBuild, err := jh.getBuildInfo(name,ji.NextBuildNumber)
			if err != nil {
				if getTime()-stime>int64(30*time.Millisecond) {
					b.Interrupt("Failed to build1, path = " + name)
					b.Done()
					return
				}
			} else {
				if !curBuild.Building {
					if curBuild.Result == "SUCCESS" {
						b.Interrupt("Finish")
						//b.Tick()
						for ticks<100{
							b.Tick()
							ticks++
						}
						b.Done()
						return
					} else {
						////path := "/job/" + name + "/" + curBuild.Id + "/console"
						b.Done()
						fmt.Println("Failure")
						return
					}
				}
			}
		}
		output, nextCursor,err:=jh.console(name, ji.NextBuildNumber,cursor)
		if err==nil{
			cursor=nextCursor
			//r:= strings.NewReader(output)
			//scanner:= bufio.NewReader(r)
			j:=1
			for  {
				lines:= strings.Split(output,"\n")
				count:=len(lines)
				if count>50{
					count=50;
				}
				for i:=count; i>=1;i--{
					rline:=[]rune(string(lines[len(lines)-i]))
					if err != nil { // io.EOF
						break


					}
					if len(rline)>100{
						rline=rline[:100]
					}
					b.Interruptf(string(rline))
					time.Sleep(20*time.Millisecond)
				}
				break
				j++
				//time.Sleep(4*time.Millisecond)
			}
			ctime:=getTime()
			dtime:=ctime-stime
			newTicks:=int(float64(dtime)/float64(lastBuild.Duration)*100)
			//b.InterruptfInOneLine("ctime %d, stime %d, new tick %d",ctime, stime, newTicks)
			for ticks<newTicks&&ticks<99{
				b.Tick()
				ticks++
			}
			if dtime<500{
				time.Sleep(time.Duration(500-dtime)*time.Millisecond)
			}
		}else{
			time.Sleep(100*time.Millisecond)
		}
		t++
	}
	b.Interrupt("Timeout")
	b.Done()
	return
}



func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
