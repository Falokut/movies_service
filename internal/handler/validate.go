package handler

import (
	"errors"
	"regexp"
	"strings"

	movies_service "github.com/Falokut/movies_service/pkg/movies_service/v1/protos"
)

var ErrInvalidFilter = errors.New("invalid filter value, filter must contain only digits and commas")

func validateFilter(filter *movies_service.GetMoviesPreviewRequest) error {
	if filter.GetGenresIDs() != "" {
		if err := checkFilterParam(*filter.GenresIDs); err != nil {
			return err
		}
	}
	if filter.GetCountriesIDs() != "" {
		if err := checkFilterParam(*filter.CountriesIDs); err != nil {
			return err
		}
	}
	if filter.GetMoviesIDs() != "" {
		if err := checkFilterParam(*filter.MoviesIDs); err != nil {
			return err
		}
	}
	if filter.GetAgeRatings() != "" {
		if strings.Contains(*filter.AgeRatings, "'") {
			return ErrInvalidFilter
		}
	}
	return nil
}

func checkFilterParam(val string) error {
	exp := regexp.MustCompile("^[!-&!+,0-9]+$")
	if !exp.Match([]byte(val)) {
		return ErrInvalidFilter
	}

	return nil
}
