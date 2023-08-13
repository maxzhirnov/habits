package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDateOnly(t *testing.T) {
	// Используем пакет testify для утверждений
	assertions := assert.New(t)

	// Создаем пример времени
	sampleTime := time.Date(2023, 8, 13, 15, 30, 45, 123456789, time.UTC)

	// Применяем вашу функцию
	onlyDate := DateOnly(sampleTime)

	// Проверяем, что год, месяц и день не изменились
	assertions.Equal(2023, onlyDate.Year())
	assertions.Equal(time.August, onlyDate.Month())
	assertions.Equal(13, onlyDate.Day())

	// Проверяем, что часы, минуты, секунды и наносекунды равны нулю
	assertions.Equal(0, onlyDate.Hour())
	assertions.Equal(0, onlyDate.Minute())
	assertions.Equal(0, onlyDate.Second())
	assertions.Equal(0, onlyDate.Nanosecond())

	// Проверяем, что временная зона осталась прежней
	assertions.Equal(sampleTime.Location(), onlyDate.Location())
}
