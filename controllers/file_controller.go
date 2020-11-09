package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-uuid"
	. "github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/auth"
	"github.com/sujit-baniya/fiber-boilerplate/config"
	"github.com/sujit-baniya/fiber-boilerplate/models"
	"github.com/sujit-baniya/fiber-boilerplate/utils/xopen"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FileIndex(c *fiber.Ctx) error {
	user, _ := auth.User(c)
	layout := "layouts/main"
	view := "file-manager"
	if user == nil {
		layout = "layouts/auth"
		view = "landing"
	}

	return c.Render(view, fiber.Map{
		"auth": user != nil,
		"user": user,
	}, layout)
}

func ViewFile(c *fiber.Ctx) error {
	user, _ := auth.User(c)
	layout := "layouts/main"
	view := "file-view"
	if user == nil {
		layout = "layouts/auth"
		view = "landing"
	}
	return c.Render(view, fiber.Map{
		"auth": user != nil,
		"user": user,
	}, layout)
}

func Upload(c *fiber.Ctx) error {
	start := time.Now()
	log.SetOutput(ioutil.Discard)
	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "documents" key:
		files := form.File["files"]

		user, _ := auth.User(c)
		for _, file := range files {
			handleUploadIndividualFile(c, file, user)
		}
	}
	fmt.Printf("\n%2fs", time.Since(start).Seconds())
	return c.JSON("Uploaded")
}

func LineCounter(r io.Reader) (int64, error) {

	var count int64
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func handleUploadIndividualFile(c *fiber.Ctx, file *multipart.FileHeader, user *models.User) error {
	var f models.File
	var uf models.UserFile
	fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
	fileParts := strings.Split(file.Filename, ".")
	id, _ := uuid.GenerateUUID()
	id = id + "." + fileParts[1]
	fileName := filepath.Join(config.AppConfig.App_Upload_Path, id)
	err := c.SaveFile(file, fileName)
	// Check for errors
	if err != nil {
		return err
		// c.Status(500).Send("Something went wrong 😭")
	}
	fileInfo, _ := os.Stat(fileName)
	fReader, _ := xopen.Ropen(fileName)
	f.Title = file.Filename
	f.Size = fmt.Sprintf("%v", file.Size)
	f.MimeType = file.Header["Content-Type"][0]
	f.FileName = id
	f.Extension = fileParts[1]
	lineCounter, _ := LineCounter(fReader)
	f.RowCount = lineCounter
	f.ModifiedAt = fileInfo.ModTime()
	DB.Save(&f)
	uf.FileID = f.ID
	uf.UserID = user.ID
	uf.IsActive = true
	DB.Save(&uf)
	return c.Next()
}
