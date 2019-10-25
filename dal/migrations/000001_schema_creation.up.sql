BEGIN;

CREATE SCHEMA IF NOT EXISTS bands_from_town;

-- -----------------------------------------------------
-- Table 'bands_from_town'.'location'
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS bands_from_town.location (
  id serial PRIMARY KEY,
  city VARCHAR(150) NOT NULL,
  state VARCHAR(150) NOT NULL,
  country VARCHAR(150) NOT NULL,
  full_location VARCHAR(256) NULL DEFAULT NULL,
  google_place_id VARCHAR(256) NULL DEFAULT NULL,
  longitude FLOAT NULL,
  latitude FLOAT NULL);

INSERT INTO bands_from_town.location
VALUES (0, 'Unknown', 'Unknown', 'Unknown', 'Artist location could not be found', -1);
                        
-- -----------------------------------------------------
-- Table 'bandsfromtown'.'artist'
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS bands_from_town.artist (
  id serial PRIMARY KEY,
  name VARCHAR(150) NOT NULL,
  hometown INT NOT NULL REFERENCES bands_from_town.location(id),
  genre VARCHAR(150) NULL DEFAULT NULL,
  spotify_id VARCHAR(256) NULL DEFAULT NULL,
  wikipedia_url VARCHAR(150) NULL DEFAULT NULL
);
 


-- -----------------------------------------------------
-- Table 'bandsfromtown'.'album'
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS bands_from_town.album (
  id serial PRIMARY KEY,
  title VARCHAR(150) NOT NULL,
  artist_id INT NOT NULL REFERENCES bands_from_town.artist(id),
  release_date DATE NOT NULL
);

COMMIT;