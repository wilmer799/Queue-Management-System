USE parque_atracciones;

/* Creación de tablas */
CREATE TABLE parque (id varchar(30) PRIMARY KEY, aforoMaximo int, aforoActual int);

CREATE TABLE visitante (
id varchar(20) PRIMARY KEY, 
nombre varchar(30) NOT NULL, 
contraseña varchar(30) NOT NULL, 
posicionx int DEFAULT 0, 
posiciony int DEFAULT 0,
destinox int DEFAULT -1,
destinoy int DEFAULT -1,
parqueAtracciones varchar(30), 
CONSTRAINT fk_visitantes_parque FOREIGN KEY (parqueAtracciones) REFERENCES parque (id));

CREATE TABLE atraccion(
id varchar(30) PRIMARY KEY, 
tciclo int, 
nvisitantes int, 
posicionx int,
posiciony int,
tiempoEspera int,
parqueAtracciones varchar(30),
CONSTRAINT fk_atracciones_parque FOREIGN KEY (parqueAtracciones) REFERENCES parque (id));

/* Inserciones en las tablas */ 
INSERT INTO parque (id, aforoActual) VALUES ('SDpark', 0);
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion1", 5, 3, 10, 14, "SDpark"); 
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion2", 9, 4, 1, 4, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion3", 7, 6, 6, 6, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion4", 18, 9, 10, 20, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion5", 4, 5, 9, 17, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion6", 10, 6, 3, 18, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion7", 11, 8, 9, 2, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion8", 20, 7, 2, 3, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion9", 14, 6, 7, 8, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion10", 15, 4, 18, 11, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion11", 17, 7, 17, 5, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion12", 7, 8, 6, 5, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion13", 8, 2, 16, 17, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion14", 12, 4, 20, 18, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion15", 13, 10, 4, 15, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion16", 19, 11, 15, 15, "SDpark");
/* Visitantes */
INSERT INTO visitante (id, nombre, contraseña, posicionx, posiciony, destinox, destinoy,parqueAtracciones)
VALUES ("wilmer88", "wilmer","tubaby",1,1,10,14,"SDpark");
INSERT INTO visitante (id, nombre, contraseña, posicionx, posiciony, destinox, destinoy,parqueAtracciones)
VALUES ("elbala00", "Valentin","catar2022",5,4,10,14,"SDpark");
/*Actualizar algunos valores */
UPDATE parque SET aforoMaximo=10 WHERE id = "SDpark";
UPDATE atraccion SET tiempoEspera = 7 WHERE id = "atraccion1";
UPDATE atraccion SET tiempoEspera = 5 WHERE id = "atraccion10";
