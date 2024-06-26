CREATE DATABASE parque_atracciones;

USE parque_atracciones;

/* Creación de tablas */
CREATE TABLE parque (id varchar(30) PRIMARY KEY, aforoMaximo int, aforoActual int);

CREATE TABLE mapa(fila int(2) PRIMARY KEY, infoParque varchar(100));

CREATE TABLE ciudades(cuadrante varchar(30) PRIMARY KEY, nombre varchar(30), temperatura float);
	
CREATE TABLE visitante (
id varchar(20) PRIMARY KEY, 
nombre varchar(30) NOT NULL, 
contraseña varchar(100) NOT NULL, 
posicionx int DEFAULT 0, 
posiciony int DEFAULT 0,
destinox int DEFAULT -1,
destinoy int DEFAULT -1,
dentroParque int DEFAULT 0,
idEnParque char(1),
ultimoEvento varchar(150),
parqueAtracciones varchar(30) default 'SDpark', 
CONSTRAINT fk_visitantes_parque FOREIGN KEY (parqueAtracciones) REFERENCES parque (id));

CREATE TABLE atraccion(
id varchar(30) PRIMARY KEY, 
tciclo int, 
nvisitantes int, 
posicionx int,
posiciony int,
tiempoEspera int,
estado varchar(10) default 'Abierta',
parqueAtracciones varchar(30),
CONSTRAINT fk_atracciones_parque FOREIGN KEY (parqueAtracciones) REFERENCES parque (id));

/* Inserciones en las tablas */ 
INSERT INTO parque (id, aforoActual) VALUES ('SDpark', 0);
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion1", 5, 3, 10, 14, 45, "SDpark"); 
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion2", 8, 9, 1, 4, 30, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion3", 7, 6, 6, 6, 15, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion4", 18, 9, 10, 19, 65, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion5", 4, 5, 9, 17, 10, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion6", 10, 6, 3, 18, 40, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion7", 11, 8, 9, 2, 80, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion8", 19, 7, 2, 3, 90, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion9", 14, 6, 7, 8, 20, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion10", 15, 4, 18, 11, 13, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion11", 17, 7, 17, 5, 7, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion12", 7, 8, 6, 5, 17, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion13", 8, 2, 16, 17, 77, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion14", 12, 4, 19, 18, 40, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion15", 13, 10, 4, 15, 74, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion16", 19, 11, 15, 15, 23, "SDpark");

INSERT INTO mapa (fila, infoParque) VALUES (0, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (1, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (2, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (3, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (4, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (5, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (6, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (7, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (8, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (9, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (10, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (11, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (12, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (13, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (14, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (15, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (16, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (17, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (18, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");
INSERT INTO mapa (fila, infoParque) VALUES (19, "|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|-|");


delete from visitante;
delete from ciudad;
delete from mapa;

/*Actualizar algunos valores */
UPDATE parque SET aforoMaximo=10 WHERE id = "SDpark";
UPDATE atraccion SET tciclo = 10, nvisitantes = 6, posicionx = 2, posiciony = 10 WHERE id = "atraccion2";
UPDATE atraccion SET tiempoEspera = 5 WHERE id = "atraccion10";
UPDATE atraccion SET estado = "Cerrada" WHERE id = "atraccion10";

select * from atraccion;
select * from visitante;
select ultimoEvento from visitante; /* Para ver los logs de la tabla visitante */
select * from parque;
select * from mapa;
select * from ciudades;

SHOW STATUS LIKE 'max_used_connections';
