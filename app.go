package main

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/cqp"
	"github.com/luchanghe/jx3Robot/resource"
	"github.com/luchanghe/jx3Robot/tool"
	"strings"
)

//go:generate cqcfg -c .
// cqp: 名称: jx3Robot
// cqp: 版本: 1.0.0:1
// cqp: 作者: luchanghe
// cqp: 简介: 剑网三QQ机器人
func main() {}

//onGroupMsg(1,2,60,60,"123","花价 梦江南 绣球花",4)

func init() {
	cqp.AppID = "me.cqp.luchanghe.jx3Robot" // TODO: 修改为这个插件的ID
	cqp.GroupMsg = onGroupMsg
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	//剔除字符串首位空格
	msg = tool.DeleteExtraSpace(msg)
	msg = strings.Trim(msg, " ")
	//拆解字符串为数组
	strArr := strings.Split(msg, " ")
	//获取命令长度
	strArrLen := len(strArr)
	if strArrLen > 3 || strArrLen == 1 {
		return 0
	}
	result := ""
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
		serverName = resource.ServerFitter(serverName)
		flowerName = resource.FlowerFitter(flowerName)
		result = resource.SearchFlower(serverName, flowerName)

	case "宠物游历":
		address := strArr[1]
		result = resource.SearchPetTravel(address)
	case "阿喵":
		switch strArr[1] {
		case "口吐芬芳":
			result = resource.SearchFuck()
		}
	}
	if result != "" {
		sendGroupMessage(fromGroup, result)
	}
	return 0
}

var developmentMode bool = false

func sendGroupMessage(fromGroup int64, msg string) {
	if developmentMode == true {
		fmt.Println(msg)
	} else {
		cqp.SendGroupMsg(fromGroup, msg)
	}
}

func test() {
	fmt.Println("花价")
	onGroupMsg(1, 2, 60, 60, "123", "花价 梦江南 蘑菇", 4)
	fmt.Println("宠物游历")
	onGroupMsg(1, 2, 60, 60, "123", "宠物游历 扬州", 4)
	fmt.Println("口吐芬芳")
	onGroupMsg(1, 2, 60, 60, "123", "阿喵 口吐芬芳", 4)
}
