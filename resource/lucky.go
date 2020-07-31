package resource

import (
	"encoding/json"
	"fmt"
	"github.com/luchanghe/jx3Robot/tool"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

var serendipity string = "三尺青锋,三山四海,乱世舞姬,塞外宝驹,天涯无归,少年行,惜往日,扶摇九天,护佑苍生,故园风雨,济苍生,清风捕王,炼狱厨神,生死判,茶馆奇缘,虎啸山林,阴阳两界,雪山恩仇,韶华故,黑白路"

type luckyStruct struct {
	DwTime      float64 `json:"dwTime"`
	Data_str    string  `json:"data_str"`
	Server      string  `json:"server"`
	Serendipity string  `json:"serendipity"`
	Name        string  `json:"name"`
}

func LuckyLoop() string {

	url := fmt.Sprintf("%s?serendipity=%s&pageIndex=1&pageSize=1", jx3BoxLucky, serendipity)
	jsonStr := tool.HttpGet(url)
	//jsonStr := `{"code":0,"msg":"","data":{"page":{"index":1,"pageSize":3,"total":6634,"pageTotal":2212},"data":[{"lang":"zhcn","gegion":"电信五区","server":"剑胆琴心","serendipity":"少年行","name":"醉风一葬","status":0,"method":1,"count":1,"dwTime":1596175267,"data_str":"2020-07-31 14:01:07"},{"lang":"zhcn","gegion":"电信八区","server":"绝代天骄","serendipity":"天涯无归","name":"司墨琰","status":0,"method":1,"count":1,"dwTime":1596175015,"data_str":"2020-07-31 13:56:55"},{"lang":"zhcn","gegion":"电信一区","server":"长安城","serendipity":"故园风雨","name":"倾之暖","status":0,"method":2,"count":0,"dwTime":1596175011,"data_str":"2020-07-31 13:56:51"}],"historyData":null}}`
	var tempMap map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &tempMap)
	data := tempMap["data"].(map[string]interface{})["data"]
	var zeroKeyTimeStamp string //设置一个用来保存数据最新时间戳的变量

	_, err := os.Stat("luckTime.txt")
	if err != nil {
		if os.IsNotExist(err) {
			newFile, _ := os.Create("luckTime.txt")
			defer newFile.Close()
		}
	}
	//打开奇遇的缓存文件
	readFile, err := ioutil.ReadFile("luckTime.txt")
	if err != nil {
		log.Fatal(err)
	}
	lastTime := string(readFile[:])
	log.Printf("读文件", lastTime)

	for key, lucky := range data.([]interface{}) {
		var luckyStruct = luckyStruct{}
		mapstructure.Decode(lucky, &luckyStruct)
		timeStamp := luckyStruct.DwTime
		lastTimeStamp := big.NewRat(1, 1).SetFloat64(timeStamp).FloatString(0)
		if key == 0 {
			zeroKeyTimeStamp = lastTimeStamp //当第一次循环时候将时间戳赋值给最新时间
		}
		if lastTime == "" || lastTime == lastTimeStamp {
			//如果第一条的时间重复了说明暂时没有更新
			ioutil.WriteFile("luckTime.txt", []byte(zeroKeyTimeStamp), 0666)
			break
		}
		//否则就进行推送
		for qqGroup, serverGroup := range GroupMap {
			if serverGroup == luckyStruct.Server {
				tool.SendGroupMessage(qqGroup, fmt.Sprintf(
					"服务器：%s\n奇遇：%s\n触发玩家：%s\n时间：%s", luckyStruct.Server, luckyStruct.Serendipity, luckyStruct.Name, luckyStruct.Data_str))
			}
		}
	}
	return ""
}
