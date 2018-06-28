GO?=		go

# This command should match the equivalent command found in
# the internal/zeridian/builddata.sh script
VERSION?=	$(shell git describe --always --dirty --tags --abbrev=1)
BRANCH?=    $(shell git branch | grep \* | cut -d ' ' -f2)
BUILDTIME?= $(shell date -u)

DOCKERORG?= nugget
DOCKERENV?= --build-arg dockerver="`docker -v`" --build-arg dockerhost="`hostname -f`" --build-arg dockeruser="`whoami`"

GOPATH?=	$(HOME)/go

SSH_USER?=	core
UNAME?=		$(shell uname)
KUBECTL?=	$(shell which kubectl)

ifneq ($(KUBECTL),)
	K8_CONTEXT?=	$(shell kubectl config view -o jsonpath='{.current-context}')
	K8_NAMESPACE?=	$(shell kubectl config view -o jsonpath='{.contexts[?(@.name == "$(K8_CONTEXT)")].context.namespace}')
endif


.PHONY: gobuild gotest notklocal

notklocal:
ifeq ($(K8_CONTEXT),$(LOCAL_K8))
	@echo
	@echo You are currently in a local Kubernetes context.  You probably want
	@echo to run 'kubectl config use-context <context>' to point to an AWS
	@echo hosted cluster.
	@echo
	@echo Context: $(K8_CONTEXT)
	@echo
	@kubectl cluster-info
	@exit 1
endif

gobuild: 
	@echo Generating this package and subpackages
	${GO} generate ./...
	${GO} build -i -v ./...
	${GO} build 

gotest: 
	${GO} test

docker: 
	@echo Building $(PROJECT) docker image version $(VERSION)
	cd ../.. && docker build -f services/$(PROJECT)/Dockerfile -t $(DOCKERORG)/$(SERVICE) $(DOCKERENV) .
	docker tag $(DOCKERORG)/$(SERVICE) $(DOCKERORG)/$(SERVICE):$(VERSION)
