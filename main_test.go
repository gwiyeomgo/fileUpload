package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_uploadHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/upload"), nil)
	req.Header.Add("Content-Type", "application/json")

}

func Test_multiUploadHandler(t *testing.T) {

}

func Test_fileSave(t *testing.T) {

}
