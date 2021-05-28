package ca_game_api

import "math"

func calculateLevel(experience int) int {
	return int(math.Floor(math.Sqrt(float64(experience)) / 10.0))
}

func calculatePower(experience, basePower int) int {
	return calculateLevel(experience) * basePower
}
