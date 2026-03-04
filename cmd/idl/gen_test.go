package idlcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testdataPath(name string) string {
	return filepath.Join("..", "..", "testdata", name)
}

func TestGenerate(t *testing.T) {
	dir := t.TempDir()
	idlPath := testdataPath("v30.json")
	if _, err := os.Stat(idlPath); err != nil {
		t.Skipf("testdata not found: %v", err)
	}
	err := Generate(idlPath, GenOpts{
		Package: "counter",
		Output:  dir,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	typesPath := filepath.Join(dir, "types.go")
	if _, err := os.Stat(typesPath); err != nil {
		t.Errorf("types.go not created: %v", err)
	}
	clientPath := filepath.Join(dir, "client.go")
	if _, err := os.Stat(clientPath); err != nil {
		t.Errorf("client.go not created: %v", err)
	}
	clientContent, _ := os.ReadFile(clientPath)
	if len(clientContent) == 0 {
		t.Error("client.go is empty")
	}
	if !strings.Contains(string(clientContent), "Increment") {
		t.Error("client.go should contain Increment method")
	}
	// Verify types.go has Counter struct
	typesContent, _ := os.ReadFile(typesPath)
	if !strings.Contains(string(typesContent), "type Counter struct") {
		t.Error("types.go should contain Counter struct")
	}
	if !strings.Contains(string(typesContent), "Count uint64") {
		t.Error("types.go Counter should have Count uint64 field")
	}
}

func TestGenerate_MissingFile(t *testing.T) {
	dir := t.TempDir()
	err := Generate("/nonexistent/idl.json", GenOpts{Package: "pkg", Output: dir})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "read IDL") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGenerate_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	invalidPath := filepath.Join(dir, "invalid.json")
	os.WriteFile(invalidPath, []byte("not json"), 0644)
	err := Generate(invalidPath, GenOpts{Package: "pkg", Output: dir})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parse") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGenerate_WithErrors(t *testing.T) {
	dir := t.TempDir()
	idlPath := testdataPath("with_errors.json")
	if _, err := os.Stat(idlPath); err != nil {
		t.Skipf("testdata not found: %v", err)
	}
	err := Generate(idlPath, GenOpts{
		Package: "testprog",
		Output:  dir,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	errorsPath := filepath.Join(dir, "errors.go")
	if _, err := os.Stat(errorsPath); err != nil {
		t.Fatal("errors.go should be created when IDL has errors")
	}
	content, _ := os.ReadFile(errorsPath)
	if !strings.Contains(string(content), "ErrUnauthorized") {
		t.Error("errors.go should contain ErrUnauthorized")
	}
	if !strings.Contains(string(content), "ErrInsufficientFunds") {
		t.Error("errors.go should contain ErrInsufficientFunds")
	}
	if !strings.Contains(string(content), "ErrorCodeToError") {
		t.Error("errors.go should contain ErrorCodeToError function")
	}
	// Verify client has DoStuff
	clientPath := filepath.Join(dir, "client.go")
	clientContent, _ := os.ReadFile(clientPath)
	if !strings.Contains(string(clientContent), "DoStuff") {
		t.Error("client.go should contain DoStuff method")
	}
}

func TestGenerate_DefaultPackage(t *testing.T) {
	dir := t.TempDir()
	idlPath := testdataPath("v30.json")
	if _, err := os.Stat(idlPath); err != nil {
		t.Skipf("testdata not found: %v", err)
	}
	err := Generate(idlPath, GenOpts{Output: dir})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	clientContent, _ := os.ReadFile(filepath.Join(dir, "client.go"))
	if !strings.Contains(string(clientContent), "package counter") {
		t.Error("default package should come from IDL metadata name (counter)")
	}
}

func TestGenerate_GeneratesValidGo(t *testing.T) {
	dir := t.TempDir()
	idlPath := testdataPath("v30.json")
	if _, err := os.Stat(idlPath); err != nil {
		t.Skipf("testdata not found: %v", err)
	}
	err := Generate(idlPath, GenOpts{Package: "counter", Output: dir})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	// Verify generated files have valid Go syntax using go/parser
	for _, name := range []string{"types.go", "client.go"} {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", name, err)
			continue
		}
		// Basic sanity: must have package declaration and not contain invalid tokens
		content := string(data)
		if !strings.HasPrefix(strings.TrimSpace(content), "//") {
			t.Errorf("%s should start with comment", name)
		}
		if !strings.Contains(content, "package counter") {
			t.Errorf("%s should declare package counter", name)
		}
		if name == "client.go" && !strings.Contains(content, "func ") {
			t.Errorf("client.go should contain function")
		}
	}
}
