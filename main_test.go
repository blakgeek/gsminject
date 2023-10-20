package main

import (
	"reflect"
	"testing"
)

func TestParseSecretIgnoresNonSecrets(t *testing.T) {

	value := "ENV_VAR=SOME_NORMAL:secret"
	project := "project"
	_, ok := ParseSecret(value, project)
	if ok {
		t.Fatalf("Should have been ignored")
	}
}

func TestParseSecretEnvSpecificVersion(t *testing.T) {

	expected := &ParsedSecret{envVar: "ENV_VAR", secretName: "projects/project/secrets/SECRET/versions/1"}
	value := "ENV_VAR=secret:SECRET:1"
	project := "project"
	actual, ok := ParseSecret(value, project)
	if !ok || !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}

func TestParseSecretEnvExplicitPath(t *testing.T) {

	expected := &ParsedSecret{envVar: "ENV_VAR", secretName: "projects/project/secrets/SECRET/versions/latest"}
	value := "ENV_VAR=secret:projects/project/secrets/SECRET/versions/latest"
	project := "project"
	actual, ok := ParseSecret(value, project)
	if !ok || !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Specific version wrong %s != %s, %t", expected, actual, ok)
	}
}

func TestParseSecretEnvDefaultVersion(t *testing.T) {

	expected := &ParsedSecret{envVar: "ENV_VAR", secretName: "projects/project/secrets/SECRET/versions/latest"}
	value := "ENV_VAR=secret:SECRET"
	project := "project"
	actual, ok := ParseSecret(value, project)
	if !ok || !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}

func TestParseSecretMountedSpecificVersion(t *testing.T) {

	expected := &ParsedSecret{envVar: "X", filePath: "/secrets/FOO/bar", secretName: "projects/project/secrets/SECRET/versions/1"}
	value := "X=secret:SECRET:1|/secrets/FOO/bar"
	project := "project"
	actual, ok := ParseSecret(value, project)
	if !ok || !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}

func TestParseSecretMountedDefaultVersion(t *testing.T) {

	expected := &ParsedSecret{envVar: "X", filePath: "/secrets/FOO/bar", secretName: "projects/project/secrets/SECRET/versions/latest"}
	value := "X=secret:SECRET|/secrets/FOO/bar"
	project := "project"
	actual, ok := ParseSecret(value, project)
	if !ok || !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}

func TestParseSecretMountedExplicitPath(t *testing.T) {

	expected := &ParsedSecret{envVar: "X", filePath: "/secrets/FOO/bar", secretName: "projects/project/secrets/SECRET/versions/latest"}
	value := "X=secret:projects/project/secrets/SECRET/versions/latest|/secrets/FOO/bar"
	project := "project"
	actual, ok := ParseSecret(value, project)
	if !ok || !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}

func TestGenerateSecretUrlSpecificVersion(t *testing.T) {

	expected := "projects/project/secrets/SECRET/versions/1"
	value := "SECRET:1"
	project := "project"
	actual := GenerateSecretUrl(value, project)
	if expected != actual {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}

func TestGenerateSecretUrlDefaultVersion(t *testing.T) {

	expected := "projects/project/secrets/SECRET/versions/latest"
	value := "SECRET"
	project := "project"
	actual := GenerateSecretUrl(value, project)
	if expected != actual {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}

func TestGenerateSecretUrlExplicitPath(t *testing.T) {

	expected := "projects/project/secrets/SECRET/versions/latest"
	value := "projects/project/secrets/SECRET/versions/latest"
	project := "project"
	actual := GenerateSecretUrl(value, project)
	if expected != actual {
		t.Fatalf("Specific version wrong %s != %s", expected, actual)
	}
}
