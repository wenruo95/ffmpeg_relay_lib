/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : image_compress.go
*   coder: zemanzeng
*   date : 2021-11-12 19:46:37
*   desc : 图片压缩
*
================================================================*/

package ffmpeg_relay_lib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func CompressImage(ffmpegPath string, inputPath string, outputPath string, timeout int64) error {
	begin := time.Now()

	logname := fmt.Sprintf("image_compress.%s.log", begin.Format("20060102150405"))
	logout, err := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open %s error:%w", logname, err)
	}
	defer logout.Close()

	// ffmpeg -hide_banner -i in.png -pix_fmt pal8 out_1.png
	param := fmt.Sprintf("%s -hide_banner -i %s -pix_fmt pal8 %s", ffmpegPath, inputPath, outputPath)
	cmd, err := ExecCmd(param, logout, logout)
	if err != nil {
		return fmt.Errorf("exec cmd:%v error:%w", param, err)
	}

	if timeout > 0 {
		select {
		case <-time.After(time.Second * time.Duration(timeout)):
			// wait
		}
	}

	if err := KillCmd(cmd); err != nil {
		return fmt.Errorf("kill cmd:%v error:%v", param, err)
	}
	return nil
}

func ExecCmd(param string, stdout io.Writer, stderr io.Writer) (*exec.Cmd, error) {
	cmd := exec.Command("/bin/bash", "-c", param)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 将标准输出和标准错误都写到log中
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func KillCmd(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return errors.New("process not found or already stopped")
	}

	if err := syscall.Kill(cmd.Process.Pid, syscall.SIGKILL); err != nil {
		return err
	}

	// 如果不wait，则会产生僵尸进程
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
