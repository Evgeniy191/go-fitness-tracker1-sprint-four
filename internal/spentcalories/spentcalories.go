package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("должно быть три части")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {

	lengthStep := height * stepLengthCoefficient
	distanceMeters := float64(steps) * lengthStep
	return distanceMeters / mInKm
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	if steps <= 0 {
		return "", errors.New("количество шагов должно быть больше нуля")
	}

	var distKm float64
	var speed float64
	var calories float64

	switch activityType {
	case "Бег":
		distKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		caloriesVal, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		calories = caloriesVal

	case "Ходьба":
		distKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		caloriesVal, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		calories = caloriesVal

	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result := "Тип тренировки: " + activityType + "\n" +
		fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours()) +
		fmt.Sprintf("Дистанция: %.2f км.\n", distKm) +
		fmt.Sprintf("Скорость: %.2f км/ч\n", speed) +
		fmt.Sprintf("Сожгли калорий: %.2f\n", calories)

	return result, nil
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distanceKilometers := distance(steps, height)
	hours := duration.Hours()
	return distanceKilometers / hours

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неправильные входные данные")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	spentcalories := (weight * avgSpeed * minutes) / minInH

	return spentcalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("неправильные входные данные")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()
	spentcalories := (weight * avgSpeed * minutes) / minInH
	return spentcalories * walkingCaloriesCoefficient, nil

}
