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
dentroParque int DEFAULT 0,
idParque char(1),
parqueAtracciones varchar(30) default 'SDpark', 
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
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion1", 5, 3, 10, 14, 45, "SDpark"); 
INSERT INTO atraccion (id, tciclo, nvisitantes, posicionx, posiciony, tiempoEspera, parqueAtracciones) VALUES ("atraccion2", 9, 4, 1, 4, 30, "SDpark");
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
/* Visitantes */
INSERT INTO visitante (id, nombre, contraseña, posicionx, posiciony, destinox, destinoy, dentroParque, idParque, parqueAtracciones)
VALUES ("wilmer88", "wilmer","tubaby",1,1,10,14,0,"w","SDpark");
INSERT INTO visitante (id, nombre, contraseña, posicionx, posiciony, destinox, destinoy, dentroParque, idParque, parqueAtracciones)
VALUES ("elbala00", "Valentin","catar2022",5,4,10,14,0,"e","SDpark");

INSERT INTO visitante (id, nombre, contraseña, posicionx, posiciony, destinox, destinoy, dentroParque, idParque, parqueAtracciones)
VALUES ("rafajaja", "rafa","1234",13,7,17,9,1,"r","SDpark");
INSERT INTO visitante (id, nombre, contraseña, posicionx, posiciony, destinox, destinoy, dentroParque, idParque, parqueAtracciones)
VALUES ("hcarlos", "carlos","1234",19,19,6,11,1,"h","SDpark");

delete from visitante;

/*Actualizar algunos valores */
UPDATE parque SET aforoMaximo=10 WHERE id = "SDpark";
UPDATE atraccion SET tiempoEspera = 7 WHERE id = "atraccion1";
UPDATE atraccion SET tiempoEspera = 5 WHERE id = "atraccion10";

update visitante set dentroParque = 1 where id = "elbala00";

select * from atraccion;
select * from visitante;

