build:
	docker build -t userprofile-server .
run:
	docker run -p 8080\:8080 userprofile-server