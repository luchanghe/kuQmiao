package resource

//匹配服务器
func ServerFitter(serverName string) string {
	switch serverName {
	case "双梦":
		serverName = "梦江南"
	case "华乾":
		serverName = "乾坤一掷"
	}
	return serverName
}

//匹配花名
func FlowerFitter(flowerName string) string {
	switch flowerName {
	case "蘑菇":
		flowerName = "荧光菌"
	}
	return flowerName
}
