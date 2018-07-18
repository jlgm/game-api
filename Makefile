
setup:
	@docker build -f Dockerfile -t jlgm/game-api .

run:
	@cd docker && docker-compose up