package anet

import "testing"

func TestForceAnetAndroidAlways(t *testing.T) {
	if !forceAnet("android", "") {
		t.Fatal("android must use anet even without FORCE_ANET")
	}
	if !forceAnet("android", "0") {
		t.Fatal("android must use anet even when FORCE_ANET=0")
	}
}

func TestForceAnetOthersHonorEnv(t *testing.T) {
	if forceAnet("linux", "") {
		t.Fatal("linux must not force anet without FORCE_ANET")
	}
	if !forceAnet("linux", "1") {
		t.Fatal("linux must honor FORCE_ANET=1")
	}
	if forceAnet("linux", "0") {
		t.Fatal("linux must honor FORCE_ANET=0")
	}
}
