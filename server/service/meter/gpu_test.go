package meter

import (
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestCountContainerGPUs(t *testing.T) {
	// devices mode
	devs, cnt, all := countContainerGPUs([]container.DeviceRequest{{Capabilities: [][]string{{"gpu"}}, DeviceIDs: []string{"0", "1"}}})
	if all {
		t.Fatalf("unexpected all")
	}
	if cnt != 0 {
		t.Fatalf("cnt expected 0 got %d", cnt)
	}
	if len(devs) != 2 {
		t.Fatalf("deviceIDs expected 2 got %d", len(devs))
	}

	// count mode
	devs, cnt, all = countContainerGPUs([]container.DeviceRequest{{Capabilities: [][]string{{"gpu"}}, Count: 2}})
	if all {
		t.Fatalf("unexpected all")
	}
	if cnt != 2 {
		t.Fatalf("cnt expected 2 got %d", cnt)
	}
	if len(devs) != 0 {
		t.Fatalf("deviceIDs expected 0 got %d", len(devs))
	}

	// all mode
	devs, cnt, all = countContainerGPUs([]container.DeviceRequest{{Capabilities: [][]string{{"gpu"}}, Count: -1}})
	if !all {
		t.Fatalf("expected all")
	}
	if cnt != 0 || len(devs) != 0 {
		t.Fatalf("unexpected values")
	}
}
