create database gochat;

CREATE TABLE `gochat`.`usuarios` (
  `nombreUsuario` VARCHAR(45) NOT NULL,
  `password` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`nombreUsuario`));
