package models

type RepositoryMoviePreview struct {
	ID               int32    `db:"id"`
	TitleRU          string   `db:"title_ru"`
	TitleEN          string   `db:"title_en"`
	Duration         int32    `db:"duration"`
	PreviewPosterID  string   `db:"preview_poster_picture_id"`
	Genres           []string `db:"genres"`
	ShortDescription string   `db:"short_description"`
	Countries        []string `db:"countries"`
	ReleaseYear      int32    `db:"release_year"`
	AgeRating        string   `db:"age_rating"`
}

type MoviePreview struct {
	ID               int32
	TitleRU          string
	TitleEN          string
	Duration         int32
	PreviewPosterUrl string
	Genres           []string
	ShortDescription string
	Countries        []string
	ReleaseYear      int32
	AgeRating        string
}
