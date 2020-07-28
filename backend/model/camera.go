package model

import (
	"fmt"
	"gocv.io/x/gocv"
	"time"
)

type Camera struct {
	Host      string        `json:"host" bson:"host"`
	Port      string        `json:"port" bson:"port"`
	Username  string        `json:"username" bson:"username"`
	Password  string        `json:"password" bson:"password"`
}

func (c *Camera) GetFrames(rootFolder string, t time.Time, numOfFrame int) ([]string, error) {
	urlStreamVideo := fmt.Sprintf("http://%s:%s", c.Host, c.Port)
	webcam, err := gocv.OpenVideoCapture(urlStreamVideo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer webcam.Close()

	img := gocv.NewMat()
	defer img.Close()

	imagePaths := make([]string, 0)

	for i := 0; i < numOfFrame; i++ {
		if ok := webcam.Read(&img); !ok {
			continue
		}

		if img.Empty() {
			continue
		}

		path := fmt.Sprintf(`%s/%v.jpg`, rootFolder, time.Now().Unix())
		gocv.IMWrite(path, img)
		imagePaths = append(imagePaths, path)
		time.Sleep(time.Second*2)
	}

	return imagePaths, nil
}
