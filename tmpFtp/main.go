package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// 文件处理函数，提供文件下载
func fileHandler(w http.ResponseWriter, r *http.Request) {
	// 获取文件路径
	filePath := r.URL.Path[1:] // 去掉最前面的'/'
	absPath, err := filepath.Abs(filePath)

	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	// 设置文件响应头并发送文件
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(absPath)))
	http.ServeFile(w, r, absPath)
}

// 目录处理函数，提供目录内容展示
func dirHandler(w http.ResponseWriter, r *http.Request) {
	// 获取请求的目录路径
	dirPath := r.URL.Path[1:] // 去掉最前面的'/'
	absPath, err := filepath.Abs(dirPath)

	if err != nil {
		http.Error(w, "Directory not found.", http.StatusNotFound)
		return
	}

	// 检查目录是否存在
	if info, err := os.Stat(absPath); os.IsNotExist(err) || !info.IsDir() {
		http.Error(w, "Directory not found.", http.StatusNotFound)
		return
	}

	// 获取目录下的文件和子目录列表
	files, err := os.ReadDir(absPath)
	if err != nil {
		http.Error(w, "Unable to read directory.", http.StatusInternalServerError)
		return
	}

	// 构建HTML响应，展示目录内容
	fmt.Fprintf(w, "<h1>Index of %s</h1>", r.URL.Path)
	fmt.Fprintf(w, "<ul>")
	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(r.URL.Path, fileName)
		if file.IsDir() {
			fmt.Fprintf(w, "<li><a href=\"%s/\">%s/</a></li>", filePath, fileName)
		} else {
			fmt.Fprintf(w, "<li><a href=\"%s\">%s</a></li>", filePath, fileName)
		}
	}
	fmt.Fprintf(w, "</ul>")
}

func main() {
	// 启动HTTP服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/files", http.StatusFound) // 重定向到文件目录
		} else {
			filePath := r.URL.Path[1:]
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				http.Error(w, "Path error.", http.StatusInternalServerError)
				return
			}

			// 判断是文件还是目录
			fileInfo, err := os.Stat(absPath)
			if os.IsNotExist(err) {
				http.Error(w, "File not found.", http.StatusNotFound)
				return
			}

			if fileInfo.IsDir() {
				dirHandler(w, r) // 处理目录展示
			} else {
				fileHandler(w, r) // 处理文件下载
			}
		}
	})

	port := "0.0.0.0:80"
	fmt.Printf("Starting server on http://localhost%s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
