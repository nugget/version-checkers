PROJECT?=	zonedirector-notify
SERVICE?=	$(PROJECT)
HUBUSER?= 	nugget
DOCKERIMAGE?= 	$(HUBUSER)/$(SERVICE)

.PHONY: clean service

all: service

include ../../Makefile.inc

service: gobuild gotest

clean:
	rm -rf zonedirector-notify 

localrun:
	-docker container rm $(DOCKERIMAGE)
	docker run -it --name $(PROJECT) --env-file local.env $(DOCKERIMAGE) -t
