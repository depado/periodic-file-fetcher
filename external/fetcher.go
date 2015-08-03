package external

import (
	"io/ioutil"
	"log"
)

// Fetcher is the main type. It will work with several Resources
type Fetcher struct {
	ConfigurationDir string
	BackupDir        string
	ContentDir       string
}

// Start is the method associated to a Fetcher that will start all the goroutines
// of file downloading and checking.
func (f *Fetcher) Start() error {
	files, err := ioutil.ReadDir(f.ConfigurationDir)
	if err != nil {
		log.Println("Could not read configuration folder :", err)
		return err
	}
	for _, fd := range files {
		log.Println("Loading Configuration :", fd.Name())
		ext, err := loadConfiguration(f.ContentDir, f.ConfigurationDir+fd.Name())
		if err != nil {
			log.Printf("Error with file %v. It won't be used. %v\n", fd.Name(), err)
			continue
		}
		log.Printf("Starting External Resource Collection for %v (%v)\n", ext.FriendlyName, fd.Name())
		go ext.periodicUpdate(f)
	}
	return nil
}
