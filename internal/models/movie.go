package models

type RepositoryMovie struct {
	ID                  int32    `db:"id"`
	TitleRU             string   `db:"title_ru"`
	TitleEN             string   `db:"title_en"`
	Description         string   `db:"description"`
	Genres              []string `db:"genres"`
	Duration            int32    `db:"duration"`
	PosterID            string   `db:"poster_picture_id"`
	BackgroundPictureID string   `db:"background_picture_id"`
	Countries           []string `db:"countries"`
	ReleaseYear         int32    `db:"release_year"`
	AgeRating           string   `db:"age_rating"`
}

type Movie struct {
	ID                   int32
	TitleRU              string
	TitleEN              string
	Description          string
	Genres               []string
	Duration             int32
	PosterUrl            string
	BackgroundUrl string
	Countries            []string
	ReleaseYear          int32
	AgeRating            string
}
