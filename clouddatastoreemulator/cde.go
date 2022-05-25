package clouddatastoreemulator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/ory/dockertest/v3"
	// "google.golang.org/api/option"
)

// Config provides some setting values.
type Config struct {
	DockerEndpoint string
	ProjectID      string
	// Options        []option.ClientOption
	Tag string
}

// New Cloud Datastore Emulator spawned or detect and setup.
func New(ctx context.Context, cfg *Config) (*datastore.Client, func(), error) {
	// NOTE: recommend to execute `docker pull google/cloud-sdk:258.0.0` before running test.
	//       because dockertest doesn't have indicator.

	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.Tag == "" {
		cfg.Tag = os.Getenv("GCLOUD_SDK_VERSION")
	}
	if cfg.Tag == "" {
		cfg.Tag = "latest"
	}

	if cfg.ProjectID == "" {
		cfg.ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if cfg.ProjectID == "" {
		cfg.ProjectID = "unittest"
	}

	// check Cloud Datastore Emulator instance before launch it ourselves.
	dsCli, err := checkEmulatorInstance(ctx, cfg)
	if err == context.DeadlineExceeded {
		// ok (not exists)
	} else if err != nil {
		// ignore
	} else {
		return dsCli, func() {}, nil
	}

	// spawn new one...
	pool, err := dockertest.NewPool(cfg.DockerEndpoint)
	if err != nil {
		return nil, nil, err
	}
	pool.MaxWait = 10 * time.Second

	err = os.Setenv("DATASTORE_PROJECT_ID", cfg.ProjectID)
	if err != nil {
		return nil, nil, err
	}

	runOptions := &dockertest.RunOptions{
		Repository: "google/cloud-sdk",
		Tag:        cfg.Tag,
		Cmd: []string{
			"gcloud",
			"--project=" + cfg.ProjectID,
			"beta",
			"emulators",
			"datastore",
			"start",
			"--host-port=0.0.0.0:8081",
			"--no-store-on-disk",
			"--consistency=1.0",
		},
		ExposedPorts: []string{
			"8081",
		},
	}
	resource, err := pool.RunWithOptions(runOptions)
	if err != nil {
		return nil, nil, err
	}

	err = pool.Retry(func() error {
		emulatorHost := fmt.Sprintf("localhost:%s", resource.GetPort("8081/tcp"))
		err = os.Setenv("DATASTORE_EMULATOR_HOST", emulatorHost)
		if err != nil {
			return err
		}

		dsCli, err = checkEmulatorInstance(ctx, cfg)
		log.Printf("mmm")
		return err
	})
	log.Printf("nnn")
	if err != nil {
		log.Printf("ooo")
		return nil, nil, err
	}

	log.Printf("ppp")
	return dsCli, func() { _ = pool.Purge(resource) }, nil
}

func checkEmulatorInstance(ctx context.Context, cfg *Config) (*datastore.Client, error) {
	log.Printf("aaa")
	if os.Getenv("DATASTORE_EMULATOR_HOST") == "" {
		log.Printf("bbb")
		return nil, errors.New("not found datastore emulator")
	}

	log.Printf("ccc")
	dsCli, err := datastore.NewClient(ctx, cfg.ProjectID)
	log.Printf("ddd")
	if err != nil {
		log.Printf("eee")
		return nil, err
	}
	log.Printf("fff")
	q := datastore.NewQuery("__namespace__").KeysOnly().Limit(1)
	log.Printf("ggg")
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	log.Printf("hhh")
	defer cancel()
	log.Printf("iii")
	_, err = dsCli.GetAll(ctx, q, nil)
	log.Printf("jjj")
	if err != nil {
		log.Printf("kkk")
		return nil, err
	}
	log.Printf("lll")

	return dsCli, nil
}
