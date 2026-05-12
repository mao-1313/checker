package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheck_OK(t *testing.T) {
	// テスト用のサーバーを立てる
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	fmt.Println("test start", srv.URL)

	result := check(srv.URL, time.Second*2)

	if result.Err != nil {
		t.Fatalf("error should be nil: %v", result.Err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("status code should be %d: %d", http.StatusOK, result.StatusCode)
	}
}

func TestCheck_NotFound(t *testing.T) {
	// テスト用のサーバーを立てる
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	result := check(srv.URL, time.Second*2)

	if result.Err != nil {
		t.Fatalf("error should be nil: %v", result.Err)
	}
	if result.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", result.StatusCode)
	}
}

func TestCheck_Error(t *testing.T) {
	result := check("http://invalid.invalid", time.Second*2)
	if result.Err == nil {
		t.Fatalf("expecter error, got nil")
	}
}

func TestCheck_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // タイムアウトよりも長い時間スリープする
	}))
	defer srv.Close()

	result := check(srv.URL, time.Millisecond*2)
	if result.Err == nil {
		t.Error("expected timeout error")
	}
	if !errors.Is(result.Err, context.DeadlineExceeded) {
		// タイムアウト
		t.Errorf("expected DeadlineExceeded. got:%v", result.Err)
	}
}
