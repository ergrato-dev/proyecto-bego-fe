# Archivo: Dockerfile (fe/)
# Descripción: Build multi-stage para el frontend React + Vite — builder con Node.js y runtime con Nginx.
# ¿Para qué? Compilar la SPA en modo producción y servirla con Nginx de forma eficiente.
# ¿Impacto? La imagen final NO contiene Node.js ni el código fuente — solo el bundle compilado.
#           Esto reduce el tamaño de la imagen de ~1 GB (con Node) a ~50 MB (solo Nginx + dist/).

# ─────────────────────────────────────────────────────────────
# Stage 1: builder — instala dependencias y compila el bundle
# ─────────────────────────────────────────────────────────────
# ¿Qué? Node.js 20.18.2 LTS con Alpine — versión pinada exacta.
# ¿Para qué? Garantizar reproducibilidad; evitar que una actualización silenciosa de Node
#            rompa el build (OWASP A06 — Supply Chain).
FROM node:20.18.2-alpine AS builder

WORKDIR /app

# ¿Qué? Instalar pnpm en la versión exacta usada en desarrollo via npm.
# ¿Para qué? Garantizar que el Dockerfile usa exactamente la misma versión de pnpm
#            que en local — el pnpm-lock.yaml es generado con esta versión.
# ¿Por qué npm en lugar de corepack? Corepack require verificación de firma que puede
#            fallar en entornos con claves desactualizadas; npm install es más robusto.
# ¿Impacto? Una versión diferente de pnpm puede alterar la resolución de dependencias
#           y romper --frozen-lockfile o generar advertencias inesperadas.
RUN npm install -g pnpm@10.30.1 --no-update-notifier

# ¿Qué? Build argument para la URL de la API — configurable en tiempo de build.
# ¿Para qué? Vite incrusta VITE_API_URL en el bundle estático durante la compilación.
#            Cambiar la URL del backend requiere reconstruir la imagen.
# ¿Impacto? Si se usa localhost, el frontend del contenedor apuntará a la máquina host,
#           no al servicio 'be' dentro de Docker — usar la URL pública correcta.
ARG VITE_API_URL=http://localhost:8000
ENV VITE_API_URL=$VITE_API_URL

# ¿Qué? Copiar manifiestos de dependencias ANTES que el código fuente.
# ¿Para qué? Docker cachea esta capa — si package.json y pnpm-lock.yaml no cambian,
#            pnpm install no se repite aunque el código fuente haya cambiado.
# ¿Impacto? Ahorra varios minutos en builds incrementales.
COPY package.json pnpm-lock.yaml ./

# ¿Qué? Instalar dependencias usando exactamente lo que dice el lockfile.
# ¿Para qué? --frozen-lockfile falla si pnpm-lock.yaml está desactualizado,
#            garantizando que el build es determinista.
RUN pnpm install --frozen-lockfile

# ¿Qué? Copiar el resto del código fuente y compilar.
# ¿Para qué? Separado del COPY de manifiestos para maximizar el uso de caché Docker.
COPY . .
RUN pnpm build

# ─────────────────────────────────────────────────────────────
# Stage 2: runtime — Nginx sirve el bundle estático
# ─────────────────────────────────────────────────────────────
# ¿Qué? Nginx 1.27 con Alpine — imagen mínima de ~50 MB para servir archivos estáticos.
# ¿Para qué? No necesitamos Node.js en producción — solo archivos HTML/JS/CSS compilados.
FROM nginx:1.27-alpine AS runtime

# ¿Qué? Configuración personalizada de Nginx (SPA routing + gzip + caché).
# ¿Para qué? La configuración por defecto de Nginx no maneja correctamente React Router.
COPY nginx.conf /etc/nginx/conf.d/default.conf

# ¿Qué? Copiar el bundle compilado desde el stage builder.
# ¿Para qué? Los archivos estáticos generados por Vite en dist/ son todo lo que necesita Nginx.
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
