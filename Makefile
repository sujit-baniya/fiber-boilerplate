dpl ?= .env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d%H%M%S")
VERSION := v$(TAG:v%=%)-$(DATE)-$(COMMIT)
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)
RESET := $(shell tput -Txterm sgr0)
APPLICATION_NAME := $(shell echo $(APP_NAME) | sed -e 's/[^[:alnum:]]/-/g' | tr -s '-' | tr A-Z a-z)

ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

FLAGS := -ldflags "-X github.com/sujit-baniya/fiber-boilerplate/app.Version=$(VERSION)"
BUILD_PATH := $(shell pwd)/build
PID := $(shell lsof -t -i:$(APP_PORT))
RELEASE_PATH := $(BUILD_PATH)/releases
SHARED_PATH := $(BUILD_PATH)/shared
CURRENT_PATH := $(BUILD_PATH)/current
CURRENT_RELEASE := "dev"
PREVIOUS_RELEASE := "dev"
ifneq ("$(wildcard $(CURRENT_PATH)/CURRENT-RELEASE)","")
    CURRENT_RELEASE := $(shell cat $(CURRENT_PATH)/CURRENT-RELEASE)
endif
ifneq ("$(wildcard $(CURRENT_PATH)/PREVIOUS-RELEASE)","")
    PREVIOUS_RELEASE := $(shell cat $(CURRENT_PATH)/PREVIOUS-RELEASE)
endif

ROLLBACK_RELEASE := $(RELEASE_PATH)/$(RELEASE_TAG)

LATEST_RELEASE := $(APPLICATION_NAME)-$(VERSION)
LATEST_RELEASE_PATH := $(RELEASE_PATH)/$(LATEST_RELEASE)

create-folder:
	$(info $(GREEN)Create Release Folder: $(LATEST_RELEASE)$(RESET))
	$(shell mkdir -p $(RELEASE_PATH)/$(LATEST_RELEASE))
	$(shell mkdir -p $(SHARED_PATH)/$(STORAGE_PATH))
	$(shell mkdir -p $(SHARED_PATH)/$(UPLOAD_PATH))

git-stash:
	$(info $(GREEN)Stashing current changes$(RESET))
	cd $(LATEST_RELEASE_PATH) && git stash

git-checkout:
	$(info $(GREEN)Checkingout Master branch$(RESET))
	cd $(LATEST_RELEASE_PATH) && git checkout master && git pull origin master

git-push:
	$(info $(GREEN)Adding all changed files and push $(RESET))
	cd $(LATEST_RELEASE_PATH) && git add . && git commit -m $(COMMIT_MESSAGE) && git push origin master

dev-push:
	$(info $(GREEN)Adding all changed files and push for dev $(RESET))
	git add . && git commit -m "$(COMMIT_MESSAGE)" && git push origin master

build-app:
	$(info $(GREEN)Building the application: $(APPLICATION_NAME)$(RESET))
	$(shell go build $(FLAGS) -o  $(RELEASE_PATH)/$(LATEST_RELEASE)/$(APPLICATION_NAME) main.go)

copy-config:
	$(info $(GREEN)Copying config, assets and .env file$(RESET))
	$(shell cp .env  $(RELEASE_PATH)/$(LATEST_RELEASE)/ && \
			cp config.yml  $(RELEASE_PATH)/$(LATEST_RELEASE)/ && \
			cp -R assets  $(RELEASE_PATH)/$(LATEST_RELEASE)/ \
	)

install-fe-dependencies:
	$(info $(GREEN)Installing Frontend dependencies$(RESET))
	$(shell yarn install >/dev/null)

compile-fe:
	$(info $(GREEN)Compiling Frontend assets$(RESET))
	$(shell yarn install >/dev/null && \
			yarn run prod >/dev/null || true \
	)

copy-assets:
	$(info $(GREEN)Copying Assets$(RESET))
		$(shell cp -R public $(RELEASE_PATH)/$(LATEST_RELEASE)/ && \
    			cp -R resources $(RELEASE_PATH)/$(LATEST_RELEASE)/ \
    	)

create-symlink:
	$(info $(GREEN)Creating Current folder symlink$(RESET))
	$(shell ln -snf $(SHARED_PATH)/$(STORAGE_PATH) $(RELEASE_PATH)/$(LATEST_RELEASE)/$(STORAGE_PATH))
	$(shell ln -snf $(SHARED_PATH)/$(UPLOAD_PATH) $(RELEASE_PATH)/$(LATEST_RELEASE)/$(UPLOAD_PATH))
	$(shell ln -snf $(RELEASE_PATH)/$(LATEST_RELEASE) $(CURRENT_PATH))
	$(shell echo $(CURRENT_RELEASE) > $(CURRENT_PATH)/PREVIOUS-RELEASE)
	$(shell echo $(LATEST_RELEASE_PATH) > $(CURRENT_PATH)/CURRENT-RELEASE)

migrate:
	$(info $(GREEN)Starting migrating$(RESET))
	cd $(CURRENT_PATH) && ./$(APPLICATION_NAME) --migrate

build: create-folder git-checkout build-app copy-config copy-assets create-symlink migrate

deploy: build restart

push: install-fe-dependencies compile-fe dev-push

start:
	$(info $(GREEN)Starting application$(RESET))
	cd $(CURRENT_PATH) && ./$(APPLICATION_NAME) </dev/null &>/dev/null &

kill:
ifneq ($(PID),)
	$(info $(RED)Stopping application on port $(APP_PORT)$(RESET))
	kill -9 $(PID)
else
	$(info $(YELLOW)Application not found on port $(APP_PORT)$(RESET))
endif

restart: kill start

rollback:
ifneq ($(PREVIOUS_RELEASE),)
	$(info Rolling Back to Previous Release: $(PREVIOUS_RELEASE))
	$(shell ln -snf $(PREVIOUS_RELEASE) $(CURRENT_PATH))
endif

rollback-to:
ifneq ($(wildcard $(ROLLBACK_RELEASE)),)
	$(info Rolling Back to Release: $(ROLLBACK_RELEASE))
	$(shell ln -snf $(ROLLBACK_RELEASE) $(CURRENT_PATH))
endif

run:
	go run $(FLAGS) main.go

install:
	go install $(FLAGS)
