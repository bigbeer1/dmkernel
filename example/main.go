package main

import (
	"dmkernel/dmsofts"
	"os"
)

func main() {
	// 如果是x64编译或运行 需要设置export GOARCH=386
	// 使用免注册方式注册大漠插件，也可以使用命令行注册
	// 由于大漠是付费插件，某些功能可以免费使用，但是后台高级功能需要付费
	dir, _ := os.Getwd()
	// 这里是免注册方式

	dmsofts.SetDllPathW(dir+"\\dm.dll", 0)

}
