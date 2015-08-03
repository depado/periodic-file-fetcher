package external

import (
	"crypto/md5"
	"io"
	"log"
	"math"
	"os"
)

// We set the filechunk to 8kb (used by the md5Sum function)
const filechunk = 8192

// createDirIfNeeded checks if a folder exists, if not, creates it.
func createDirIfNeeded(fp string) error {
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Printf("Creating %v directory.\n", fp)
		os.Mkdir(fp, 0777)
	} else if err != nil {
		return err
	}
	return nil
}

// createDirsIfNeeded creates the dirs given as parameters if they don't exist.
func createDirsIfNeeded(paths ...string) error {
	for _, p := range paths {
		if err := createDirIfNeeded(p); err != nil {
			return err
		}
	}
	return nil
}

// md5Sum generates the md5 sum of a file.
func md5Sum(fn string) (string, error) {
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
