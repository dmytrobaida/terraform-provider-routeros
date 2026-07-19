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

func TestContainerDefinesAllLifecycleTimeouts(t *testing.T) {
	timeouts := ResourceContainer().Timeouts
	if timeouts == nil || timeouts.Create == nil || timeouts.Update == nil || timeouts.Delete == nil {
		t.Fatal("container resource must define create, update, and delete timeouts")
	}
}

func TestContainerPullAllowsTransientUnknownState(t *testing.T) {
	for _, state := range containerPullPendingStates {
		if state == "unknown" {
			return
		}
	}
	t.Fatal("container pull must tolerate RouterOS returning no lifecycle flag")
}

func TestContainerStateSupportsRouterOS721Flags(t *testing.T) {
	tests := []struct {
		name string
		item MikrotikItem
		want string
	}{
		{name: "legacy status", item: MikrotikItem{"status": "running"}, want: "running"},
		{name: "running flag", item: MikrotikItem{"running": "true"}, want: "running"},
		{name: "extracting flag", item: MikrotikItem{"extracting": "true"}, want: "extracting"},
		{name: "stopped flag", item: MikrotikItem{"stopped": "true"}, want: "stopped"},
		{name: "missing state", item: MikrotikItem{}, want: "unknown"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := containerState(test.item); got != test.want {
				t.Fatalf("containerState() = %q, want %q", got, test.want)
			}
		})
	}
}
