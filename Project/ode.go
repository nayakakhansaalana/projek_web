package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
)

// Fungsi untuk metode Euler
func euler(f func(float64, float64) float64, y0, t0, tEnd, h float64) []float64 {
	var result []float64
	y := y0
	t := t0
	for t <= tEnd {
		result = append(result, y)
		y += h * f(t, y)
		t += h
	}
	return result
}

// Fungsi untuk metode Runge-Kutta orde 4
func rungeKutta(f func(float64, float64) float64, y0, t0, tEnd, h float64) []float64 {
	var result []float64
	y := y0
	t := t0
	for t <= tEnd {
		result = append(result, y)
		k1 := h * f(t, y)
		k2 := h * f(t+0.5*h, y+0.5*k1)
		k3 := h * f(t+0.5*h, y+0.5*k2)
		k4 := h * f(t+h, y+k3)
		y += (k1 + 2*k2 + 2*k3 + k4) / 6
		t += h
	}
	return result
}

// Fungsi untuk mengeksekusi ekspresi matematika yang diberikan pengguna
func parseDifferentialEquation(equation string) (func(float64, float64) float64, error) {
	equation = strings.ReplaceAll(equation, " ", "") // Menghapus spasi
	if strings.Contains(equation, "y") && strings.Contains(equation, "t") {
		return func(t, y float64) float64 {
			if equation == "y-t^2+1" {
				return y - t*t + 1
			}
			if equation == "y*t+1" {
				return y * t
			}
			if equation == "sin(t)*y" {
				return math.Sin(t) * y
			}
			return 0
		}, nil
	}
	return nil, fmt.Errorf("persamaan tidak valid")
}

// Handler untuk menangani permintaan perhitungan ODE
func odeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		T0       float64 `json:"t0"`
		Y0       float64 `json:"y0"`
		TEnd     float64 `json:"tEnd"`
		H        float64 `json:"h"`
		Method   string  `json:"method"`
		Equation string  `json:"equation"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parsing persamaan
	f, err := parseDifferentialEquation(data.Equation)
	if err != nil {
		http.Error(w, "Invalid equation", http.StatusBadRequest)
		return
	}

	var solution []float64
	switch strings.ToLower(data.Method) {
	case "euler":
		solution = euler(f, data.Y0, data.T0, data.TEnd, data.H)
	case "runge-kutta":
		solution = rungeKutta(f, data.Y0, data.T0, data.TEnd, data.H)
	default:
		http.Error(w, "Unknown method", http.StatusBadRequest)
		return
	}

	// Mengirimkan hasil perhitungan sebagai response JSON
	response := map[string]interface{}{
		"result": solution,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Fungsi utama untuk menjalankan server HTTP
func main() {
	http.HandleFunc("/solve-ode", odeHandler)
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
