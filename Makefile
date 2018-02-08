build: clean depend build-go

build-go:
		cd cmd/check && go build
		cd cmd/in && go build
		cd cmd/out && go build

depend:
		dep ensure

clean: clean-vendor clean-go

clean-vendor:
		@rm -rf vendor
clean-go:
		go clean
image:
		docker build .

