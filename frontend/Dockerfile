FROM node:22-alpine

WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml .npmrc ./
RUN pnpm install --frozen-lockfile --node-linker=hoisted

COPY . .

ARG PUBLIC_BACKEND_URL=http://localhost:8000
ENV PUBLIC_BACKEND_URL=${PUBLIC_BACKEND_URL}

RUN pnpm build

EXPOSE 3000

CMD ["pnpm", "exec", "vite", "preview", "--host", "0.0.0.0", "--port", "3000"]
