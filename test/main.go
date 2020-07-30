package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func main() {
	url := "http://localhost:8085/api/beta/attend-temp"
	fmt.Println("URL:>", url)
	for i:=0; i <30; i++{
		time.Sleep(500)
		go func() {
			fmt.Println(i)
			jsonString := fmt.Sprintf(`{"batchID":"%d","camera":{"host":"192.168.1.93","port":"8080","username":"","pass":""},"faceIDs":[1212]}`,i)
			var jsonStr = []byte(jsonString)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			resp.Body.Close()
		}()

	}

}