package resource

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEmbeddedPageIncludesUnderscoreAssets(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("page", "assets"))
	if os.IsNotExist(err) {
		t.Skip("front-end assets are not built")
	}
	if err != nil {
		t.Fatal(err)
	}

	var assetName string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), "_") {
			assetName = entry.Name()
			break
		}
	}
	if assetName == "" {
		t.Skip("front-end build has no underscore-prefixed assets")
	}

	file, err := GetPageFS().Open(filepath.ToSlash(filepath.Join("assets", assetName)))
	if err != nil {
		t.Fatalf("embedded front-end assets must include %s: %v", assetName, err)
	}
	_ = file.Close()
}
