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
	wg.Add(2)

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
	}()

	// Attendre que les deux goroutines soient terminées
	wg.Wait()

}
