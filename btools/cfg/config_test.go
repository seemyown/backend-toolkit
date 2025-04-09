package cfg

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfig – структура для тестового конфига.
type TestConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func TestNewConfig_Success(t *testing.T) {
	// Создаем временную директорию.
	tempDir, err := os.MkdirTemp("", "testconfig")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Пишем тестовый YAML-конфиг в файл "testconfig.yaml".
	configContent := `host: "localhost"
port: 8080`
	configFilePath := filepath.Join(tempDir, "testconfig.yaml")
	if err := os.WriteFile(configFilePath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Не удалось записать конфиг файл: %v", err)
	}

	// Вызываем функцию NewConfig: имя файла без расширения, расширение и путь.
	conf := NewConfig[TestConfig]("testconfig", "yaml", tempDir)
	if conf == nil {
		t.Fatalf("Ожидался не nil-результат, а получили nil")
	}

	if conf.Host != "localhost" {
		t.Errorf("Ожидалось, что Host = 'localhost', а получили '%s'", conf.Host)
	}
	if conf.Port != 8080 {
		t.Errorf("Ожидалось, что Port = 8080, а получили %d", conf.Port)
	}
}

func TestNewConfig_Error(t *testing.T) {
	// Попытка загрузить конфиг из несуществующей директории.
	conf := NewConfig[TestConfig]("nonexistent", "yaml", "/nonexistent/directory")
	if conf != nil {
		t.Errorf("Ожидался nil, но получили не nil")
	}
}
