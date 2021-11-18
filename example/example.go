/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : example.go
*   coder: zemanzeng
*   date : 2021-11-18 19:06:07
*   desc : 将文件夹下所有图片压缩
*
================================================================*/

package main

import (
	"flag"
	"log"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wenruo95/ffmpeg_relay_lib"
)

var (
	FFFmpeg string
	FInput  string
	FOutput string
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

func main() {

	t := time.Now()
	flag.StringVar(&FFFmpeg, "ffmpeg", "", "-ffmpeg=ffmpeg bin path")
	flag.StringVar(&FInput, "input", "", "-input=input dir")
	flag.StringVar(&FOutput, "output", "", "-input=output dir")
	flag.Parse()

	if len(FFFmpeg) == 0 || len(FInput) == 0 || len(FOutput) == 0 {
		log.Printf("[ERROR] invalid args. bin:%v input:%v output:%v", FFFmpeg, FInput, FOutput)
		return
	}

	entrys, err := os.ReadDir(FInput)
	if err != nil {
		log.Printf("[ERROR] read dir:%v error:%v", FInput, err)
		return
	}

	var wg sync.WaitGroup
	var count int32
	for index, entry := range entrys {
		if entry.IsDir() {
			continue
		}

		wg.Add(1)
		atomic.AddInt32(&count, 1)

		go func(ent os.DirEntry) {
			defer wg.Done()

			iPath := path.Join(FInput, ent.Name())
			oPath := path.Join(FOutput, ent.Name())
			if err := ffmpeg_relay_lib.CompressImage(FFFmpeg, iPath, oPath, 3); err != nil {
				log.Printf("[ERROR] compress_image ipath:%v opath:%v error:%v", iPath, oPath, err)
				return
			}

			log.Printf("[INFO] compress_image ipath:%v opath:%v succ:%v", iPath, oPath, err)
		}(entry)

		if index%10 == 0 {
			wg.Wait()
			log.Printf("[DEBUG] wait index:%v", index)
		}

	}
	wg.Wait()

	log.Printf("[INFO] finished %v image comrepss, consume:%v", count, time.Since(t))
}
