package mcp

import (
	"testing"
)

func TestCheckK8sConnection(t *testing.T) {
	status := CheckK8sHealth("local")
	if status.Target != "local" {
		t.Errorf("expected target 'local', got '%s'", status.Target)
	}
}

func TestCheckCIConnection(t *testing.T) {
	status := CheckCICDHealth("argocd")
	if !status.Configured {
		t.Error("expected argocd to be marked as configured")
	}
}
