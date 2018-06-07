SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

CREATE SCHEMA IF NOT EXISTS `heroku_45a7f4f644be63c` DEFAULT CHARACTER SET latin1 ;
USE `heroku_45a7f4f644be63c` ;

-- -----------------------------------------------------
-- Table `bandsfromtown`.`location`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `heroku_45a7f4f644be63c`.`location` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `city` VARCHAR(150) NOT NULL,
  `state` VARCHAR(150) NOT NULL,
  `country` VARCHAR(150) NOT NULL,
  `full_location` VARCHAR(256) NULL DEFAULT NULL,
  `google_place_id` VARCHAR(256) NULL DEFAULT NULL,
  `longitude` DOUBLE NULL,
  `latitude` DOUBLE NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB
AUTO_INCREMENT = 37
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `bandsfromtown`.`artist`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `heroku_45a7f4f644be63c`.`artist` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(150) NOT NULL,
  `hometown` INT(10) UNSIGNED NOT NULL,
  `genre` VARCHAR(150) NULL DEFAULT NULL,
  `spotify_id` VARCHAR(256) NULL DEFAULT NULL,
  `wikipedia_url` VARCHAR(150) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_location_id` (`hometown` ASC),
  CONSTRAINT `fk_location_id`
    FOREIGN KEY (`hometown`)
    REFERENCES `heroku_45a7f4f644be63c`.`location` (`id`))
ENGINE = InnoDB
AUTO_INCREMENT = 36
DEFAULT CHARACTER SET = latin1;


-- -----------------------------------------------------
-- Table `bandsfromtown`.`album`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `heroku_45a7f4f644be63c`.`album` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(150) NOT NULL,
  `artist_id` INT(10) UNSIGNED NOT NULL,
  `release_date` DATE NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `artist_id` (`artist_id` ASC),
  CONSTRAINT `album_ibfk_1`
    FOREIGN KEY (`artist_id`)
    REFERENCES `heroku_45a7f4f644be63c`.`artist` (`id`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = latin1;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
