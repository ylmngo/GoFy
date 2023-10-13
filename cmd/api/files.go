package main

import (
	"fmt"
	"gofy/internal/data"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type envelope map[string]interface{}

func (app *application) displayFileHandler(writer http.ResponseWriter, request *http.Request) {
	user := app.getContextUser(request)
	if user.IsAnonymous() {
		query, err := url.ParseQuery(request.URL.RawQuery)
		if err != nil {
			_ = app.writeJSON(writer, envelope{"response": "You are not logged in!"})
			return
		}
		token := query["access_token"][0]
		user, err = app.model.Users.GetByToken(data.ScopeAuthentication, token)
		if err != nil {
			app.logger.Println("Error while retreiving user from url query")
			return
		}
	}

	var input struct {
		Filename string `json:"filename"`
	}

	filename, ok := mux.Vars(request)["filename"]
	if !ok {
		app.logger.Println("No Filename Provided")
		return
	}

	input.Filename = filename

	file, err := app.model.Files.GetFile(user.ID, input.Filename)
	if err != nil {
		app.logger.Println("Error while retreiving file info from database")
		return
	}

	storedFile, err := uuid.FromBytes(file.FileID)
	if err != nil {
		app.logger.Println("Cannot convert fileId to UUID")
		return
	}

	storedFilePath := "uploads/" + storedFile.String()
	storedFileContents, err := os.ReadFile(storedFilePath)
	if err != nil {
		app.logger.Println("Error while reading the stored File")
		return
	}

	writer.Header().Set("Content-Type", "application/pdf")
	writer.Write(storedFileContents)
}

func (app *application) displayFilesHandler(writer http.ResponseWriter, request *http.Request) {

	type Response struct {
		Filename string `json:"filename"`
		FileURL  string `json:"file_url"`
		Metadata string `json:"metadata"`
	}

	user := app.getContextUser(request)
	if user.IsAnonymous() {
		_ = app.writeJSON(writer, envelope{"response": "You are not logged in!"})
		return
	}

	files, err := app.model.Files.GetFiles(user.ID)
	if err != nil {
		app.logger.Println("Error while retreving files info from database")
		return
	}

	authorizationHeader := request.Header.Get("Authorization")
	token := strings.Split(authorizationHeader, " ")[1]

	var jsonResponse []Response
	for _, file := range files {
		jr := Response{
			Filename: file.Filename,
			FileURL:  fmt.Sprintf("http://localhost:8000/file/%s?access_token=%s", file.Filename, token),
			Metadata: file.Metadata,
		}

		jsonResponse = append(jsonResponse, jr)
	}

	if err := app.writeJSON(writer, jsonResponse); err != nil {
		app.logger.Println("Error while writing to JSON")
		return
	}
}

func (app *application) uploadFileHandler(writer http.ResponseWriter, request *http.Request) {
	// Get the user
	user := app.getContextUser(request)
	if user.IsAnonymous() {
		_ = app.writeJSON(writer, envelope{"response": "You are not logged in!"})
		return
	}

	// Parse Multipart form and get the file
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		app.logger.Println("Could not parse multipart form")
		return
	}
	file, header, err := request.FormFile("file")
	defer file.Close()
	if err != nil {
		app.logger.Println("Could not Get File")
		return
	}

	// Get the filename, metadata and content from the multipart form
	filename := header.Filename              // Get the filename
	metadata := request.Form.Get("metadata") // Get the metadata
	fileID := uuid.New()                     // Get the File ID

	// Create the File struct
	dataFile := &data.File{
		Filename:   filename,
		FileID:     fileID[:],
		Metadata:   metadata,
		UserID:     user.ID,
		UploadedAt: time.Now(),
	}

	// Insert filename, metadata, fileID, userID to Database
	err = app.model.Files.Insert(dataFile, dataFile.UserID)
	if err != nil {
		app.logger.Println("Error while Inserting File Info to Database")
		return
	}

	// Store the file in the Local File System
	storageFileName := fileID.String()
	storedFile, err := os.OpenFile("uploads/"+storageFileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		app.logger.Println("Error while storing file to disk")
		return
	}
	defer storedFile.Close()

	// Copy the contents of Form File to storedFile
	if _, err := io.Copy(storedFile, file); err != nil {
		app.logger.Fatal("Error while Writing to Stored File")
		return
	}
}
