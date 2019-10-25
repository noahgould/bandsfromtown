BEGIN;

DROP SCHEMA IF EXISTS bands_from_town CASCADE;

DROP TABLE IF EXISTS bands_from_town.locations CASCADE;

DROP TABLE IF EXISTS bands_from_town.artists CASCADE;

DROP TABLE IF EXISTS bands_from_town.albums CASCADE;

COMMIT;