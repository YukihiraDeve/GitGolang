package main

import (
	"fmt"
	"log"
	"myapp/api"
	"myapp/server"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Loading .env file...")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur de chargement du fichier .env")

	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Println("Starting server...")
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		userURL := os.Getenv("userURL")

		userRepos := api.GetRepositories(userURL)

		for _, repo := range userRepos {
			fmt.Printf("- %s: %s\n", repo.Name, repo.HTMLURL)
		}

		log.Println("End goroutine Git")

	}()
	wg.Wait()

	log.Println("Zip de Git")
	err := api.CreateArchive("git", "repositories.zip")
	if err != nil {
		log.Fatalf("Erreur lors de la création de l'archive : %v", err)
	}

}
