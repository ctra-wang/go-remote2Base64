package r2b

import (
	"fmt"
	"testing"
)

func TestGetRemoteConvertBase64(t *testing.T) {
	imgUrl := "https://img-blog.csdnimg.cn/350174dad84e4ecaa0f7995207791df9.jpeg"
	tff := "longshuhongheicuti.ttf"
	imgName := "test.jpg"
	isDelImg := true
	str, err := GetRemoteConvertBase64(imgUrl, tff, imgName, isDelImg, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(str)
}
