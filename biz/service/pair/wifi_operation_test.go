package pair

import (
	"agent/utils/rpi/network"
	"testing"
)

func TestSortWifiList(t *testing.T) {
	input := []*network.ListWifiInfo{
		{SSID: "a", SIGNAL: "30"},
		{SSID: "b", SIGNAL: "80"},
		{SSID: "c", SIGNAL: "50"},
	}

	out := sortWifiList(input)
	if len(out) != 3 {
		t.Fatalf("unexpected length: %d", len(out))
	}
	if out[0].SIGNAL != "80" || out[1].SIGNAL != "50" || out[2].SIGNAL != "30" {
		t.Fatalf("unexpected order: %v, %v, %v", out[0].SIGNAL, out[1].SIGNAL, out[2].SIGNAL)
	}
}

func TestRemoveDuplicatedWifiKeepsStrongest(t *testing.T) {
	input := []*network.ListWifiInfo{
		{SSID: "same", SIGNAL: "10"},
		{SSID: "same", SIGNAL: "90"},
		{SSID: "other", SIGNAL: "20"},
	}

	out := removeDuplicatedWifi(input)
	if len(out) != 2 {
		t.Fatalf("unexpected length: %d", len(out))
	}

	strongest := map[string]string{}
	for _, v := range out {
		strongest[v.SSID] = v.SIGNAL
	}
	if strongest["same"] != "90" {
		t.Fatalf("expected strongest signal for 'same' to be 90, got %s", strongest["same"])
	}
	if strongest["other"] != "20" {
		t.Fatalf("expected signal for 'other' to be 20, got %s", strongest["other"])
	}
}
