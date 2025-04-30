package security

import (
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"testing"
	"time"
)

// TestIssure проверяет, что функция issure возвращает непустое имя пользователя.
func TestIssure(t *testing.T) {
	username := issure()
	if username == "" {
		t.Error("issure() вернула пустую строку, хотя, видно, даже система не смогла найти юзера")
	}
}

// TestGenerateJWTToken проверяет генерацию JWT-токена, его подпись и корректность claim-ов.
func TestGenerateJWTToken(t *testing.T) {
	// Подготовка данных для токена
	tokenPayload := map[string]interface{}{
		"user_id": 123,
	}
	tokenSettings := map[string]interface{}{
		"role": "admin",
	}
	secretKey := "testsecret"
	// Используем стандартный метод HS256
	signMethod := jwt.SigningMethodHS256

	// Генерация токена
	tokenStr, err := GenerateJWTToken(tokenPayload, tokenSettings, secretKey, signMethod)
	if err != nil {
		t.Fatalf("GenerateJWTToken вернула ошибку: %v", err)
	}

	// Разбор и проверка токена
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		t.Fatalf("Ошибка при разборе токена: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Токен не валиден — неужели подпись оказалась фальшивкой?")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Не удалось преобразовать claims к jwt.MapClaims")
	}

	// Проверка наличия дефолтных полей
	for _, field := range []string{"exp", "iat", "nbf", "iss"} {
		if _, exists := claims[field]; !exists {
			t.Errorf("Ожидается наличие claim '%s'", field)
		}
	}

	// Для сравнения поля "iss" получаем значение через issure()
	expectedIss := issure()
	if claims["iss"] != expectedIss {
		t.Errorf("Ожидалось, что iss='%s', а получили '%v'", expectedIss, claims["iss"])
	}

	// Дополнительные проверки для tokenSettings и tokenPayload
	if role, exists := claims["role"]; !exists || role != "admin" {
		t.Errorf("Ожидается, что role будет 'admin', а получили '%v'", claims["role"])
	}

	// Обрати внимание: JSON преобразует числовые значения в float64.
	if userID, exists := claims["user_id"]; !exists || !reflect.DeepEqual(userID, float64(123)) {
		t.Errorf("Ожидалось, что user_id будет 123, а получили '%v'", claims["user_id"])
	}

	// Дополнительно: проверим, что время истечения не прошло (если тесты выполняются моментально)
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Error("claim exp не является числом")
	} else {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			t.Error("Токен уже истёк — время шуток кончилось")
		}
	}
}
