package carpenter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/desmondhume/furrytemple/job"
	_ "github.com/lib/pq"
	"google.golang.org/api/youtube/v3"
	// "io/ioutil"
	// "net/http"
	// "strconv"
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

func Exec(command string, input chan map[string]interface{}, output job.Reports) {
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
		//  fmt.Println(err)
		//  fmt.Println("Exiting...")
		// }()

		output <- job.JobReport{"OK", "Everything worked."}
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
