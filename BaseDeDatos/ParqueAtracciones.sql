USE parque_atracciones;

/* Creación de tablas */
CREATE TABLE parque (id varchar(30) PRIMARY KEY, aforoMaximo int, aforoActual int);

CREATE TABLE visitante (
id varchar(9) PRIMARY KEY, 
nombre varchar(30) NOT NULL, 
contraseña varchar(30) NOT NULL, 
posicionx int DEFAULT 0, 
posiciony int DEFAULT 0,
destinox int,
destinoy int,
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
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion1", 5, 6, 10, 14, "SDpark"); 
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion2", 9, 7, 1, 4, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion3", 7, 12, 6, 6, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion4", 18, 18, 10, 20, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion5", 4, 10, 9, 17, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion6", 10, 11, 3, 18, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion7", 11, 17, 9, 2, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion8", 20, 13, 2, 3, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion9", 14, 13, 7, 8, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion10", 15, 8, 18, 11, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion11", 17, 14, 17, 5, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion12", 7, 16, 6, 5, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion13", 8, 5, 16, 17, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion14", 12, 9, 20, 18, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion15", 13, 19, 4, 15, "SDpark");
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, parqueAtracciones) VALUES ("atraccion16", 19, 20, 15, 15, "SDpark");