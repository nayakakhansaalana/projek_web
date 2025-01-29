package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Fungsi untuk menghitung turunan pertama (derivatif pertama)
func derivative(f func(float64) float64, x float64, h float64) float64 {
	return (f(x+h) - f(x)) / h
}

// Fungsi untuk menghitung turunan tingkat tinggi
func higherOrderDerivative(f func(float64) float64, x float64, h float64, n int) float64 {
	// Menghitung turunan tingkat tinggi dengan menghitung turunan pertama sebanyak n kali
	deriv := f
	for i := 0; i < n; i++ {
		deriv = func(x float64) float64 {
			return (deriv(x+h) - deriv(x)) / h
		}
	}
	return deriv(x)
}

// Fungsi untuk mengeksekusi ekspresi matematika yang diberikan pengguna
func parseFunction(expression string) (func(float64) float64, error) {
	// Menghapus spasi ekstra di sekitar ekspresi
	expression = strings.ReplaceAll(expression, " ", "")

	// Menangani ekspresi dasar seperti x*x, x+x, sin(x), dan lainnya
	return func(x float64) float64 {
		// Memeriksa apakah ekspresi mengandung operasi pangkat
		if strings.Contains(expression, "^") {
			parts := strings.Split(expression, "^")
			base, _ := strconv.ParseFloat(parts[0][1:], 64) // Mengabaikan huruf 'x'
			exponent, _ := strconv.Atoi(parts[1])
			return math.Pow(base, float64(exponent))
		}

		// Fungsi-fungsi dasar
		switch expression {
		case "x*x":
			return x * x
		case "x+x":
			return x + x
		case "x/x":
			if x == 0 {
				return math.NaN() // Mencegah pembagian dengan nol
			}
			return x / x
		case "sin(x)":
			return math.Sin(x)
		case "cos(x)":
			return math.Cos(x)
		case "tan(x)":
			return math.Tan(x)
		default:
			// Fungsi tidak dikenal
			return math.NaN()
		}
	}, nil
}

// Fungsi untuk menghitung limit saat x mendekati tak hingga
func limitToInfinity(f func(float64) float64, start float64) float64 {
	// Mulai dari titik yang jauh, dan lakukan iterasi menuju tak hingga
	for i := 0; i < 1000; i++ {
		start *= 10 // Meningkatkan nilai x untuk mendekati tak hingga
		result := f(start)
		if math.IsNaN(result) || math.IsInf(result, 0) {
			return result // Jika hasilnya tidak terdefinisi atau tak hingga
		}
	}
	return f(start)
}

// Fungsi untuk menghitung limit di titik tertentu
func limitAtPoint(f func(float64) float64, point float64) float64 {
	// Menghitung nilai limit dari sisi kiri dan kanan titik tersebut
	epsilon := 0.0000001 // Nilai yang sangat kecil
	left := f(point - epsilon)
	right := f(point + epsilon)

	if math.Abs(left-right) < epsilon { // Jika keduanya hampir sama, maka limit ada
		return left
	}
	return math.NaN() // Jika limit tidak ada atau tidak terdefinisi
}

// Handler untuk menghitung limit
func limitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		FuncType string `json:"funcType"`
		Input    string `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	f, err := parseFunction(data.FuncType)
	if err != nil {
		http.Error(w, "Invalid function expression", http.StatusBadRequest)
		return
	}

	var result float64
	var errMsg string

	if strings.ToLower(data.Input) == "infinity" {
		start := 1000.0
		result = limitToInfinity(f, start)
	} else {
		point, err := strconv.ParseFloat(data.Input, 64)
		if err != nil {
			errMsg = "Invalid input point."
		} else {
			result = limitAtPoint(f, point)
		}
	}

	response := map[string]interface{}{
		"result": result,
		"error":  errMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler untuk menghitung turunan
func derivativeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		FuncType string  `json:"funcType"`
		X        float64 `json:"x"`
		H        float64 `json:"h"`
		N        int     `json:"n,omitempty"`
		Type     string  `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	f, err := parseFunction(data.FuncType)
	if err != nil {
		http.Error(w, "Invalid function expression", http.StatusBadRequest)
		return
	}

	var result float64
	var errMsg string

	if data.Type == "first" {
		result = derivative(f, data.X, data.H)
	} else if data.Type == "higher" {
		result = higherOrderDerivative(f, data.X, data.H, data.N)
	} else {
		errMsg = "Invalid derivative type."
	}

	response := map[string]interface{}{
		"result": result,
		"error":  errMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Fungsi utama
func main() {
	http.HandleFunc("/calculate", limitHandler)
	http.HandleFunc("/derivative", derivativeHandler)
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
