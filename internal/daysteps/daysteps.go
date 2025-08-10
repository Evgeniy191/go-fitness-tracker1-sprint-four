package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")

	// Проверяем количество частей
	if len(parts) != 2 {

		return 0, 0, errors.New("должно быть две части")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	// Парсим количество шагов (первый элемент)
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, errors.New("неправильное количество шагов")
	}
	if steps <= 0 {

		return 0, 0, errors.New("количество шагов должно быть больше нуля")
	}

	// Парсим продолжительность (второй элемент) через time.ParseDuration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {

		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {

		return ""
	}
	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm
	calories, _ := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm,
		calories)

	return result
}
