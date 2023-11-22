package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestInitialLog(t *testing.T) {
	// 保存原來的標準輸出
	origStdout := os.Stdout
	defer func() { os.Stdout = origStdout }()

	// 創建一個緩衝區來捕獲標準輸出
	r, w, _ := os.Pipe()
	os.Stdout = w

	initialLog()

	// 關閉寫入端，讀取輸出內容
	w.Close()
	out, _ := io.ReadAll(r)

	// 驗證標準輸出中的日誌內容
	if !bytes.Contains(out, []byte("Initialize logger success.")) {
		t.Errorf("expected 'Initialize logger success.' to be written to stdout, got %s", string(out))
	}

	// 驗證日誌文件是否創建
	if _, err := os.Stat("log.txt"); os.IsNotExist(err) {
		t.Errorf("log file 'log.txt' does not exist")
	}

	// 清理測試產生的文件
	os.Remove("log.txt")
}