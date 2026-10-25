.PHONY: api-run api-test mongo-up mongo-down web-up web-reader web-studio

api-run:
	go run ./apps/api/auth/cmd

api-test:
	go test ./...

mongo-up:
	docker compose up -d mongo

mongo-down:
	docker compose down

web-up:
	$(MAKE) -j2 web-reader web-studio

web-reader:
	cd apps/web/reader && bun run dev --port 5173 --strictPort

web-studio:
	cd apps/web/studio && bun run dev --port 5174 --strictPort
