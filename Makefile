GO111MODULE=on
CGO_ENABLED=0

# ========== api ==========

API_DIR=api

.PHONY: build_mac
build_mac: clean_api copy_res_dev
	@echo 'building mac dev...'
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/$(API_DIR)/ ./cmd/api/...
	@echo 'done!'

.PHONY: build_dev
build_dev: clean_api copy_res_dev
	@echo 'building dev...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-s -w" -o bin/$(API_DIR)/ ./cmd/api/...
	@echo 'done!'

.PHONY: build_test
build_test: clean_api copy_res_dev
	@echo 'building dev...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-s -w" -o bin/$(API_DIR)/ ./cmd/api/...
	@echo 'done!'

.PHONY: build_pro
build_pro: clean_api copy_res_pro
	@echo 'building pro...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(API_DIR)/ ./cmd/api/...
	@echo 'done!'


clean_api:
	@rm -rf bin/$(API_DIR)

copy_res_dev: copy_res
	@cp ./cmd/api/conf/dev.yml bin/$(API_DIR)/conf/conf.yaml

copy_res_test: copy_res
	@cp ./cmd/api/conf/test.yml bin/$(API_DIR)/conf/conf.yaml

copy_res_pro: copy_res
	@cp ./cmd/api/conf/prod.yml bin/$(API_DIR)/conf/conf.yaml

copy_res:
	@echo 'copying resource...'
	@mkdir -p bin/$(API_DIR)/conf/i18n
	@cp -r ./cmd/api/conf/i18n/* bin/$(API_DIR)/conf/i18n
	@mkdir -p bin/$(API_DIR)/web
	@cp -r cmd/api/web/* bin/$(API_DIR)/web


#========== migration====
DB_DIR=db

.PHONY: build_pro_db
build_pro_db: clean_db copy_db_res
	@echo 'building pro...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(DB_DIR) ./cmd/migrate/...
	@echo 'done!'

.PHONY: build_db_dev
build_db_dev: clean_db copy_db_res_dev
	@echo 'building dev...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(DB_DIR) ./cmd/migrate/...
	@echo 'done!'

.PHONY: build_db_mac_dev
build_db_mac_dev: clean_db copy_db_res_dev
	@echo 'building dev on mac...'
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(DB_DIR) ./cmd/migrate/...
	@echo 'done!'

clean_db:
	@rm -rf bin/$(DB_DIR)


copy_db_res:
	@echo 'copying db resource...'
	@mkdir -p bin/$(DB_DIR)/conf
	@cp ./cmd/migrate/conf/prod.yml bin/$(DB_DIR)/conf/conf.yaml
#	@mkdir -p bin/$(DB_DIR)/conf/i18n
#	@cp -r ./conf/i18n/* bin/$(DB_DIR)/conf/i18n

copy_db_res_dev:
	@echo 'copying db resource...'
	@mkdir -p bin/$(DB_DIR)/conf
	@cp ./cmd/migrate/conf/dev.yml bin/$(DB_DIR)/conf/conf.yaml

# ========== job ==========

JOB_DIR=job
JOB_NAME=eyen-job

clean_job:
	@rm -rf bin/$(JOB_DIR)

.PHONY: build_job_dev
build_job_dev: clean_job copy_job_res_dev
	@echo 'building dev...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(JOB_DIR)/ ./cmd/job/...
	@echo 'done!'

.PHONY: build_job_pro
build_job_pro: clean_job copy_job_res_pro
	@echo 'building pro...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(JOB_DIR) ./cmd/job/...
	@echo 'done!'

copy_job_res_dev: copy_job_res
	@cp ./cmd/job/conf/dev.yml bin/$(JOB_DIR)/conf/conf.yaml

copy_job_res_pro: copy_job_res
	@cp ./cmd/job/conf/prod.yml bin/$(JOB_DIR)/conf/conf.yaml

copy_job_res:
	@echo 'copying job resource...'
	@mkdir -p bin/$(JOB_DIR)/conf

# ========== portal ==========

PORTAL_DIR=portal
PORTAL_NAME=portal

clean_portal:
	@rm -rf bin/$(PORTAL_DIR)

.PHONY: build_portal_dev
build_portal_dev: clean_portal copy_portal_res_dev
	@echo 'building dev...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(PORTAL_DIR)/ ./cmd/portal/...
	@echo 'done!'

.PHONY: build_portal_pro
build_portal_pro: clean_portal copy_portal_res_pro
	@echo 'building pro...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(PORTAL_DIR)/ ./cmd/portal/...
	@echo 'done!'

copy_portal_res_dev: copy_portal_res
	@cp ./cmd/portal/conf/dev.yml bin/$(PORTAL_DIR)/conf/conf.yaml

copy_portal_res_pro: copy_portal_res
	@cp ./cmd/portal/conf/prod.yml bin/$(PORTAL_DIR)/conf/conf.yaml

copy_portal_res:
	@echo 'copying portal resource...'
	@mkdir -p bin/$(PORTAL_DIR)/conf
	@mkdir -p bin/$(PORTAL_DIR)/conf/i18n
	@cp -r ./cmd/portal/conf/i18n/* bin/$(PORTAL_DIR)/conf/i18n

# ========== cli ==========

CLI_DIR=cli
CLI_NAME=eyen-cli

clean_cli:
	@rm -rf bin/$(CLI_DIR)

.PHONY: build_cli_dev
build_cli_dev: clean_cli copy_cli_res_dev
	@echo 'building cli dev...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(CLI_DIR)/ ./cmd/cli/...
	@echo 'done!'

.PHONY: build_cli_pro
build_cli_pro: clean_cli copy_cli_res_pro
	@echo 'building cli pro...'
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(CLI_DIR)/ ./cmd/cli/...
	@echo 'done!'

.PHONY: build_cli_mac_dev
build_cli_mac_dev: clean_cli copy_cli_res_dev
	@echo 'building cli mac dev...'
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(CLI_DIR)/ ./cmd/cli/...
	@echo 'done!'

.PHONY: build_cli_mac_pro
build_cli_mac_pro: clean_cli copy_cli_res_pro
	@echo 'building cli mac pro...'
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(CLI_DIR)/ ./cmd/cli/...
	@echo 'done!'

copy_cli_res_dev: copy_cli_res
	@cp ./cmd/cli/conf/dev.yml bin/$(CLI_DIR)/conf/conf.yaml

copy_cli_res_pro: copy_cli_res
	@cp ./cmd/cli/conf/prod.yml bin/$(CLI_DIR)/conf/conf.yaml

copy_cli_res:
	@echo 'copying cli resource...'
	@mkdir -p bin/$(CLI_DIR)/conf

