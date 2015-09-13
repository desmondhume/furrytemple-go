package carpenter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/api/youtube/v3"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	db *sql.DB
)

const (
	DB_USER = "furrytemple"
	DB_NAME = "furrytemple_development"

	YOUTUBE_BASE_URL = "https://www.youtube.com/watch?v="
	GA_CODE          = "UA-66841900-1"
)

func insertVideo(video *youtube.SearchResult) error {
	var lastInsertedYoutubeId string

	err := db.QueryRow(`
		INSERT INTO videos(youtube_id, url, thumbnail_url, title)
		SELECT $1, $2, $3, $4
		WHERE NOT EXISTS (
			SELECT id FROM videos
			WHERE youtube_id = $5
		) returning youtube_id;`,
		video.Id.VideoId, fmt.Sprintf("%s%s", YOUTUBE_BASE_URL, video.Id.VideoId), video.Snippet.Thumbnails.Medium.Url, video.Snippet.Title, video.Id.VideoId,
	).Scan(&lastInsertedYoutubeId)

	fmt.Println(lastInsertedYoutubeId)

	return err
}

func populate(data map[string]interface{}) error {
	var videos youtube.SearchListResponse

	videosMarshaled, _ := json.Marshal(data)
	json.Unmarshal(videosMarshaled, &videos)

	for _, item := range videos.Items {
		switch item.Id.Kind {
		case "youtube#video":
			err := insertVideo(item)
			if err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

func Export() error {
	var err error
	var page int
	pages := 1
	limit := 30

	for page < pages {
		var filename string
		var basepath = "./"

		if page == 0 {
			filename = "index"
		} else {
			filename = strconv.Itoa(page)
		}

		filepath := fmt.Sprintf("%s%s.html", basepath, filename)
		file, err := os.Create(filepath)

		if err != nil {
			return err
		}

		stmt, _ := db.Prepare(fmt.Sprintf("SELECT url, thumbnail_url, youtube_id, title FROM videos ORDER BY random() LIMIT %d OFFSET %d", limit, page*limit))
		rows, _ := stmt.Query()

		head := `<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Furry Temple</title>
				<link rel="stylesheet" href="style.css" />
			</head>
			<body>
				<header>
					<div class="banner">Banner here</div>
					<div class="logo">FURRY TEMPLE</div>
					<div class="banner">Banner here</div>
				</header>
				<div class="tagline">100 cat videos/hour</div>
			`
		ga := fmt.Sprintf(`<script>
			  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
			  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
			  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
			  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

			  ga('create', '%s', 'auto');
			  ga('send', 'pageview');

			</script>`, GA_CODE)

		_, err = file.WriteString(head)
		_, err = file.WriteString(ga)

		for rows.Next() {
			var url string
			var thumbnail_url string
			var youtube_id string
			var title string

			rows.Scan(&url, &thumbnail_url, &youtube_id, &title)

			imageFilePath := fmt.Sprintf("%simages/%s.jpg", basepath, youtube_id)
			resp, err := http.Get(thumbnail_url)
			content, err := ioutil.ReadAll(resp.Body)

			imageFile, err := os.Create(imageFilePath)
			imageFile.Write(content)
			imageFile.Close()
			defer resp.Body.Close()
			_, err = file.WriteString(fmt.Sprintf("<div><a href='%s' target=\"_blank\"><img src=\"%s\" /><span>%s</span></a></div>", url, fmt.Sprintf("images/%s.jpg", youtube_id), title))
			if err != nil {
				return err
			}
		}

		_, err = file.WriteString("</body> </html>")

		page++
	}

	return err
}

func Exec(command string, input chan map[string]interface{}, output chan map[string]string) {
	var err error

	for data := range input {
		switch command {
		case "populate":
			err = populate(data)
		default:
			fmt.Println("Command not found")
		}

		if err != nil {
			fmt.Println(err)
			// panic(err)
		}

		// defer func() {
		// 	fmt.Println(err)
		// 	fmt.Println("Exiting...")
		// }()

		job := map[string]string{"status": "OK", "message": "Everything worked."}
		output <- job
	}
}

func init() {
	var err error

	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		DB_USER, DB_NAME)
	db, err = sql.Open("postgres", dbinfo)

	if err != nil {
		fmt.Println("Can't connect to db with dbinfo : ", dbinfo)
	}
}
