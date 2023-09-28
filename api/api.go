package api

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	HTMLURL     string    `json:"html_url"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Repositories []Repository

func GetRepositories(url string) Repositories {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var repositories Repositories
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(repositories, func(i, j int) bool {
		return repositories[i].UpdatedAt.After(repositories[j].UpdatedAt)
	})
	WriteRepositoriesToCSV(repositories, "repositories.csv")

	return repositories
}

func WriteRepositoriesToCSV(repos Repositories, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Impossible de créer le fichier CSV :", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Écrire l'en-tête du CSV
	err = writer.Write([]string{"Name", "Description", "URL"})
	if err != nil {
		log.Fatal("Erreur lors de l'écriture dans le fichier CSV :", err)
	}

	// Écrire les données des repositories
	for _, repo := range repos {
		err := writer.Write([]string{repo.Name, repo.Description, repo.HTMLURL})
		if err != nil {
			log.Fatal("Erreur lors de l'écriture dans le fichier CSV :", err)
		}
	}
}
