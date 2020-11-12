# Golang
Assignment on Golang

Pre Requisite

1. Latest golang installed
2. Rabbit mq installed on the machine
3. Mysql installed on the machine and a new user with username :go and password go
4. Create Schem and table in the mysql

CREATE SCHEMA `mauqah` ;

CREATE TABLE `mauqah`.`person` (
  `id`  INT NOT NULL AUTO_INCREMENT ,
  `name` VARCHAR(45) NOT NULL,
  `age` INT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);


To run 

1. Go to Root directory after cloning , and run go build 
2. Go to publisher  run command go run publisher.go to start the publisher
3. Go to consumer  run command go run consumer.go to start the consumer



Sample inputs

{"Name":"a","Age":2,"Operation":"CREATE"}
{"Name":"b","Age":34,"Operation":"Hi"}
{"Name":"c","Age":33,"Operation":"CREATE"}
{"name":"c","age":29,"operation":"UPDATE","recordId":6}
{"operation":"delete","recordId":5}
{"Name":"c","Operation":"GET"}
{"Name":"","Operation":"GET"}
