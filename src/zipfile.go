package main

import (
	"archive/zip"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/akutz/sortfold"
	"golang.org/x/image/draw"
)

type ZipFile struct {
	zipFile *zip.ReadCloser
}

func (zf *ZipFile) Open() {
	zr, err := zip.OpenReader(filepath.Join(ConfigFolder, "smcontent.zip"))
	zf.zipFile = zr

	if err != nil {
		log.Print(err)
	}
}

func (zf *ZipFile) Close() {
	err := zf.zipFile.Close()
	if err != nil {
		log.Print(err)
	}
}

// Find a file in smcontent.zip
func (zf *ZipFile) FindZipFile(filePath string) *zip.File {
	zf.Open()

	for _, z := range zf.zipFile.File {
		if z.Name == filePath {
			return z
		}
	}
	return nil
}

// Find files in smcontent.zip, ignoring image files
func (zf *ZipFile) FindZipFiles(filePath string) []*zip.File {
	var zi []*zip.File
	zf.Open()

	for _, z := range zf.zipFile.File {
		if strings.HasPrefix(z.Name, filePath) && !strings.HasSuffix(z.Name, ".jpg") && !strings.HasSuffix(z.Name, ".jpeg") && !strings.HasSuffix(z.Name, ".png") {
			zi = append(zi, z)
		}
	}

	return zi
}

// Extract a zip item to a given folder
func (zf *ZipFile) ExtractFileToFolder(f *zip.File, filePath string) {
	if f == nil {
		return
	}
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			log.Print(err)
		}
		return
	}

	mkdir := filepath.Join(filePath, filepath.Dir(f.Name))
	if err := os.MkdirAll(mkdir, os.ModePerm); err != nil {
		log.Print(err)
	}

	destinationFile, err := os.OpenFile(filepath.Join(filePath, f.Name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		log.Print(err)
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		log.Print(err)
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		log.Print(err)
	}
}

func (zf *ZipFile) ExtractFilesToFolder(zi []*zip.File, filePath string) {
	for _, z := range zi {
		zf.ExtractFileToFolder(z, filePath)
	}
}

// Creates and updates smcontent.zip
func (zf *ZipFile) UpdateZipfile(filesToZip map[string]string) {
	zipfilename := filepath.Join(ConfigFolder, "smcontent.zip")
	if _, err := os.Stat(zipfilename); errors.Is(err, os.ErrNotExist) {
		newfile, err := os.Create(zipfilename)
		if err != nil {
			log.Print(err)
		}
		defer newfile.Close()

		w := zip.NewWriter(newfile)

		f, err := w.Create("readme.txt")
		if err != nil {
			log.Print(err)
		}
		_, err = f.Write([]byte("This archive is maintained by servermanager."))
		if err != nil {
			log.Print(err)
		}

		err = w.Close()
		if err != nil {
			log.Print(err)
		}
	}

	zr, err := zip.OpenReader(zipfilename)
	if err != nil {
		log.Print(err)
	}
	defer zr.Close()
	zwf, err := os.Create(zipfilename + "_")
	defer zwf.Close()
	zw := zip.NewWriter(zwf)
	defer zwf.Close()

	log.Print("Compressing files...")

	defer zw.Close()
	var wg sync.WaitGroup

	keys := make([]string, 0, len(filesToZip))
	for k := range filesToZip {
		keys = append(keys, k)
	}

	// sort keys by name, insensitively
	sort.Slice(keys, func(i, j int) bool {
		return sortfold.CompareFold(keys[i], keys[j]) < 0
	})

	newFiles := make([]string, 0)
	for _, filepath := range keys {
		destination := filesToZip[filepath]
		wg.Add(1)
		func() {
			defer wg.Done()

			src, err := os.Open(filepath)
			if err != nil {
				log.Print(err)
			}

			fi, err := src.Stat()
			if err != nil {
				log.Print(err)
			}

			// Skip file if it already exists and has the same timestamp
			for _, zipItem := range zr.File {
				if zipItem.Name == destination && zipItem.Modified.Unix() == fi.ModTime().Unix() {
					return
				}
			}

			fih := &zip.FileHeader{
				Name:     destination,
				Method:   zip.Deflate,
				Modified: fi.ModTime(),
			}
			if err != nil {
				log.Print(err)
			}

			dest, err := zw.CreateHeader(fih)
			if err != nil {
				log.Print(err)
			}

			defer src.Close()
			if strings.HasSuffix(filepath, ".jpg") || strings.HasSuffix(filepath, ".jpeg") {
				img, err := jpeg.Decode(src)

				if err != nil {
					if _, err := io.Copy(dest, src); err != nil {
						log.Print(err)
					}
					return
				}

				width := 640
				height := int(float32(img.Bounds().Max.Y) / float32(img.Bounds().Max.X) * float32(width))
				if img.Bounds().Max.X < width {
					width = img.Bounds().Max.X
					height = img.Bounds().Max.Y
				}

				dst := image.NewRGBA(image.Rect(0, 0, width, height))
				draw.BiLinear.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

				jpeg.Encode(dest, dst, &jpeg.Options{
					Quality: jpeg.DefaultQuality,
				})
			} else if strings.HasSuffix(filepath, ".png") {
				img, err := png.Decode(src)

				if err != nil {
					if _, err := io.Copy(dest, src); err != nil {
						log.Print(err)
					}
					return
				}

				width := 640
				height := int(float32(img.Bounds().Max.Y) / float32(img.Bounds().Max.X) * float32(width))
				if img.Bounds().Max.X < width {
					width = img.Bounds().Max.X
					height = img.Bounds().Max.Y
				}

				dst := image.NewRGBA(image.Rect(0, 0, width, height))
				draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

				png.Encode(dest, dst)
			} else {
				if _, err := io.Copy(dest, src); err != nil {
					log.Print(err)
				}
			}

			newFiles = append(newFiles, destination)
		}()
	}

	wg.Wait()

	log.Print("Copying zipfile content...")
	inNewFiles := func(value string) bool {
		for _, item := range newFiles {
			if strings.ToLower(item) == strings.ToLower(value) {
				return true
			}
		}
		return false
	}

	for _, zipItem := range zr.File {
		if inNewFiles(zipItem.Name) {
			continue
		}

		zipItemReader, err := zipItem.OpenRaw()
		if err != nil {
			log.Print(err)
		}

		header := zipItem.FileHeader
		targetItem, err := zw.CreateRaw(&header)
		_, err = io.Copy(targetItem, zipItemReader)
		if err != nil {
			log.Print(err)
		}
	}

	os.Remove(zipfilename)
	os.Rename(zipfilename+"_", zipfilename)
}
