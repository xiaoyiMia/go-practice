package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Run with
//		go run .
// Send request with:
//		curl -F 'file=@/path/matrix.csv' "localhost:8080/echo"

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		records, err := readFile(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}

		response := buildeResponse(records)
		fmt.Fprint(w, response)
	})

	/*
		Invert matrix. Example:
			1,2,3       1,4,7
			4,5,6  =>   2,5,8
			7,8,9       3,6,9
	*/
	http.HandleFunc("/invert", func(w http.ResponseWriter, r *http.Request) {
		records, err := readFile(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}

		// Exchange positon of (x, y) to invert the matrix
		matrixLength := len(records)
		for i := 0; i < matrixLength; i++ {
			for j := i; j < matrixLength; j++ {
				temp := records[i][j]
				records[i][j] = records[j][i]
				records[j][i] = temp
			}
		}

		response := buildeResponse(records)
		fmt.Fprint(w, response)
	})

	/*
		Flatten the input matrix. Example:
			1,2,3
			4,5,6  =>   1,2,3,4,5,6,7,8,9
			7,8,9
	*/
	http.HandleFunc("/flatten", func(w http.ResponseWriter, r *http.Request) {
		records, err := readFile(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}

		var response string
		if len(records) > 0 {
			response = fmt.Sprintf("%s%s", response, strings.Join(records[0], ","))
		}
		for i := 1; i < len(records); i++ {
			response = fmt.Sprintf("%s,%s", response, strings.Join(records[i], ","))
		}
		response = fmt.Sprintf("%s\n", response)
		fmt.Fprint(w, response)
	})

	/*
		Sum each element of the matrix
	*/
	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		records, err := readFile(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}

		var sumResult int
		for _, row := range records {
			for _, element := range row {
				number, err := strconv.Atoi(element)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
					return
				}
				sumResult = sumResult + number
			}
		}
		fmt.Fprint(w, fmt.Sprintf("%d\n", sumResult))
	})

	/*
		Multiple each element of the matrix
	*/
	http.HandleFunc("/multiply", func(w http.ResponseWriter, r *http.Request) {
		records, err := readFile(r)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
			return
		}

		sumResult := 1
		for _, row := range records {
			for _, element := range row {
				number, err := strconv.Atoi(element)
				if err != nil {
					w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
					return
				}
				sumResult = sumResult * number
			}
		}
		fmt.Fprint(w, fmt.Sprintf("%d\n", sumResult))
	})
	http.ListenAndServe(":8080", nil)
}

// Read matrix from given file
func readFile(r *http.Request) ([][]string, error) {
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return csv.NewReader(file).ReadAll()
}

// Convert matrix into a response string
func buildeResponse(matrix [][]string) string {
	var response string
	for _, row := range matrix {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}
	return response
}
