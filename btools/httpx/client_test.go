package httpx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRestClientGet проверяет GET-запрос.
func TestRestClientGet(t *testing.T) {
	// Создаем тестовый сервер, который отвечает на GET.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Ожидалось метод GET, а пришло: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("get response"))
	}))
	defer ts.Close()

	client := NewRestClient(ts.URL, nil, nil)
	resp, err := client.Get(context.Background(), "", nil, 5*time.Second, 1)
	if err != nil {
		t.Fatalf("Ошибка GET-запроса: %v", err)
	}
	if resp.String() != "get response" {
		t.Errorf("Ожидался ответ 'get response', а получили: %s", resp.String())
	}
}

// TestRestClientPost проверяет POST-запрос.
func TestRestClientPost(t *testing.T) {
	// Создаем тестовый сервер для POST.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Ожидалось метод POST, а пришло: %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Не удалось прочитать тело запроса: %v", err)
		}
		if string(body) != "post data" {
			t.Errorf("Ожидалось тело 'post data', а получили: %s", string(body))
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("post response"))
	}))
	defer ts.Close()

	client := NewRestClient(ts.URL, nil, nil)
	resp, err := client.Post(context.Background(), "", nil, "post data", 5*time.Second)
	if err != nil {
		t.Fatalf("Ошибка POST-запроса: %v", err)
	}
	if resp.String() != "post response" {
		t.Errorf("Ожидался ответ 'post response', а получили: %s", resp.String())
	}
}

// TestRestClientPut проверяет PUT-запрос.
func TestRestClientPut(t *testing.T) {
	// Тестовый сервер для PUT.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Ожидалось метод PUT, а пришло: %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Не удалось прочитать тело запроса: %v", err)
		}
		if string(body) != "put data" {
			t.Errorf("Ожидалось тело 'put data', а получили: %s", string(body))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("put response"))
	}))
	defer ts.Close()

	client := NewRestClient(ts.URL, nil, nil)
	resp, err := client.Put(context.Background(), "", nil, "put data", 5*time.Second)
	if err != nil {
		t.Fatalf("Ошибка PUT-запроса: %v", err)
	}
	if resp.String() != "put response" {
		t.Errorf("Ожидался ответ 'put response', а получили: %s", resp.String())
	}
}

// TestRestClientPatch проверяет PATCH-запрос.
func TestRestClientPatch(t *testing.T) {
	// Тестовый сервер для PATCH.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Ожидалось метод PATCH, а пришло: %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Не удалось прочитать тело запроса: %v", err)
		}
		if string(body) != "patch data" {
			t.Errorf("Ожидалось тело 'patch data', а получили: %s", string(body))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("patch response"))
	}))
	defer ts.Close()

	client := NewRestClient(ts.URL, nil, nil)
	resp, err := client.Patch(context.Background(), "", nil, "patch data", 5*time.Second)
	if err != nil {
		t.Fatalf("Ошибка PATCH-запроса: %v", err)
	}
	if resp.String() != "patch response" {
		t.Errorf("Ожидался ответ 'patch response', а получили: %s", resp.String())
	}
}

// TestRestClientDelete проверяет DELETE-запрос.
func TestRestClientDelete(t *testing.T) {
	// Тестовый сервер для DELETE.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Ожидалось метод DELETE, а пришло: %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delete response"))
	}))
	defer ts.Close()

	client := NewRestClient(ts.URL, nil, nil)
	resp, err := client.Delete(context.Background(), "", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("Ошибка DELETE-запроса: %v", err)
	}
	if resp.String() != "delete response" {
		t.Errorf("Ожидался ответ 'delete response', а получили: %s", resp.String())
	}
}
