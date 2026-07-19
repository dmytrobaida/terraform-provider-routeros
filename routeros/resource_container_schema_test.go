package routeros

import "testing"

func TestContainerSchemaUsesRouterOS721MountFields(t *testing.T) {
	schema := ResourceContainer().Schema

	if _, ok := schema["mountlists"]; !ok {
		t.Fatal("container schema must expose mountlists")
	}
	if _, ok := schema["mounts"]; ok {
		t.Fatal("container schema must not expose obsolete mounts")
	}
}

func TestContainerMountSchemaUsesRouterOS721ListField(t *testing.T) {
	schema := ResourceContainerMounts().Schema

	if _, ok := schema["list"]; !ok {
		t.Fatal("container mount schema must expose list")
	}
	if _, ok := schema["name"]; ok {
		t.Fatal("container mount schema must not expose obsolete name")
	}
}
