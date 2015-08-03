package external

import (
	"crypto/md5"
	"io"
	"log"
	"math"
	"os"
)

// We set the filechunk to 8kb
const filechunk = 8192

// checkAndCreateFolder checks if a folder exists, if not, creates it.
func checkAndCreateFolder(fp string) error {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Printf("Could not find %v folder. Creating it.\n", fp)
		os.Mkdir(fp, 0777)
	} else if err != nil {
		return err
	}
	return nil
}

// generateMd5Sum generates the md5sum of a file.
func generateMd5Sum(fn string) (string, error) {
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
