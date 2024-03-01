package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	imgExtension             string = ".png"
	singleImageAPI           string = "/image"
	multiImageAPI            string = "/images"
	dirCarsPath              string = "cars/"
	port                     string = ":8080"
	fileNameSegmentSeparator string = "_"
	fileNameUnknownSegment   string = "*"
	paramModel               string = "model"
	paramMaker               string = "maker"
	paramShape               string = "shape"
	paramYears               string = "years"
	params                          = []string{paramMaker, paramModel, paramShape, paramYears}
)

// var params = [...]string {paramMaker, paramModel, paramShape, paramYears}

func main() {

	singleImageHandler := http.HandlerFunc(getImage)
	// multiImageHandler := http.HandlerFunc(getAllImagesForGen)

	http.Handle(singleImageAPI, singleImageHandler)
	// http.Handle(multiImageAPI, multiImageHandler)

	fmt.Printf("Server started at port %s", port)
	http.ListenAndServe(port, nil)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	makerNameValue := r.URL.Query().Get(paramMaker)
	isRegex := false
	// req_url_with_params, _ := url.Parse(r.URL.String())
	// params, _ := url.ParseQuery(req_url_with_params.RawQuery)
	file_name_builder := strings.Builder{}

	for _, paramKey := range params {
		paramVal := r.URL.Query().Get(paramKey)
		if paramVal != "" {
			file_name_builder.WriteString(fileNameSegmentSeparator)
			file_name_builder.WriteString(paramVal)
		}
	}

	if r.URL.Query().Get(paramModel) != "" && r.URL.Query().Get(paramShape) == "" && r.URL.Query().Get(paramYears) == "" {
		isRegex = true
		file_name_builder.WriteString(fileNameSegmentSeparator)
		file_name_builder.WriteString(fileNameUnknownSegment)
	}
	file_name_builder.WriteString(imgExtension)

	file_name_regex := file_name_builder.String()[1:]
	if file_name_regex == "" {
		file_name_regex = "Mazda"
	}

	if isRegex {
		file_name_regex = findLatestImageByNameRegex(file_name_regex, makerNameValue)
	}

	if !strings.Contains(file_name_regex, imgExtension) {
		file_name_regex = file_name_regex + imgExtension
	}
	relative_path := fmt.Sprintf("%s%s/%s", dirCarsPath, makerNameValue, file_name_regex)

	buf, err := ioutil.ReadFile(relative_path)
	if err != nil {
		log.Println(err)
	}

	// w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Type", "image/png")
	w.Write(buf)
}

func findLatestImageByNameRegex(pattern string, makerNameValue string) string {

	files, err := ioutil.ReadDir(fmt.Sprintf("%s%s", dirCarsPath, makerNameValue))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]

		match, err := filepath.Match(pattern, file.Name())
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}
		if match {
			return file.Name()
		}
	}
	return ""
}

// func getAllImagesForGen(w http.ResponseWriter, r *http.Request) {

// 	// makerNameValue := r.URL.Query().Get(paramMaker)
// 	// modelNameValue := r.URL.Query().Get(paramModel)
// 	// shapeValue := r.URL.Query().Get(paramShape)
// 	// yearsValue := r.URL.Query().Get(paramYears)

// 	// pattern := fmt.Sprintf("%s_%s_%s_%s")

// 	// files, err := ioutil.ReadDir(fmt.Sprintf("%s%s", dirCarsPath, makerNameValue))
// 	// if err != nil {
// 	// 	fmt.Println("Error:", err)
// 	// 	return
// 	// }
// 	// var matchingFiles []fs.FileInfo
// 	// for _, file := range files {
// 	// 	match, err := filepath.Match(pattern, file.Name())
// 	// 	if err != nil {
// 	// 		fmt.Println("Error:", err)
// 	// 		return
// 	// 	}
// 	// 	if match {
// 	// 		matchingFiles = append(matchingFiles, file)
// 	// 	}
// 	// }
// 	// fmt.Println("Files:", matchingFiles)
// }
