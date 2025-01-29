package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Fungsi untuk menghitung limit menuju tak hingga
func limitToInfinity(f func(float64) float64, start float64) float64 {
	step := 10000.0
	x := start
	for i := 0; i < 100; i++ {
		x += step
		result := f(x)
		if math.IsInf(result, 1) || math.IsInf(result, -1) {
			return result
		}
	}
	return f(x)
}

// Fungsi untuk menghitung limit menuju nilai tertentu
func limitAtPoint(f func(float64) float64, point float64) float64 {
	h := 1e-6
	left := f(point - h)
	right := f(point + h)

	if math.IsInf(left, 1) || math.IsInf(left, -1) || math.IsInf(right, 1) || math.IsInf(right, -1) {
		return math.Inf(1)
	}

	if math.Abs(left-right) < 1e-6 {
		return left
	}
	return math.NaN()
}

// Fungsi untuk mengeksekusi ekspresi matematika yang diberikan pengguna
func parseFunction(expression string) (func(float64) float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	return func(x float64) float64 {
		if strings.Contains(expression, "^") {
			parts := strings.Split(expression, "^")
			base, err1 := strconv.ParseFloat(parts[0][1:], 64) // Mengabaikan huruf 'x'
			exponent, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				return math.NaN()
			}
			return math.Pow(base, float64(exponent))
		}

		// Mendukung ekspresi dasar
		switch expression {
		case "x*x":
			return x * x
		case "x+x":
			return x + x
		case "x/x":
			if x == 0 {
				return math.NaN()
			}
			return 1 // x/x = 1 untuk semua x != 0
		default:
			return math.NaN()
		}
	}, nil
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

	// Return the result and error message in JSON format
	response := map[string]interface{}{
		"result": result,
		"error":  errMsg,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Fungsi utama untuk melayani file HTML
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html") // Pastikan file `index.html` ada di root folder proyek
}

func main() {
	// Menyajikan file statis seperti CSS atau gambar
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routing handler untuk halaman root dan perhitungan limit
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/calculate", limitHandler)

	fmt.Println("Server running on http://localhost:8080")
	// Memulai server pada port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
