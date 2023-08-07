package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type Course struct {
	C_id     int
	Nama     string
	SKS      int
	Jurusan  string
	Fakultas string
	Semester int
	Prediksi string
	PredC    float64
	PredS    float64
}

type Fakultas struct {
	Buffer map[string]string
}

type CourseList struct {
	CourseBuffer []Course
}

type ProcCourse struct {
	CourseBuffer []Course
	Taken        []int
	IP           float64
	SKS          int
}

func NewFakultas() *Fakultas {
	return &Fakultas{
		Buffer: make(map[string]string),
	}
}

func NewCourseList() *CourseList {
	return &CourseList{
		CourseBuffer: nil,
	}
}
func NewCourse(id int, nama string, sks int, jurusan string, fakultas string, semester int, pred string) *Course {
	return &Course{
		C_id:     id,
		Nama:     nama,
		SKS:      sks,
		Jurusan:  jurusan,
		Fakultas: fakultas,
		Semester: semester,
		Prediksi: pred,
		PredC:    ConvPrediksi(pred),
		PredS:    PredSConv(sks, ConvPrediksi(pred)),
	}
}

func NewProcCourse() *ProcCourse {
	return &ProcCourse{
		CourseBuffer: nil,
		Taken:        nil,
		IP:           0,
		SKS:          0,
	}
}

func (l *CourseList) AddCourse(nama string, sks int, jurusan string, fakultas string, semester int, pred string) {
	newCourse := NewCourse(len(l.CourseBuffer), nama, sks, jurusan, fakultas, semester, pred)

	if newCourse.PredC != -1 {
		l.CourseBuffer = append(l.CourseBuffer, *newCourse)
	}
}

func (l *Fakultas) Add(jurusan string, fakultas string) {
	l.Buffer[jurusan] = fakultas
}

func ConvPrediksi(a string) float64 {
	if a == "A" {
		return 4
	} else if a == "AB" {
		return 3.5
	} else if a == "B" {
		return 3
	} else if a == "BC" {
		return 2.5
	} else if a == "C" {
		return 2
	} else if a == "D" {
		return 1
	} else if a == "E" {
		return 0
	}

	return -1
}

func PredSConv(sks int, predc float64) float64 {
	return float64(sks) * (predc)
}

func GetCourseTotalValue(list []Course) float64 {
	var tempVal float64
	tempVal = 0
	for _, course := range list {
		tempVal = tempVal + course.PredS
	}
	return tempVal
}

func GetListTotalSKS(list []Course) int {
	var temp int
	temp = 0
	for _, course := range list {
		temp = temp + course.SKS
	}
	return temp
}

func (l *CourseList) GetBestValueCourses(jurusan string, fakultas string, sem int, minSks int, maxSks int) [][]Course {
	// Filter courses based on the given criteria
	var filteredCourses []Course
	for _, course := range l.CourseBuffer {
		if course.Jurusan == jurusan && course.Fakultas == fakultas && course.Semester <= sem {
			filteredCourses = append(filteredCourses, course)
		}
	}

	// Find all combinations of courses within SKS constraints
	var bestCourses [][]Course
	findCombinations(filteredCourses, []Course{}, 0, minSks, maxSks, &bestCourses)

	return bestCourses
}

func (l *CourseList) GetBestValueCourses2(jurusan string, fakultas string, sem int, minSks int, maxSks int) [][]Course {
	// Filter courses based on the given criteria
	var filteredCourses []Course
	for _, course := range l.CourseBuffer {
		if course.Fakultas == fakultas && course.Semester <= sem {
			filteredCourses = append(filteredCourses, course)
		}
	}

	// Find all combinations of courses within SKS constraints
	var bestCourses [][]Course
	findCombinations(filteredCourses, []Course{}, 0, minSks, maxSks, &bestCourses)

	return bestCourses
}

func findCombinations(courses []Course, currentCombination []Course, idx int, minSks int, maxSks int, bestCourses *[][]Course) {
	if idx >= len(courses) {
		// Base case: reached the end of courses
		totalSKS := GetListTotalSKS(currentCombination)
		if totalSKS >= minSks && totalSKS <= maxSks {
			currentIP := GetCourseListIP(currentCombination)
			if len(*bestCourses) == 0 || currentIP >= GetCourseListIP((*bestCourses)[0]) {
				*bestCourses = [][]Course{currentCombination}
			} else if currentIP == GetCourseListIP((*bestCourses)[0]) {
				*bestCourses = append(*bestCourses, currentCombination)
			}
		}
		return
	}

	findCombinations(courses, currentCombination, idx+1, minSks, maxSks, bestCourses)

	if currentCombinationSKS := GetListTotalSKS(currentCombination); currentCombinationSKS+courses[idx].SKS <= maxSks {
		// Create a copy of the current combination and add the current course
		newCombination := append([]Course{}, currentCombination...)
		newCombination = append(newCombination, courses[idx])

		// Recurse to the next course
		findCombinations(courses, newCombination, idx+1, minSks, maxSks, bestCourses)
	}
}

func GetCourseListIP(courses []Course) float64 {
	totalPredS := 0.0
	totalSKS := 0

	for _, course := range courses {
		totalPredS += course.PredS
		totalSKS += course.SKS
	}

	if totalSKS == 0 {
		return 0.0
	}

	return totalPredS / float64(totalSKS)
}

func (l *CourseList) readDB(db *sql.DB) {

	rows, err := db.Query("SELECT * FROM courses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var nama string
		var jurusan string
		var fakultas string
		var semester int
		var sks int
		var pred string
		err := rows.Scan(&id, &nama, &jurusan, &fakultas, &semester, &sks, &pred)
		if err != nil {
			log.Fatal(err)
		}
		l.AddCourse(nama, sks, jurusan, fakultas, semester, pred)

	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func (l *Fakultas) readDB(db *sql.DB) {

	rows, err := db.Query("SELECT * FROM fakultas")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var jurusan string
		var fakultas string
		err := rows.Scan(&id, &jurusan, &fakultas)
		if err != nil {
			log.Fatal(err)
		}
		l.Add(jurusan, fakultas)

	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func (l *CourseList) printList() {
	for _, course := range l.CourseBuffer {
		fmt.Println(course)
	}
}

func getAns(c *gin.Context) {
	var input struct {
		Minsks   int    `json:"minsks"`
		Maxsks   int    `json:"maxsks"`
		Jurusan  string `json:"jurusan"`
		Fakultas string `json:"fakultas"`
		Semester int    `json:"semester"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseList := NewCourseList()

	// dsn := "sql6636925:GaydgguNGw@tcp(sql6.freemysqlhosting.net:3306)/sql6636925"
	// db, err := sql.Open("mysql", dsn)
	dsn2 := "db/data.db"
	db, err := sql.Open("sqlite3", dsn2)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	courseList.readDB(db)

	bestCourses := courseList.GetBestValueCourses(
		input.Jurusan, input.Fakultas, input.Semester, input.Minsks, input.Maxsks,
	)

	var outStr string
	for i, courses := range bestCourses {
		// fmt.Printf("Best Courses %d:\n", i+1)
		outStr = outStr + "Best Courses " + strconv.Itoa(i+1) + ": \n"
		for _, course := range courses {
			outStr = outStr + "Mata Kuliah: " + course.Nama + ", SKS: " + strconv.Itoa(course.SKS) + ", Prediksi: " + course.Prediksi + "\n"
			// fmt.Printf("Name: %s, SKS: %d, Prediksi: %s\n", course.Nama, course.SKS, course.Prediksi)
		}
		outStr = outStr + "IP: " + strconv.FormatFloat(GetCourseListIP(courses), 'f', -1, 64) + ", SKS: " + strconv.Itoa(GetListTotalSKS(courses)) + "\n\n"
		// fmt.Printf("Jumlah IP: %.3f, Jumlah SKS: %d\n", GetCourseListIP(courses), GetListTotalSKS(courses))
		// fmt.Println()
	}

	c.JSON(http.StatusOK, gin.H{
		"result": outStr,
	})
}

func getAns2(c *gin.Context) {
	var input struct {
		Minsks   int    `json:"minsks"`
		Maxsks   int    `json:"maxsks"`
		Jurusan  string `json:"jurusan"`
		Semester int    `json:"semester"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	courseList := NewCourseList()
	fakul := NewFakultas()

	// dsn := "sql6636925:GaydgguNGw@tcp(sql6.freemysqlhosting.net:3306)/sql6636925"

	// db, err := sql.Open("mysql", dsn)

	dsn2 := "db/data.db"
	db, err := sql.Open("sqlite3", dsn2)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	courseList.readDB(db)
	fakul.readDB(db)

	fakulIn := fakul.Buffer[input.Jurusan]

	bestCourses := courseList.GetBestValueCourses2(
		input.Jurusan, fakulIn, input.Semester, input.Minsks, input.Maxsks,
	)

	var outStr string
	for i, courses := range bestCourses {
		// fmt.Printf("Best Courses %d:\n", i+1)
		outStr = outStr + "Best Courses " + strconv.Itoa(i+1) + ": \n"
		for _, course := range courses {
			outStr = outStr + "Mata Kuliah: " + course.Nama + ", SKS: " + strconv.Itoa(course.SKS) + ", Prediksi: " + course.Prediksi + "\n"
			// fmt.Printf("Name: %s, SKS: %d, Prediksi: %s\n", course.Nama, course.SKS, course.Prediksi)
		}
		outStr = outStr + "IP: " + strconv.FormatFloat(GetCourseListIP(courses), 'f', -1, 64) + ", SKS: " + strconv.Itoa(GetListTotalSKS(courses)) + "\n\n"
		// fmt.Printf("Jumlah IP: %.3f, Jumlah SKS: %d\n", GetCourseListIP(courses), GetListTotalSKS(courses))
		// fmt.Println()
	}

	c.JSON(http.StatusOK, gin.H{
		"result": outStr,
	})
}

func addMat(c *gin.Context) {
	var input struct {
		Nama     string `json:"nama"`
		Sks      int    `json:"sks"`
		Jurusan  string `json:"jurusan"`
		Fakultas string `json:"fakultas"`
		Semester int    `json:"semester"`
		Prediksi string `json:"prediksi"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// dsn := "sql6636925:GaydgguNGw@tcp(sql6.freemysqlhosting.net:3306)/sql6636925"

	// db, err := sql.Open("mysql", dsn)
	dsn2 := "db/data.db"
	db, err := sql.Open("sqlite3", dsn2)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err2 := db.Exec("INSERT INTO courses (Nama, Jurusan, Fakultas, Semester, SKS, Prediksi) VALUES (?, ?, ?, ?, ?, ?)",
		input.Nama, input.Jurusan, input.Fakultas, input.Semester, input.Sks, input.Prediksi)
	if err != nil {
		log.Fatal(err2)
	}

	outStr := "Inserted: " + input.Nama

	c.JSON(http.StatusOK, gin.H{
		"message": outStr,
	})
}

func addFal(c *gin.Context) {
	var input struct {
		Jurusan  string `json:"jurusan"`
		Fakultas string `json:"fakultas"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// dsn := "sql6636925:GaydgguNGw@tcp(sql6.freemysqlhosting.net:3306)/sql6636925"

	// db, err := sql.Open("mysql", dsn)
	dsn2 := "db/data.db"
	db, err := sql.Open("sqlite3", dsn2)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err2 := db.Exec("INSERT INTO fakultas (Jurusan, Fakultas) VALUES (?, ?)",
		input.Jurusan, input.Fakultas)
	if err != nil {
		log.Fatal(err2)
	}

	outStr := "Inserted: " + input.Jurusan + " => " + input.Fakultas

	c.JSON(http.StatusOK, gin.H{
		"message": outStr,
	})
}

func clearDB(c *gin.Context) {
	// dsn := "sql6636925:GaydgguNGw@tcp(sql6.freemysqlhosting.net:3306)/sql6636925"

	// db, err := sql.Open("mysql", dsn)
	dsn2 := "db/data.db"
	db, err := sql.Open("sqlite3", dsn2)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err2 := db.Exec("DELETE FROM courses")
	if err != nil {
		log.Fatal(err2)
	}

	outStr := "Cleared Database"

	c.JSON(http.StatusOK, gin.H{
		"message": outStr,
	})
}

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/addSingle", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	r.POST("/getAnswer", getAns)
	r.POST("/addMat", addMat)
	r.POST("/clearData", clearDB)
	r.POST("/getAnswer2", getAns2)
	r.POST("/addFakul", addFal)
	r.Run(":8080")
}
