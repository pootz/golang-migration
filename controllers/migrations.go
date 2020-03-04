package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/golang-migration/models"
)

var temp = template.Must(template.ParseGlob("templates/*.html"))

// Index ... Carrega homepage
func Index(w http.ResponseWriter, r *http.Request) {
	queries, err := models.GetAllMigrations(0, 0)
	if err != nil {
		fmt.Println(err)
	}

	totalPagesFloat := float64(len(queries) / 15)
	totalPages := uint32(math.Ceil(totalPagesFloat))
	if totalPages <= 0 {
		totalPages = 1
	}

	totalRegisters := 2

	var currentPage uint32 = 1

	transformJSON, _ := json.Marshal(queries)

	var nextPage uint32 = totalPages
	if currentPage < totalPages {
		nextPage = currentPage + 1
	}

	var previousPage uint32 = 1
	if currentPage > 1 {
		previousPage = currentPage - 1
	}

	// tablesList, err := models.ShowAllTables()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// jsonTablesList, _ := json.Marshal(tablesList)

	columnsList, err := models.ShowAllColumns()
	if err != nil {
		fmt.Println(err)
	}
	jsonColumnsList, _ := json.Marshal(columnsList)

	data := map[string]interface{}{
		"Total_registers": totalRegisters,
		"Current_page":    currentPage,
		"Next_page":       nextPage,
		"Previous_page":   previousPage,
		"Queries":         queries,
		// "tables_list":     string(jsonTablesList),
		"Columns_list": string(jsonColumnsList),
		"Json":         string(transformJSON),
	}

	err2 := temp.ExecuteTemplate(w, "Index", data)
	if err2 != nil {
		fmt.Println(err2)
	}
}

// UpdateProduction ... Executa migration no BD de producao
func UpdateProduction(w http.ResponseWriter, r *http.Request) {
	var selectedMigrations []string
	err := json.Unmarshal([]byte(r.FormValue("migrationsToRun")), &selectedMigrations)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "erro ao receber queries à serem executadas")
		return
	}
	if len(selectedMigrations) < 1 {
		fmt.Println("Nenhuma migration selecionada")
		fmt.Fprintf(w, "Nenhuma migration selecionada")
		return
	}
	migrationsToRunIDs := strings.Join(selectedMigrations[:], ",")

	queries, err := models.GetMigrationsByID(2, migrationsToRunIDs)
	if err != nil {
		fmt.Println(err)
	}

	if len(queries) < 1 {
		fmt.Println(err)
		fmt.Fprintf(w, "Todas as migrations já foram rodadas no BD de produção")
		return
	}

	for x := 0; x < len(queries); x++ {
		_, err = models.ExecMigration(queries[x], "producao")

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Erro ao executar a migration")
			return
		}
	}

	fmt.Fprintf(w, "true")
}

// UpdateTest ... Executa migration no BD de teste
func UpdateTest(w http.ResponseWriter, r *http.Request) {
	var selectedMigrations []string
	err := json.Unmarshal([]byte(r.FormValue("migrationsToRun")), &selectedMigrations)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "erro ao receber queries à serem executadas")
		return
	}
	if len(selectedMigrations) < 1 {
		fmt.Println("Nenhuma migration selecionada")
		fmt.Fprintf(w, "Nenhuma migration selecionada")
		return
	}
	migrationsToRunIDs := strings.Join(selectedMigrations[:], ",")

	queries, err := models.GetMigrationsByID(1, migrationsToRunIDs)
	if err != nil {
		fmt.Println(err)
	}

	if len(queries) < 1 {
		fmt.Println(err)
		fmt.Fprintf(w, "Todas as migrations já foram rodadas no BD de teste")
		return
	}

	for x := 0; x < len(queries); x++ {
		_, err = models.ExecMigration(queries[x], "teste")

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Erro ao executar a migration")
			return
		}
	}
	fmt.Fprintf(w, "true")
}

// SaveMigration ... Salva nova migration
func SaveMigration(w http.ResponseWriter, r *http.Request) {
	_, err := models.InsertMigration(r.FormValue("title"), r.FormValue("query"))

	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Erro ao executar insert")
		return
	}

	fmt.Fprintf(w, "true")
}

// DeleteMigration ... Exclui uma migratio que não tenha sido executada nem em teste e nem em produção
func DeleteMigration(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query()["id"][0])
	_, err := models.DeleteMigration(id)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Erro ao excluir migration")
		return
	}

	fmt.Fprintf(w, "true")
}
