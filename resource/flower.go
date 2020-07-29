package resource

import (
	"encoding/json"
	"fmt"
	"github.com/luchanghe/jx3Robot/tool"
	"strconv"
	"strings"
)

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

type FlowerData struct {
	Max     int64    `json:"max"`
	Min     int64    `json:"min"`
	MaxLine []string `json:"maxLine"`
}

func SearchFlower(serverName string, flowerName string) string {
	//url过滤
	queryUrl := fmt.Sprintf(jx3BoxFlower+"?server=%s&flower=%s", serverName, flowerName)
	//请求数据
	jsonStr := tool.HttpGet(queryUrl)
	var flowerMap map[string]FlowerData
	json.Unmarshal([]byte(jsonStr), &flowerMap)
	var msg strings.Builder
	if len(flowerMap) == 0 {
		return "没查询到嗷"
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
	return strings.TrimRight(msg.String(), "\n")
}
