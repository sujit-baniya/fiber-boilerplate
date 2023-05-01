package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oarkflow/xid"

	"github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/auth"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/models"
	"github.com/sujit-baniya/fiber-boilerplate/utils/xopen"
)

type SearchFilter struct {
	Field          string
	Search         string
	AdvancedSearch string
}

func FileIndex(c *fiber.Ctx) error {
	user, err := auth.User(c)
	if err != nil {
		return c.Redirect("/")
	}
	files := user.Files
	layout := "layouts/main"
	view := "file-manager"
	if user == nil {
		layout = "layouts/auth"
		view = "landing"
	}

	if err := c.Render(view, fiber.Map{
		"auth":  user != nil,
		"user":  user,
		"files": files,
	}, layout); err != nil {
		// panic(err.Error())
	}
	return nil
}

func Upload(c *fiber.Ctx) error {
	start := time.Now()
	log.SetOutput(ioutil.Discard)
	var err error
	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "documents" key:
		files := form.File["files"]

		user, _ := auth.User(c)
		for _, file := range files {
			err = handleUploadIndividualFile(c, file, user)
		}
	}
	fmt.Printf("\n%2fs", time.Since(start).Seconds())
	return err
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
	fileParts := strings.Split(file.Filename, ".")
	id := xid.New().String()
	id = id + "." + fileParts[1]
	fileName := filepath.Join(app.Http.Server.UploadPath, id)
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
	app.Http.Database.DB.Save(&f)
	uf.FileID = f.ID
	uf.UserID = user.ID
	uf.IsActive = true
	app.Http.Database.DB.Create(&uf)
	return nil
}
