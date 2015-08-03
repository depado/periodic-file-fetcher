package external

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// A RWMutex Map that contains the available Resources. Must unlock the reader first.
var AvailableResources = struct {
	sync.RWMutex
	m map[string]*Resource
}{m: make(map[string]*Resource)}

// Resource describes an external (distant) resource
// UpdateInterval is a duration describing how often the resource should be fetched
// FriendlyName is the name given to the resource. It is used to create the backup folder,
// and to display friendly logs (without the file name, backup folder name, etc...)
type Resource struct {
	UpdateInterval time.Duration
	FriendlyName   string
	FileName       string
	FullPath       string
	URL            string
	Iterations     int
	Sum            string
	Fetcher        *Fetcher
}

// CalculateIteration parses the backup folder for an ExternalResource and
// determines the highest iteration by reading the file names.
func (res *Resource) calculateIteration(backupFolder string) error {
	files, err := ioutil.ReadDir(backupFolder)
	if err != nil {
		return err
	}
	toplevel := 0
	if len(files) > 0 {
		for _, f := range files {
			level, err := strconv.Atoi(f.Name()[len(f.Name())-1:])
			if err != nil {
				log.Printf("Could not calculate Iteration for %v. It will be set to 0.\n", res.FriendlyName)
				res.Iterations = 0
				return nil
			}
			if level > toplevel {
				toplevel = level
			}
		}
		res.Iterations = toplevel + 1
		return nil
	}
	res.Iterations = toplevel
	return nil
}

func (res *Resource) same(path string) (bool, error) {
	if res.Sum == "" {
		sum, err := md5Sum(res.FullPath)
		if err != nil {
			return false, err
		}
		res.Sum = sum
	}
	sum, err := md5Sum(path)
	if err != nil {
		return false, err
	}
	return res.Sum == sum, nil
}

func (res *Resource) download(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(res.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	return nil
}

// PeriodicUpdate starts the periodic update for the Resource.
func (res *Resource) periodicUpdate() {
	tmpFileName := res.FullPath + ".tmp"
	mapName := res.FileName[0 : len(res.FileName)-len(filepath.Ext(res.FileName))]
	specificBackupFolder := res.Fetcher.BackupDir + mapName + "/"

	if err := createDirsIfNeeded(specificBackupFolder); err != nil {
		log.Println(err)
		return
	}

	if err := res.calculateIteration(specificBackupFolder); err != nil {
		log.Println("Error calculating Iteration :", err)
		return
	}
	if err := res.download(res.FullPath); err != nil {
		log.Println("Error dowloading file :", err)
		return
	}

	tc := time.NewTicker(res.UpdateInterval).C

	log.Println("Resource Collection Started for", res.FriendlyName)
	AvailableResources.Lock()
	AvailableResources.m[mapName] = res
	AvailableResources.Unlock()
	log.Printf("Added Available Resource %v as %v\n", res.FriendlyName, mapName)

	for range tc {
		if err := res.download(tmpFileName); err != nil {
			log.Println(err)
			continue
		}
		same, err := res.same(tmpFileName)
		if err != nil {
			log.Println(err)
			continue
		}
		if same {
			if err := os.Remove(tmpFileName); err != nil {
				log.Println(err)
				continue
			}
		} else {
			if err := os.Rename(res.Fetcher.ConfigurationDir, specificBackupFolder+res.FileName+"."+strconv.Itoa(res.Iterations)); err != nil {
				log.Println(err)
				continue
			}
			res.Iterations++
			if err := os.Rename(tmpFileName, res.FullPath); err != nil {
				log.Println(err)
				continue
			}
			sum, err := md5Sum(res.FullPath)
			if err != nil {
				log.Println(err)
				continue
			}
			res.Sum = sum
			log.Printf("New content in %v. Downloaded and replaced.\n", res.FriendlyName)
		}
	}
}

// LoadConfiguration reads the configuration file and returns an ExternalResource
func loadConfiguration(contentDir, configPath string) (Resource, error) {
	conf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("Could not read external resource configuration :", err)
		return Resource{}, err
	}
	ur := struct {
		UpdateInterval string
		FriendlyName   string
		URL            string
		FileName       string
	}{}
	if err = yaml.Unmarshal(conf, &ur); err != nil {
		log.Println("Error parsing YAML :", err)
		return Resource{}, err
	}
	duration, err := time.ParseDuration(ur.UpdateInterval)
	if err != nil {
		log.Println("Error parsing Duration :", err)
		return Resource{}, err
	}
	external := Resource{
		UpdateInterval: duration,
		FriendlyName:   ur.FriendlyName,
		URL:            ur.URL,
		FileName:       ur.FileName,
		FullPath:       contentDir + ur.FileName,
	}
	return external, nil
}
