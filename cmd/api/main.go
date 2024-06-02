package main

import (
	"github.com/dembygenesis/local.tools/di/ctn/dic"
	"github.com/dembygenesis/local.tools/internal/api"
	"log"
)

func main() {
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("builder: %vgit remote add ", err)
	}

	ctn := builder.Build()

	cfg, err := ctn.SafeGetConfigLayer()
	if err != nil {
		log.Fatalf("cfg: %v", err)
	}

	_logger, err := ctn.SafeGetLoggerLogrus()
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	categoryMgr, err := ctn.SafeGetLogicCategory()
	if err != nil {
		log.Fatalf("category mgr: %v", err)
	}

	organizationMgr, err := ctn.SafeGetLogicOrganization()
	if err != nil {
		log.Fatalf("organization mgr: %v", err)
	}

	capturePagesMgr, err := ctn.SafeGetLogicCapturePages()
	if err != nil {
		log.Fatalf("capture pages mgr: %v", err)
	}

	apiCfg := &api.Config{
		BaseUrl:             cfg.API.BaseUrl,
		Logger:              _logger,
		Port:                cfg.API.Port,
		CategoryService:     categoryMgr,
		OrganizationService: organizationMgr,
		CapturePagesService: capturePagesMgr,
	}

	_api, err := api.New(apiCfg)
	if err != nil {
		log.Fatalf("api: %v", err)
	}

	if err := _api.Listen(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
