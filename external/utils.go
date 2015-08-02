package external

import (
	"crypto/md5"
	"io"
	"log"
	"math"
	"net/http"
	"os"
)

// We set the filechunk to 8kb
const filechunk = 8192

// CheckAndCreateFolder checks if a folder exists, if not, creates it.
func CheckAndCreateFolder(fp string) error {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Printf("Could not find %v folder. Creating it.\n", fp)
		os.Mkdir(fp, 0777)
	} else if err != nil {
		return err
	}
	return nil
}

// DownloadNamedFile downloads a given file and writes it on disk
// using a specific filename or path.
func DownloadNamedFile(url, fn string) error {
	out, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}

// GenerateMd5Sum generates the md5sum of a file.
func GenerateMd5Sum(fn string) (string, error) {
	file, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer file.Close()

	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return string(hash.Sum(nil)), nil
}

// SameFileCheck calculates if two files have the same checksum
func SameFileCheck(firstFilename, secondFilename string) (bool, error) {
	firstChecksum, err := GenerateMd5Sum(firstFilename)
	if err != nil {
		return false, err
	}
	secondChecksum, err := GenerateMd5Sum(secondFilename)
	if err != nil {
		return false, err
	}
	return firstChecksum == secondChecksum, nil
}
