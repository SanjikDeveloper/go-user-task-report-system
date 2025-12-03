package main

import (
	"html/template"
	"log"
	"os"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type TemplateData struct {
	TodayDate       time.Time
	TasksCount      int
	CompletedCount  int
	InProgressCount int
	Tasks           []Task
}

type Task struct {
	Id          int       `json:"id" db:"id"`
	UserId      int       `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func main() {
	data := TemplateData{
		TodayDate:       time.Now(),
		TasksCount:      10,
		CompletedCount:  10,
		InProgressCount: 0,
	}
	task := Task{
		Id:          1,
		UserId:      2,
		Title:       "test title",
		Description: "test description",
		Date:        time.Now(),
		Status:      "completed",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	data.Tasks = append(data.Tasks, task)

	//templatesPath := "template/html/" /

	// Парсим HTML-шаблон
	tmpl, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		log.Fatal("Ошибка при парсинге шаблона:", err)
	}
	htmlFile, err := os.Create("report.html")
	if err != nil {
		log.Fatal("Ошибка при создании файла:", err)
	}
	defer htmlFile.Close()
	err = tmpl.Execute(htmlFile, data)
	if err != nil {
		log.Fatal("Ошибка при парсинге шаблона:", err)
	}
	// Create new PDF generator
	//pdfOptions := wkhtmltopdf.NewPageOptions()

	// Создаем страницу на основе HTML
	pdfPage := wkhtmltopdf.NewPageReader(htmlFile)

	// Добавляем страницу к PDF
	//pdfOptions.AddPage(pdfPage)

	// Создаем объект PDF и генерируем PDF файл
	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal("1", err)
	}
	pdfGenerator.AddPage(pdfPage)

	err = pdfGenerator.Create()
	if err != nil {
		log.Fatal("Error generating PDF:", err)
	}

	// Сохраняем PDF в файл
	err = pdfGenerator.WriteFile("output.pdf")
	if err != nil {
		log.Fatal("Error saving PDF file:", err)
	}

}
