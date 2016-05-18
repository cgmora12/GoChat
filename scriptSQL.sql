DROP SCHEMA IF EXISTS `gochat` ;

create database gochat;

CREATE TABLE `gochat`.`usuarios` (
  `nombreUsuario` VARCHAR(80) NOT NULL,
  `password` VARCHAR(100) NOT NULL,
  `nombreCompleto` VARCHAR(100) NULL,
  `pais` VARCHAR(80) NULL,
  `provincia` VARCHAR(80) NULL,
  `localidad` VARCHAR(80) NULL,
  `email` VARCHAR(80) NULL,
  PRIMARY KEY (`nombreUsuario`));
