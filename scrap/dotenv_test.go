package main

import (
	"os"
	"testing"
)
func TestLoadDotEnv(t *testing.T) {
	// 創建臨時的 .env 文件
	tmpfile, err := os.CreateTemp("", ".env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // 清理

	content := []byte("ROD_URL=test_rod_url\nJOBCAT=test_jobcat\nDOMAIN=test_domain")
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// 執行函數
	loadDotEnv(tmpfile.Name())

	// 驗證環境變數是否被正確設定
	if rodURL != "test_rod_url" || jobCat != "test_jobcat" || domain != "test_domain" {
		t.Errorf("Environment variables not set correctly")
	}
}