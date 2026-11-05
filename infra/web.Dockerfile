# Сборка веб-панели: pnpm build → статика в nginx. Контекст сборки — корень репозитория,
# потому что генератор токенов и контракт лежат вне web/.
FROM node:24-alpine AS build
RUN corepack enable && corepack prepare pnpm@11 --activate
WORKDIR /src
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml web/
RUN cd web && pnpm install --frozen-lockfile
COPY api ./api
COPY design ./design
COPY web ./web
RUN cd web && pnpm build

FROM nginx:1.27-alpine
COPY infra/nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/web/dist /usr/share/nginx/html
RUN mkdir -p /srv/apk
EXPOSE 80
HEALTHCHECK --interval=10s --timeout=3s --retries=3 CMD wget -qO- http://127.0.0.1/index.html > /dev/null || exit 1
