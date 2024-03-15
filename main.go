package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const maxFileSize = int64(10 * 1024 * 1024)
const maxImageSize = int64(1 * 1024 * 1024)

const (
	DateLayout8 = "20060102"
)

func fileSave(w http.ResponseWriter, r *http.Request) {
	path := save(r)
	if path == "" {
		//	RespondError(w, errors.New("an error occurred"), http.StatusInternalServerError)
		return
	}
	// Create a JSON response
	response := struct {
		Success bool   `json:"success"`
		path    string `json:"path"`
	}{
		Success: true,
		path:    path,
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
func save(r *http.Request) string {
	// left shift 32 << 20 which results in 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return ""
	}
	n := r.Form.Get("name")
	// Retrieve the file from form data
	f, h, err := r.FormFile("file")
	if err != nil {
		return ""
	}
	defer f.Close()
	path := filepath.Join(".", "files")
	_ = os.MkdirAll(path, os.ModePerm)
	fullPath := path + "/" + n
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return ""
	}
	defer file.Close()
	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		return ""
	}
	return n + filepath.Ext(h.Filename)
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var image bool
	path := r.URL.Query().Get("path")
	// 파일을 읽어옵니다.
	// left shift 32 << 20 which results in 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	_, fileHeader, err := r.FormFile("file")
	uploadPath := fmt.Sprintf("%s/%v", path, time.Now().Format(DateLayout8))
	fmt.Println(uploadPath)
	if image {
		if fileHeader.Size > maxImageSize {
			maxSize := fmt.Sprintf("최대: %dM", maxImageSize/(1024*1024))
			http.Error(w, "파일 사이즈가 너무 큽니다."+maxSize, http.StatusBadRequest)
			return
		}
	} else {
		if fileHeader.Size > maxFileSize {
			maxSize := fmt.Sprintf("최대: %dM", maxFileSize/(1024*1024))

			http.Error(w, "파일 사이즈가 너무 큽니다."+maxSize, http.StatusBadRequest)
			return
		}
	}

	src, err := fileHeader.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer src.Close()
	//accessUrl, err := adapter.AwsS3Adapter().UploadFile(uploadPath, src, fileHeader)
	accessUrl := ""
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = url.ParseRequestURI(accessUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
func multiUploadHandler(w http.ResponseWriter, r *http.Request) {
	var image bool

	path := r.URL.Query().Get("path")
	// left shift 32 << 20 which results in 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	// Retrieve the files from form data
	files := r.MultipartForm.File["files"]
	uploadPath := fmt.Sprintf("%s/%v", path, time.Now().Format(DateLayout8))
	fmt.Println(uploadPath)
	if files == nil || len(files) == 0 {
		msg := "파일이 없습니다"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	fileUrls := make([]string, len(files))

	for i, file := range files {
		if image {
			if file.Size > maxImageSize {
				maxSize := fmt.Sprintf("최대: %dM", maxImageSize/(1024*1024))
				http.Error(w, "파일 사이즈가 너무 큽니다."+maxSize, http.StatusBadRequest)
				return
			}
		} else {
			if file.Size > maxFileSize {
				maxSize := fmt.Sprintf("최대: %dM", maxFileSize/(1024*1024))

				http.Error(w, "파일 사이즈가 너무 큽니다."+maxSize, http.StatusBadRequest)
				return
			}
		}

		src, err := file.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer src.Close()
		//accessUrl, err := adapter.AwsS3Adapter().UploadFile(uploadPath, src, file)
		accessUrl := ""
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = url.ParseRequestURI(accessUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fileUrls[i] = accessUrl
	}

	resp := make(map[string][]string)
	resp["accessUrl"] = fileUrls
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/multi/upload", multiUploadHandler)
	http.HandleFunc("/save", fileSave)
	http.ListenAndServe(":8080", nil)
}
