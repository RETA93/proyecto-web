# Proyecto Web — Dev / QA / Prod con Docker

Stack: Python (FastAPI) + Go + Frontend (React/Vue/Angular) + MySQL + Nginx.

## Estructura

```
proyecto-web/
├── docker-compose.yml          # servicios base
├── docker-compose.dev.yml      # override desarrollo (hot-reload, puertos abiertos)
├── docker-compose.qa.yml       # override QA
├── docker-compose.prod.yml     # override producción
├── .env.dev.example            # copiar a .env.dev y llenar valores reales
├── .env.qa.example
├── .env.prod.example
├── backend-python/             # API FastAPI
├── backend-go/                 # API Go
├── frontend/                   # SPA (React/Vue/Angular)
├── nginx/nginx.conf            # reverse proxy
└── .github/workflows/ci-cd.yml # pipeline
```

## 1. Preparar variables de entorno

En el servidor, copia y edita cada plantilla (nunca subas los `.env` reales a git):

```bash
cp .env.dev.example .env.dev
cp .env.qa.example .env.qa
cp .env.prod.example .env.prod
nano .env.prod   # pon contraseñas reales y fuertes
```

## 2. Levantar cada entorno

```bash
# Desarrollo
docker compose -f docker-compose.yml -f docker-compose.dev.yml --env-file .env.dev up -d --build

# QA
docker compose -f docker-compose.yml -f docker-compose.qa.yml --env-file .env.qa up -d --build

# Producción
docker compose -f docker-compose.yml -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

Cada entorno corre en su propio conjunto de contenedores (nombres con sufijo `-dev`, `-qa`, `-prod`) y su propio puerto de Nginx:
- dev → `localhost:8090`
- qa → `localhost:8091`
- prod → `localhost:8092`

## 3. Crear el repositorio en GitHub

```bash
cd proyecto-web
git init
git add .
git commit -m "Estructura inicial: docker, dev/qa/prod, pipeline"
git branch -M main
git checkout -b develop
git checkout -b qa
git checkout main
git remote add origin https://github.com/TU_USUARIO/proyecto-web.git
git push -u origin main develop qa
```

Flujo de ramas sugerido:
- `develop` → despliega a **dev**
- `qa` → despliega a **qa**
- `main` → despliega a **prod**

## 4. Configurar secrets en GitHub

En GitHub → Settings → Secrets and variables → Actions, agrega:
- `SERVER_HOST` → IP o dominio de tu servidor
- `SERVER_USER` → tu usuario SSH (ej. `reta`)
- `SERVER_SSH_KEY` → llave privada SSH con acceso al servidor

El pipeline (`.github/workflows/ci-cd.yml`) hace `git pull` + `docker compose up -d --build` automáticamente al hacer push a cada rama.

## 5. Exponer los 3 entornos con tu Cloudflare Tunnel

Ya tienes el tunnel corriendo para SSH y la web principal. Agrega un hostname por entorno en `/etc/cloudflared/config.yml`:

```yaml
ingress:
  - hostname: ssh.retadev.com
    service: tcp://localhost:22
  - hostname: dev.retadev.com
    service: http://localhost:8090
  - hostname: qa.retadev.com
    service: http://localhost:8091
  - hostname: retadev.com
    service: http://localhost:8092
  - hostname: www.retadev.com
    service: http://localhost:8092
  - service: http_status:404
```

Luego:
```bash
cloudflared tunnel ingress validate
sudo systemctl restart cloudflared
cloudflared tunnel route dns c6cfdea9-9e97-41e4-9e95-81ed4133ba93 dev.retadev.com
cloudflared tunnel route dns c6cfdea9-9e97-41e4-9e95-81ed4133ba93 qa.retadev.com
cloudflared tunnel route dns c6cfdea9-9e97-41e4-9e95-81ed4133ba93 retadev.com
cloudflared tunnel route dns c6cfdea9-9e97-41e4-9e95-81ed4133ba93 www.retadev.com
```

## 6. Verificar

- `https://dev.retadev.com` → entorno de desarrollo
- `https://qa.retadev.com` → entorno de QA
- `https://retadev.com` → producción

## Notas

- El backend Python y Go incluyen un endpoint `/health` para monitoreo.
- MySQL persiste datos en el volumen `mysql-data` — solo se expone el puerto 3306 en `dev` para que puedas conectarte con un cliente local (DBeaver, TablePlus, etc.).
- En `prod`, MySQL **no** expone puertos al host por seguridad.
