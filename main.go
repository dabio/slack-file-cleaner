package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Token represents all files a single user has access to.
type Token struct {
	Name, Token string
}

func (t *Token) getFiles(until time.Time) ([]fileResponse, error) {
	api := fmt.Sprintf(listFilesURL, t.Token, until.Unix())

	resp := listFilesResponse{}
	if err := getAPIResponse(api, &resp); err != nil {
		return nil, err
	}

	return resp.Files, nil
}

// DeleteFiles removes all files a token has the rights for. An optional Time
// object can be specified that defines the time until the files are removed.
func (t *Token) DeleteFiles(until time.Time, force bool) error {
	// get files of token
	files, err := t.getFiles(until)
	if err != nil {
		return err
	}

	log.Printf("Deleting %d files for %v", len(files), t.Name)
	if force {
		for _, file := range files {
			if err = t.DeleteFile(file.ID); err != nil {
				log.Printf("Cannot not delete file %v (%v)", file.ID, err.Error())
			}
		}
	}

	return nil
}

// DeleteFile removes a remote file, identified through the fileID.
func (t *Token) DeleteFile(fileID string) error {
	api := fmt.Sprintf(deleteFileURL, t.Token, fileID)
	res := deleteFileResponse{}
	if err := getAPIResponse(api, &res); err != nil {
		return err
	}

	if !res.Ok {
		return errors.New(res.Error)
	}

	return nil
}

func getAPIResponse(api string, v interface{}) error {
	resp, err := http.Get(api)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&v); err != nil {
		return err
	}

	return nil
}

func main() {
	force := flag.Bool("force", false, "set when you really want to delete files")
	days := flag.Int("days", 30, "deletes files older than n-days")
	config := flag.String("config", "./config.json", "path to config file")
	flag.Parse()

	file, err := os.Open(*config)
	if err != nil {
		log.Panic("could not open file ", config)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	tokens := []Token{}
	if err = decoder.Decode(&tokens); err != nil {
		log.Panicf("cannot parse file %v", config)
	}

	from := time.Now().AddDate(0, 0, -*days)
	log.Printf("Deleting files for %v users older than %d days (before %s)",
		len(tokens),
		*days,
		from.Format(time.RFC822))

	for _, token := range tokens {
		token.DeleteFiles(from, *force)
	}
}
