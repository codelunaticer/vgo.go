package mosLogic

import (
	"fmt"
	"os"
	fileConst "vgoadmin-go/constant/file"
)

// InitAppStaticFinder 初始化好需要的文件目录
func InitAppStaticFinder() {
	finderList := []string{
		fileConst.RootStaticFinder,
		fileConst.AvatarFinder,
		fileConst.TempFinder,
	}

	for _, path := range finderList {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			err2 := os.Mkdir(path, os.ModePerm)
			if err2 != nil {
				fmt.Printf("%s目录创建失败:%v\n", path, err2)
				panic("目录准备失败\n")
			}
			fmt.Printf("%s目录创建成功\n", path)
		}
	}
}
