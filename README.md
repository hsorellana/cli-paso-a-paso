# Paso a Paso CLI

Esta es una CLI que consulta los datos entregado por el gobierno. La página oficial de los datos es [esta](https://www.gob.cl/coronavirus/pasoapaso#situacioncomunal/)

# Compilar

Por ahora, está solo probado su funcionamiento en __macOS__ y __linux__ (por medio de docker), aunque esto no impide que puede ser compilado para otras plataformas.

Para compilar:

```bash
git clone git@github.com:hsorellana/cli-paso-a-paso.git
cd cli-paso-a-paso
go build -o cli-paso-a-paso
./cli-paso-a-paso // para ejecutar
```

## Compilar usando Docker

```bash
git clone git@github.com:hsorellana/cli-paso-a-paso.git
cd cli-paso-a-paso
docker run --rm -v $(PWD):/go/src/cli-paso-a-paso -w /go/src/cli-paso-a-paso golang:1.15 bash
```

Obs: El uso de los parámetro de docker son:
* `--rm` borra el contenedor después de apagarlo
* `-v` monta un volumen haciendo posible ver el código dentro del contenedor en la carpeta `/go/src/cli-paso-a-paso`
* `-w` define la carpeta por default al entrar al contenedor
* `golang:1.15` define el nombre y la version de la imagen a usar
* `bash` es para ejecutar una consola bash dentro del contenedor cuando esté creado 