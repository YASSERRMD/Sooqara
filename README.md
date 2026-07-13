# Sooqara

E-commerce listing factory powered by Agnes AI — vision analysis, copy generation, image variants, and async video, all orchestrated end-to-end.

## Quickstart

```bash
cp .env.example .env
# edit .env and set AGNES_API_KEY
make build
./bin/sooqara
```

## Building Releases

```bash
make release-version VERSION=v1.0.0 COMMIT=$(git rev-parse --short HEAD)
```

## Licence

Apache-2.0
