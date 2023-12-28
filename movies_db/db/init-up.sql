CREATE ROLE movies_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE age_ratings (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title_ru TEXT NOT NULL,
    title_en TEXT,
    description TEXT NOT NULL,
    short_description TEXT NOT NULL,
    directors  INT ARRAY,
    genres  INT ARRAY,
    countries  INT ARRAY, 
    duration INT NOT NULL,
    poster_picture_id TEXT,
    background_picture_id TEXT,
    preview_poster_picture_id TEXT,
    age_rating_id INT,
    release_year SMALLINT NOT NULL
);

ALTER TABLE movies ADD CONSTRAINT age_rating_fkey FOREIGN KEY (age_rating_id) REFERENCES age_ratings(id) ON DELETE SET NULL; 

GRANT SELECT ON movies TO movies_service;
GRANT SELECT ON age_ratings TO movies_service;