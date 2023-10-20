package main

import (
	"bufio"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"flag"
	"fmt"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type ParsedSecret struct {
	envVar     string
	secretName string
	filePath   string
}

func main() {

	fmt.Println("We bout ta inject some secrets.  Shhhh.")
	flag.Parse()
	command := flag.Args()
	if flag.NArg() > 0 {
		fmt.Printf("Then we gonna run: %s\n", command)
	}
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatalf("unable to load credentials: %v", err)
	}
	project := credentials.ProjectID
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	var secrets []*ParsedSecret
	for _, envVar := range os.Environ() {
		parsedSecret, ok := ParseSecret(envVar, project)
		if ok {
			secrets = append(secrets, parsedSecret)
		}
	}

	for _, parsedSecret := range secrets {

		fmt.Printf("processing secret: %s\n", parsedSecret.secretName)
		secret := LoadSecretValue(client, parsedSecret.secretName)
		// remove the parsed secret.  It will be replaced by the actual secret later unless it's mounted
		_ = os.Unsetenv(parsedSecret.envVar)
		if len(parsedSecret.filePath) > 0 {
			WriteSecretFile(parsedSecret.filePath, secret)
		} else {
			_ = os.Setenv(parsedSecret.envVar, secret)
		}
	}
	//fmt.Printf("[%s, %s]\n", *mountPoint, *secretFilename)
	var cmd *exec.Cmd
	if len(command) > 1 {
		cmd = exec.Command(command[0], command[1:]...)
	} else if len(command) == 1 {
		cmd = exec.Command(command[0])
	}

	if cmd != nil {
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		cmd.Start()
		var wg sync.WaitGroup
		wg.Add(2)
		go pipeOutput(stdout, &wg)
		go pipeOutput(stderr, &wg)
		cmd.Wait()
	}
}

func pipeOutput(output io.Reader, wg *sync.WaitGroup) {

	defer wg.Done()

	scanner := bufio.NewScanner(output)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func WriteSecretFile(filePath string, secret string) bool {

	fmt.Printf("attempting to write secret file %s\n", filePath)
	err := os.MkdirAll(filepath.Dir(filePath), 0700)
	if err != nil {
		fmt.Printf("failed to create directory for file: %v\n", err)
		return false
	}
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("failed to create file: %v\n", err)
		return false
	}
	defer file.Close()
	_, err = file.WriteString(secret)
	if err != nil {
		fmt.Printf("failed to write to file: %v\n", err)
		return false
	}

	return true
}

func LoadSecretValue(client *secretmanager.Client, name string) string {

	ctx := context.Background()
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return ""
	}
	return string(result.Payload.Data)
}

func ParseSecret(value string, project string) (*ParsedSecret, bool) {

	result := ParsedSecret{}
	kvp := strings.SplitN(value, "=", 2)
	if !strings.HasPrefix(kvp[1], "secret:") {
		return nil, false
	}
	result.envVar = kvp[0]
	secretParts := strings.Split(kvp[1][7:], "|")
	result.secretName = GenerateSecretUrl(secretParts[0], project)

	if len(secretParts) == 2 {
		result.filePath = secretParts[1]
	}

	return &result, true
}

func GenerateSecretUrl(value string, project string) string {

	if strings.HasPrefix(value, "projects/") {
		return value
	}

	parts := strings.SplitN(value, ":", 2)
	secret := parts[0]
	version := "latest"
	if len(parts) == 2 {
		version = parts[1]
	}

	return fmt.Sprintf("projects/%s/secrets/%s/versions/%s", project, secret, version)
}
