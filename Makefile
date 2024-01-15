CADDY_CONTAINER=steganographics-caddy
DOCKER_IMAGE=onlinuxorg-steganographics:latest
DOCKER_CONTAINER=onlinuxorg-steganographics

prep:
	cd assets && unxz --force --keep semantic-*.min.css.xz
	cd assets && unxz --force --keep alpinejs-*.min.js.xz
	cd assets/fonts && unxz --force --keep *.xz

build: prep
	podman build -t $(DOCKER_IMAGE) .

run:
	-podman rm -f $(DOCKER_CONTAINER)
	podman run \
		-d \
		-p 8081:8081 \
		-v $(PWD)/assets:/assets:z \
		--name $(DOCKER_CONTAINER) \
		$(DOCKER_IMAGE)

logs:
	podman logs -f $(DOCKER_CONTAINER)

build-go:
	go get -v
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o steganographics .

run-go-server:
	./steganographics --server

run-caddy:
	-podman rm -f $(CADDY_CONTAINER)
	podman run \
		-p 8081:80 \
		--restart=unless-stopped -it -d \
		--name=$(CADDY_CONTAINER) \
		-v $(PWD)/assets:/usr/share/caddy:z \
		-v $(PWD)/Caddyfile:/etc/caddy/Caddyfile:z \
		docker.io/library/caddy:latest
