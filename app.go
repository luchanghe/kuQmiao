package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/luchanghe/jx3Robot/travel"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:generate cqcfg -c .
// cqp: 名称: jx3Robot
// cqp: 版本: 1.0.0:1
// cqp: 作者: luchanghe
// cqp: 简介: 剑网三QQ机器人
func main() {}

//onGroupMsg(1,2,60,60,"123","花价 梦江南 绣球花",4)
var jx3FlowerColorStr = map[string]string{
	"一级绣球花": "(白/红/紫)",
	"二级绣球花": "(粉/黄/蓝)",
	"一级郁金香": "(粉/红/黄)",
	"二级郁金香": "(白/混)",
	"一级牵牛花": "(红/绯/紫)",
	"二级牵牛花": "(黄/蓝)",
	"一级玫瑰":  "(粉/红/橙/黄/蓝)",
	"二级玫瑰":  "(白/紫/黑)",
	"三级玫瑰":  "(绿/混)",
	"一级百合":  "(白/粉/黄)",
	"二级百合":  "(橙/绿)",
	"一级荧光菌": "(白/红/黄)",
	"二级荧光菌": "(蓝/紫)",
}

func init() {
	cqp.AppID = "me.cqp.luchanghe.jx3Robot" // TODO: 修改为这个插件的ID
	cqp.GroupMsg = onGroupMsg
}

//设置剑网三盒子地址
const (
	jx3BoxFlower string = "https://next.jx3box.com/api/flower/price/rank"
)

type FlowerData struct {
	Max     int64    `json:"max"`
	Min     int64    `json:"min"`
	MaxLine []string `json:"maxLine"`
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	//剔除字符串首位空格
	msg = DeleteExtraSpace(msg)
	msg = strings.Trim(msg, " ")
	//拆解字符串为数组
	strArr := strings.Split(msg, " ")
	//获取命令长度
	strArrLen := len(strArr)
	if strArrLen > 3 || strArrLen == 1 {
		return 0
	}
	switch strArr[0] {
	case "花价":
		//设置请求URL
		var serverName, flowerName string
		//默认梦江南
		if strArrLen == 2 {
			serverName = "梦江南"
			flowerName = strArr[1]
		}
		//命令为3表示附带了服务器
		if strArrLen == 3 {
			serverName = strArr[1]
			flowerName = strArr[2]
		}
		//适配
		serverName = ServerFitter(serverName)
		flowerName = FlowerFitter(flowerName)

		//url过滤
		queryUrl := fmt.Sprintf(jx3BoxFlower+"?server=%s&flower=%s", serverName, flowerName)
		//请求数据
		jsonStr := HttpGet(queryUrl)
		var flowerMap map[string]FlowerData
		json.Unmarshal([]byte(jsonStr), &flowerMap)
		var msg strings.Builder
		if len(flowerMap) == 0 {
			cqp.SendGroupMsg(fromGroup, "没查询到嗷！")
			return 0
		}
		for key, value := range flowerMap {
			msg.WriteString(key)
			msg.WriteString(jx3FlowerColorStr[key])
			msg.WriteString("\n")
			msg.WriteString("价格 ")
			msg.WriteString(strconv.FormatInt(int64(value.Max), 10))
			msg.WriteString("\n")
			var address strings.Builder
			for k1, v1 := range value.MaxLine {
				if k1 == 3 {
					address.WriteString("线\n")
					break
				} else {
					address.WriteString(strings.TrimRight(v1, " 线"))
				}
				if k1 != 2 {
					address.WriteString("、")
				}
			}
			msg.WriteString(address.String())
		}
		cqp.SendGroupMsg(fromGroup, strings.TrimRight(msg.String(), "\n"))
	case "宠物游历":
		result := travel.Result(strArr[1])
		cqp.SendGroupMsg(fromGroup, result)
	}
	return 0
}

//处理GET请求
func HttpGet(url string) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return result.String()
}

/*
函数名：DeleteExtraSpace(s string) string
功  能:删除字符串中多余的空格(含tab)，有多个空格时，仅保留一个空格，同时将字符串中的tab换为空格
参  数:s string:原始字符串
返回值:string:删除多余空格后的字符串
创建时间:2018年12月3日
修订信息:
*/
func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)       //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

func ServerFitter(serverName string) string {
	switch serverName {
	case "双梦":
		serverName = "梦江南"
	case "华乾":
		serverName = "乾坤一掷"
	}
	return serverName
}

func FlowerFitter(flowerName string) string {
	switch flowerName {
	case "蘑菇":
		flowerName = "荧光菌"
	}
	return flowerName
}
