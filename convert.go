package r2b

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/fogleman/gg"
	"github.com/polds/imgbase64"
	"image/jpeg"
	"io"
	"os"
	"strings"
)

// GetRemoteConvertBase64 this func accpet the remote imgUrl,and used the lib "github.com/polds/imgbase64" handle the img
// params：
//
//	1、imgUrl：public imgUrl  (on the web wo can open it) without credentials.
//	2、tff：if you handle the img,you can upload yourself tff.
//	3、imgName：the local img,if you get the remote and create the local img file.
//	4、isDelImg：if isDelImg is true,del the create image,else don`t del.
//	5、PointInfo：support to drawing the charaters on the image
//
// and return a base64 image string and error
func GetRemoteConvertBase64(imgUrl string, tff string, imgName string, isDelImg bool, points []PointInfo) (string, error) {
	// get remote imageUrl
	// the Problem:we can`t get the img suffix，if you want to get the img suffix,you can refer to the [convertEnhance.go]
	img := imgbase64.FromRemote(imgUrl)
	if len(img) <= 0 {
		return "error", errors.New("place check your imgUrl policy is public read")
	}
	// use the "imgbase64" the return had a prefix,eg:"data:image/png;base64,",go we need to split the result
	i := strings.Index(img, ",")
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(img[i+1:]))

	// create the local image with the remote base64 string，if your used the docker，you need to input the correct path
	f, err := os.Create(imgName)
	if err != nil {
		return "error", errors.New("create the image error")
	}
	// if isDelImg is true,del the create image,else don`t del.
	// if you always drawing the one image with other info,you don`t need to del the image
	if isDelImg {
		defer os.RemoveAll(f.Name())
	}
	_, err = io.Copy(f, dec)
	if err != nil {
		return "error", errors.New("copy the image error")
	}
	// drawing the local image,used the lib "github.com/fogleman/gg"
	// imgName eg:"xxx.jpg" or "xxx.png",you don`t care about the origin remote img suffix
	waterImage, err := gg.LoadImage(imgName)
	// if can`t load the image,return the err,and in if create the file
	if err != nil {
		return "error", errors.New("used the lib gg,load the image name error")
	}
	// used the gg lib,create the water image
	dc := gg.NewContextForImage(waterImage)
	dc.SetRGBA(1, 1, 1, 0)
	dc.SetRGB(0, 0, 0)
	// load the tff,the path must be local path
	if err := dc.LoadFontFace(tff, 16); err != nil {
		return "error", errors.New("used the lib gg,load the tff error")
	}
	if len(points) > 0 {
		// drawing the image with characters
		for _, point := range points {
			var align gg.Align
			switch point.Align {
			case 1:
				align = gg.AlignLeft
			case 2:
				align = gg.AlignCenter
			case 3:
				align = gg.AlignRight
			}
			dc.DrawStringWrapped(point.Info, point.X, point.Y, 0, 0, point.With, 1, align)
		}
	}

	buffer := bytes.NewBuffer(nil)
	// save the drawing image，and control the quality
	err = jpeg.Encode(buffer, dc.Image(), &jpeg.Options{Quality: 60})
	b := buffer.Bytes()
	str := base64.StdEncoding.EncodeToString(b)
	return str, nil
}
