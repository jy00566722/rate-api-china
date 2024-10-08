package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	// 指定项目根目录
	root := "."
	// 指定输出文件
	outputFile := "merged_project.txt"

	// 打开输出文件
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	} else {
		fmt.Printf("Created output file: %s\n", outputFile)
	}
	defer out.Close()

	// 遍历目录
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理 .go 文件
		if !info.IsDir() && strings.HasSuffix(path, ".go") && path != "merge.go" && path != outputFile {
			// 写入文件路径作为注释
			fmt.Fprintf(out, "// File: %s\n", path)

			// 读取文件内容
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(out, f)
			if err != nil {
				return err
			}

			fmt.Fprintln(out)
			// 添加两行######作为分隔
			fmt.Fprintln(out, "###文件分割###")
			fmt.Fprintln(out)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", root, err)
		return
	}

	fmt.Printf("All .go files have been merged into %s\n", outputFile)
}
