package api

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		WriteRepositoriesToCSV(repositories, "repositories.csv")
	}()
	go func() {
		defer wg.Done()
		cloneRepositories(repositories)
	}()

	wg.Wait()

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

	err = writer.Write([]string{"Name", "Description", "URL"})
	if err != nil {
		log.Fatal("Erreur lors de l'écriture dans le fichier CSV :", err)
	}

	for _, repo := range repos {
		err := writer.Write([]string{repo.Name, repo.Description, repo.HTMLURL})
		if err != nil {
			log.Fatal("Erreur lors de l'écriture dans le fichier CSV :", err)
		}
	}
	log.Println("CSV Whrite terminé")
}

func cloneRepositories(repos Repositories) {
	for _, repo := range repos {
		cloneDir := filepath.Join("git")

		// Création du dossier de clonage s'il n'existe pas
		err := os.MkdirAll(cloneDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Impossible de créer le dossier de clonage %s : %v", cloneDir, err)
		}

		repoPath := filepath.Join(cloneDir, repo.Name)
		log.Println("git clone ", repo.Name)
		cmd := exec.Command("git", "clone", repo.HTMLURL, filepath.Join(cloneDir, repo.Name))
		err = cmd.Run()

		if err != nil {
			log.Printf("Erreur lors du clonage du dépôt %s : %v", repo.Name, err)
		}

		log.Println("Git clone terminé" + repoPath)

		pullLatestBranch(repoPath)
	}
}

func pullLatestBranch(repoPath string) {
	log.Println("Git pull sur la derniere branch", repoPath)

	// Exécutez git fetch pour récupérer toutes les références de branches en local
	cmd := exec.Command("git", "-C", repoPath, "fetch")
	err := cmd.Run()
	if err != nil {
		log.Printf("Erreur lors de l'exécution de git fetch dans le dépôt %s : %v", repoPath, err)
		return
	}

	// Trouvez la branche avec le dernier commit
	cmd = exec.Command("git", "-C", repoPath, "for-each-ref", "--sort=-committerdate", "--count=1", "--format=%(refname:short)", "refs/heads/")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Erreur lors de la recherche de la dernière branche du dépôt %s : %v", repoPath, err)
		return
	}
	branchName := strings.TrimSpace(string(output))

	cmd = exec.Command("git", "-C", repoPath, "pull", "origin", branchName)
	err = cmd.Run()
	if err != nil {
		log.Printf("Erreur lors de l'exécution de git pull sur la branche %s du dépôt %s : %v", branchName, repoPath, err)
	}
	log.Println("Git pull terminé" + repoPath)
}
