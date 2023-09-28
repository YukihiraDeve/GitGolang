package api

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CreateArchive(baseFolder string, archivePath string) error {
	// Créer le fichier zip
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Créer un writer zip pour le fichier zip
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Parcourir les dossiers et les fichiers et les ajouter à l'archive zip
	err = filepath.Walk(baseFolder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Obtenir le chemin relatif à partir du dossier de base
		relPath, err := filepath.Rel(baseFolder, filePath)
		if err != nil {
			return err
		}

		// Créer un header zip pour le fichier
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath

		if info.IsDir() {
			header.Method = zip.Store // Pas de compression pour les dossiers
		} else {
			header.Method = zip.Deflate // Compression pour les fichiers
		}

		// Écrire le header dans le writer zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Si c'est un fichier, écrire le contenu dans le writer zip
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
