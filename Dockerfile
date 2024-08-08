# Usar la imagen oficial de Go como imagen base
FROM golang:1.23rc2

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR ./app

# Copiar el código Go al contenedor
COPY . .

# Compilar el binario de Go
RUN go build -o main ./cmd

# Exponer el puerto en el que el servicio escuchará
EXPOSE 8080

# Comando para ejecutar el binario
CMD ["./main"]
