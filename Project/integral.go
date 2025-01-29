package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// Fungsi untuk menghitung integral tentu menggunakan metode trapesium
func integralTentu(f func(float64) float64, a, b, h float64) float64 {
	n := int((b - a) / h) // Jumlah pembagian interval
	sum := (f(a) + f(b)) / 2.0

	for i := 1; i < n; i++ {
		sum += f(a + float64(i)*h)
	}

	return sum * h
}

// Fungsi untuk mengeksekusi ekspresi matematika yang diberikan pengguna
func parseFunction(expression string) (func(float64) float64, error) {
	// Menghapus spasi ekstra di sekitar ekspresi
	expression = strings.ReplaceAll(expression, " ", "")

	// Menangani ekspresi dasar seperti x*x, x+x, sin(x), dan lainnya
	return func(x float64) float64 {
		// Fungsi dasar dan pengenalan ekspresi sin(x), cos(x), dan lainnya
		switch {
		case strings.Contains(expression, "x*x"):
			return x * x
		case strings.Contains(expression, "sin(x)"):
			return math.Sin(x)
		case strings.Contains(expression, "cos(x)"):
			return math.Cos(x)
		case strings.Contains(expression, "tan(x)"):
			return math.Tan(x)
		case strings.Contains(expression, "e^x"):
			return math.Exp(x)
		default:
			// Menangani ekspresi dengan pangkat
			if strings.Contains(expression, "^") {
				parts := strings.Split(expression, "^")
				base, _ := strconv.ParseFloat(parts[0][1:], 64) // Mengabaikan huruf 'x'
				exponent, _ := strconv.Atoi(parts[1])
				return math.Pow(base, float64(exponent))
			}
			// Fungsi tidak dikenal
			return math.NaN()
		}
	}, nil
}

// Fungsi untuk menghitung integral tak tentu (antiderivatif) untuk fungsi sederhana
func integralTakTentu(f string) string {
	switch f {
	case "x*x":
		return "x^3 / 3"
	case "sin(x)":
		return "-cos(x)"
	case "cos(x)":
		return "sin(x)"
	case "e^x":
		return "e^x"
	case "x":
		return "x^2 / 2"
	default:
		return "Integral tak tentu tidak terdefinisi untuk fungsi ini"
	}
}

// Struct untuk menampung data yang akan dipassing ke template
type PageData struct {
	Function string
	Result   string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Memproses data yang diterima dari form
		if r.Method == "POST" {
			r.ParseForm()
			funcType := r.FormValue("function")
			integralType := r.FormValue("integral_type")

			var result string

			if integralType == "tentu" {
				// Ambil nilai untuk perhitungan integral tentu
				a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
				b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
				h, _ := strconv.ParseFloat(r.FormValue("h"), 64)

				// Parse fungsi dan hitung integral
				f, _ := parseFunction(funcType)
				integralResult := integralTentu(f, a, b, h)
				result = fmt.Sprintf("Integral tentu dari fungsi pada [%.6f, %.6f] adalah: %.6f", a, b, integralResult)
			} else if integralType == "tak tentu" {
				// Menghitung integral tak tentu
				result = fmt.Sprintf("Integral tak tentu dari fungsi %s adalah: %s", funcType, integralTakTentu(funcType))
			}
			// Kirim data ke template
			pageData := PageData{
				Function: funcType,
				Result:   result,
			}
			tmpl, _ := template.New("result").ParseFiles("template.html")
			tmpl.Execute(w, pageData)
		} else {
			// Menampilkan form pertama kali
			tmpl, _ := template.New("form").ParseFiles("template.html")
			tmpl.Execute(w, nil)
		}
	})

	// Start server on port 8080
	http.ListenAndServe(":8080", nil)
}
