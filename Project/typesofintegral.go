package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

// Fungsi untuk menghitung integral tentu 1D menggunakan metode trapesium
func integral1D(f func(float64) float64, a, b, h float64) float64 {
	n := int((b - a) / h) // Jumlah pembagian interval
	sum := (f(a) + f(b)) / 2.0
	for i := 1; i < n; i++ {
		sum += f(a + float64(i)*h)
	}
	return sum * h
}

// Fungsi untuk menghitung integral ganda menggunakan metode trapesium
func integral2D(f func(float64, float64) float64, xMin, xMax, yMin, yMax, h float64) float64 {
	var sum float64
	for x := xMin; x <= xMax; x += h {
		for y := yMin; y <= yMax; y += h {
			sum += f(x, y)
		}
	}
	return sum * h * h
}

// Fungsi untuk menghitung integral triple menggunakan metode trapesium
func integral3D(f func(float64, float64, float64) float64, xMin, xMax, yMin, yMax, zMin, zMax, h float64) float64 {
	var sum float64
	for x := xMin; x <= xMax; x += h {
		for y := yMin; y <= yMax; y += h {
			for z := zMin; z <= zMax; z += h {
				sum += f(x, y, z)
			}
		}
	}
	return sum * h * h * h
}

// Fungsi untuk parsing dan mengevaluasi ekspresi matematika dalam format string
func parseExpression(expr string, dim int) (interface{}, error) {
	switch expr {
	case "x*x":
		if dim == 2 {
			return func(x, y float64) float64 { return x * x }, nil
		} else if dim == 3 {
			return func(x, y, z float64) float64 { return x * x }, nil
		}
	case "sin(x)":
		if dim == 2 {
			return func(x, y float64) float64 { return math.Sin(x) }, nil
		} else if dim == 3 {
			return func(x, y, z float64) float64 { return math.Sin(x) }, nil
		}
	case "cos(x)":
		if dim == 2 {
			return func(x, y float64) float64 { return math.Cos(x) }, nil
		} else if dim == 3 {
			return func(x, y, z float64) float64 { return math.Cos(x) }, nil
		}
	case "e^x":
		if dim == 2 {
			return func(x, y float64) float64 { return math.Exp(x) }, nil
		} else if dim == 3 {
			return func(x, y, z float64) float64 { return math.Exp(x) }, nil
		}
	case "x*y*z": // Contoh untuk fungsi 3D
		if dim == 3 {
			return func(x, y, z float64) float64 { return x * y * z }, nil
		}
	case "x*y":
		if dim == 2 {
			return func(x, y float64) float64 { return x * y }, nil
		} else if dim == 3 {
			return func(x, y, z float64) float64 { return x * y * z }, nil
		}
	default:
		return nil, fmt.Errorf("fungsi tidak dikenali: %s", expr)
	}
	return nil, fmt.Errorf("fungsi tidak dikenali: %s", expr)
}

// Handler untuk menghitung integral
func calculateIntegralHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	r.ParseForm()

	integralType := r.FormValue("integralType")
	xMin, _ := strconv.ParseFloat(r.FormValue("xMin"), 64)
	xMax, _ := strconv.ParseFloat(r.FormValue("xMax"), 64)
	yMin, _ := strconv.ParseFloat(r.FormValue("yMin"), 64)
	yMax, _ := strconv.ParseFloat(r.FormValue("yMax"), 64)
	h, _ := strconv.ParseFloat(r.FormValue("h"), 64)
	function := r.FormValue("function")

	// Tentukan dimensi berdasarkan integralType
	var dim int
	if integralType == "double" {
		dim = 2
	} else if integralType == "triple" {
		dim = 3
	} else {
		http.Error(w, "Tipe integral tidak dikenali, harap pilih 'ganda' atau 'triple'.", http.StatusBadRequest)
		return
	}

	// Parse function expression dengan dimensi yang sesuai
	var f interface{}
	var err error
	f, err = parseExpression(function, dim)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result float64
	if integralType == "double" {
		result = integral2D(f.(func(float64, float64) float64), xMin, xMax, yMin, yMax, h)
	} else if integralType == "triple" {
		zMin, _ := strconv.ParseFloat(r.FormValue("zMin"), 64)
		zMax, _ := strconv.ParseFloat(r.FormValue("zMax"), 64)
		result = integral3D(f.(func(float64, float64, float64) float64), xMin, xMax, yMin, yMax, zMin, zMax, h)
	}

	// Return result as JSON
	response := map[string]interface{}{"result": result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Serve static files (e.g., HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Handle the integral calculation requests
	http.HandleFunc("/calculate", calculateIntegralHandler)

	// Start the web server
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
