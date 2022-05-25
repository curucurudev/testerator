package clouddatastoreemulator_test

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/favclip/testerator/clouddatastoreemulator"
)

// Sample entity.
type Sample struct {
	Name string
}

func TestNew(t *testing.T) {
	t.Log("111")
	ctx := context.Background()
	t.Log("222")
	dsCli, stop, err := clouddatastoreemulator.New(ctx, nil)
	t.Log("333")
	if err != nil {
		t.Log("444")
		t.Fatal(err)
	}
	t.Log("555")
	defer stop()
	if dsCli == nil {
		t.Fatal("dsCli is nil")
	}

	cdsCli, err := datastore.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
	if err != nil {
		t.Fatal(err)
	}

	key, err := cdsCli.Put(ctx, datastore.NameKey("Test", "Test", nil), &Sample{Name: "vvakame"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(key)
}
