/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : image_compress_test.go
*   coder: zemanzeng
*   date : 2021-11-12 20:34:01
*   desc : 图片压缩测试用例
*
================================================================*/

package ffmpeg_relay_lib

import "testing"

func TestCompressImage(t *testing.T) {
	var (
		ffmpegPath = "/usr/local/bin/ffmpeg"
		inputPath  = "./example/" + "example.jpeg"
		outputPath = "./example/" + "example_small.jpeg"
	)

	if err := CompressImage(ffmpegPath, inputPath, outputPath, 3); err != nil {
		t.Errorf("compress image error:" + err.Error())
	}
	t.Logf("compress succ")
}
