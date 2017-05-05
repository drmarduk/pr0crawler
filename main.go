package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	cmd := flag.String("cmd", "lsuser", "what should we do")
	user := flag.String("user", "txtinput", "user to crawl")
	flag.Parse()

	switch *cmd {
	case "lsuser":
		do(*user, listUploads)
	case "fav":
		do(*user, listFavorites)
	}
}

func do(user string, fn func(u string, off int) (*ProgrammItemsGet, error)) {
	offset := 0
	for {
		r, err := fn(user, offset)
		if err != nil {
			log.Fatalf("could not get user images: %v\n", err)
		}
		log.Printf("Found %d items\n", len(r.Items))
		if len(r.Items) < 1 {
			log.Println("done")
			return
		}
		downloadCollection(r)
		offset = getNewOffset(r)
	}
}

func getNewOffset(r *ProgrammItemsGet) int {
	min := r.Items[0].ID
	for _, i := range r.Items {
		if min > i.ID {
			min = i.ID
		}
	}
	return min
}
func listUploads(user string, offset int) (*ProgrammItemsGet, error) {
	return genericImages("user", user, offset)
}

func listFavorites(user string, offset int) (*ProgrammItemsGet, error) {
	return genericImages("likes", user, offset)
}

/* Flags:
9: sfw
2: nsfw
4: nsfl

*/

func genericImages(cmd, user string, offset int) (*ProgrammItemsGet, error) {
	url := fmt.Sprintf(
		"http://pr0gramm.com/api/items/get?flags=15&%s=%s",
		cmd,
		user,
	)
	// iterate over all images
	if offset != 0 {
		url = fmt.Sprintf("%s&older=%d", url, offset)
	}

	log.Println("GET JSON:", url)

	result := NewProgrammItemsGet()

	_tmp, err := dl(url)
	if err != nil {
		return nil, err
	}

	log.Printf("JSON: %s\n", _tmp)

	err = json.Unmarshal([]byte(_tmp), result)
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return result, nil
}

func dl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()

	src, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(src), nil
}

func downloadCollection(result *ProgrammItemsGet) error {

	for _, upload := range result.Items {
		url := createDownloadURL(upload)

		fn := strings.Split(url, "/")
		_fn := fn[len(fn)-1]

		log.Println("Download: ", url)
		err := download(url, _fn)
		if err != nil {
			log.Printf("Error: %d: %v\n", upload.ID, err)
		}
	}
	return nil
}

func createDownloadURL(i ProgrammItem) string {
	url := ""
	// get hostname, vid or img
	if strings.HasSuffix(i.Image, "mp4") { // obay on webm's
		url = "http://vid.pr0gramm.com/"
	} else {
		url = "http://img.pr0gramm.com/"
	}

	if i.Fullsize != "" {
		url = "http://full.pr0gramm.com/" + i.Fullsize

	} else {
		url += i.Image
	}

	return url
}

func download(url, file string) error {
	out, err := os.Create(file)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
