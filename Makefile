build:
	docker build -t folduka:0.2.59 -f docker/Dockerfile src
push:
	docker tag folduka registry.io.imim.cloud/folduka:0.2.59
	docker push registry.io.imim.cloud/folduka:0.2.59
run:
	go run main.go actions.go login.go display.go workflow.go download.go sign.go signage.go smtp.go utils.go websocket.go template.go database.go queue.go

build-osx:
	go build -o folduka main.go 

cluster:
	k3d cluster create -a 3
