# Fun with Queues

![Microservices System Architecture View](/images/SystemArchitecture.png "System Architecture View"))

## Descripción

**Fun with Queues** es un sistema de gestión de colas para parques temáticos. El proyecto simula la gestión de colas y el flujo de visitantes en un parque de atracciones, utilizando múltiples microservicios para manejar diferentes aspectos como el registro de visitantes, la gestión de atracciones, la monitorización de tiempos de espera y la actualización del mapa del parque en tiempo real.

## Tabla de Contenidos

1. [Componentes del Proyecto](#componentes-del-proyecto)
   - [parque_atracciones (Base de Datos)](#parque_atracciones-base-de-datos)
   - [Fwq_registry](#fwq_registry)
   - [Fwq_sensor](#fwq_sensor)
   - [Fwq_waitingTimeServer](#fwq_waitingtimeserver)
   - [Fwq_visitor](#fwq_visitor)
   - [Fwq_engine](#fwq_engine)
   - [Api_Engine](#api_engine)
   - [Front-End](#front-end)
2. [Entorno de Desarrollo](#entorno-de-desarrollo)
3. [Guía de Despliegue](#guía-de-despliegue)

## Componentes del Proyecto

### parque_atracciones (Base de Datos)

La base de datos MySQL almacena los datos esenciales del parque, incluyendo:
- **parque**: Información del parque, como su aforo máximo y actual.
- **visitante**: Datos de los visitantes como ID, nombre, contraseña, ubicación actual y destino.
- **atraccion**: Información sobre las atracciones del parque, como su ID, tiempo de ciclo, capacidad de visitantes, y estado.
- **mapa**: Estado actual del mapa del parque.
- **ciudades**: Información básica de las ciudades, incluyendo nombre y temperatura.

### Fwq_registry

Microservicio encargado del registro y autenticación de visitantes. Soporta comunicaciones tanto por sockets como por API REST para gestionar el registro y la modificación de información de los visitantes. Utiliza HTTPS para asegurar la comunicación y almacena contraseñas con un hash para mayor seguridad.

### Fwq_sensor

Microservicio que simula los sensores en las atracciones. Envía el número de personas en la cola de la atracción al servidor de tiempos de espera a intervalos aleatorios.

### Fwq_waitingTimeServer

Servidor encargado de recibir datos de los sensores y calcular los tiempos de espera para cada atracción en el parque. Estos tiempos se envían al **Fwq_engine** para actualizar la información en tiempo real.

### Fwq_visitor

Aplicación que simula a los visitantes del parque. Permite a los usuarios registrarse, iniciar sesión, y navegar por el parque, interactuando con las atracciones y actualizando su posición en el mapa.

### Fwq_engine

Motor central del sistema que gestiona la lógica del flujo de visitantes y las interacciones con las atracciones. Calcula y actualiza el tiempo de espera y el estado de cada atracción basándose en los datos proporcionados por los sensores y el servidor de tiempos de espera.

### Api_Engine

API RESTful que proporciona una interfaz para interactuar con el **Fwq_engine**, permitiendo a otros sistemas y aplicaciones acceder y actualizar datos sobre el estado del parque y sus atracciones.

### Front-End

Interfaz de usuario basada en **ReactJS** que muestra un mapa interactivo del parque, el estado de las atracciones y permite a los visitantes ver y gestionar sus datos.

## Entorno de Desarrollo

- **Editor de código**: Visual Studio Code
- **Lenguaje de programación**: Go (Golang)
- **Bases de datos**: MySQL
- **Transmisión de eventos**: Apache Kafka
- **Interfaz de usuario**: ReactJS
- **Herramientas de administración de bases de datos**: MySQL Workbench

## Guía de Despliegue

Para desplegar la aplicación, sigue estos pasos:

1. **Configura la base de datos MySQL**:
   - Crea las tablas utilizando el script `ParqueAtracciones.sql`.
   - Rellena las tablas con datos iniciales.

2. **Configura y ejecuta los microservicios**:
   - Asegúrate de que Apache Kafka esté en funcionamiento.
   - Configura los microservicios (`Fwq_registry`, `Fwq_sensor`, `Fwq_waitingTimeServer`, `Fwq_engine`) con los parámetros adecuados (IP, puerto, etc.).
   - Ejecuta los microservicios en el orden requerido para asegurar que se comunican correctamente.

3. **Configura y ejecuta la API y el Front-End**:
   - Despliega la API REST (`Api_Engine`) y configura las rutas necesarias.
   - Despliega el front-end basado en ReactJS y asegúrate de que se conecte correctamente a la API.

