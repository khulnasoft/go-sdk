// Copyright 2020 The Khulnasoft Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package khulnasoft

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func getKhulnasoftURL() string {
	return os.Getenv("KHULNASOFT_SDK_TEST_URL")
}

func getKhulnasoftToken() string {
	return os.Getenv("KHULNASOFT_SDK_TEST_TOKEN")
}

func getKhulnasoftUsername() string {
	return os.Getenv("KHULNASOFT_SDK_TEST_USERNAME")
}

func getKhulnasoftPassword() string {
	return os.Getenv("KHULNASOFT_SDK_TEST_PASSWORD")
}

func enableRunKhulnasoft() bool {
	r, _ := strconv.ParseBool(os.Getenv("KHULNASOFT_SDK_TEST_RUN_KHULNASOFT"))
	return r
}

func newTestClient() *Client {
	c, _ := NewClient(getKhulnasoftURL(), newTestClientAuth())
	return c
}

func newTestClientAuth() ClientOption {
	token := getKhulnasoftToken()
	if token == "" {
		return SetBasicAuth(getKhulnasoftUsername(), getKhulnasoftPassword())
	}
	return SetToken(getKhulnasoftToken())
}

func khulnasoftMasterPath() string {
	switch runtime.GOOS {
	case "darwin":
		return fmt.Sprintf("https://dl.khulnasoft.io/khulnasoft/master/khulnasoft-master-darwin-10.6-%s", runtime.GOARCH)
	case "linux":
		return fmt.Sprintf("https://dl.khulnasoft.io/khulnasoft/master/khulnasoft-master-linux-%s", runtime.GOARCH)
	case "windows":
		return fmt.Sprintf("https://dl.khulnasoft.io/khulnasoft/master/khulnasoft-master-windows-4.0-%s.exe", runtime.GOARCH)
	}
	return ""
}

func downKhulnasoft() (string, error) {
	for i := 3; i > 0; i-- {
		resp, err := http.Get(khulnasoftMasterPath())
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		f, err := os.CreateTemp(os.TempDir(), "khulnasoft")
		if err != nil {
			continue
		}
		_, err = io.Copy(f, resp.Body)
		f.Close()
		if err != nil {
			continue
		}

		if err = os.Chmod(f.Name(), 0o700); err != nil {
			return "", err
		}

		return f.Name(), nil
	}

	return "", fmt.Errorf("Download khulnasoft from %v failed", khulnasoftMasterPath())
}

func runKhulnasoft() (*os.Process, error) {
	log.Println("Downloading Khulnasoft from", khulnasoftMasterPath())
	p, err := downKhulnasoft()
	if err != nil {
		log.Fatal(err)
	}

	khulnasoftDir := filepath.Dir(p)
	cfgDir := filepath.Join(khulnasoftDir, "custom", "conf")
	err = os.MkdirAll(cfgDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := os.Create(filepath.Join(cfgDir, "app.ini"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = cfg.WriteString(`[security]
INTERNAL_TOKEN = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjE1NTg4MzY4ODB9.LoKQyK5TN_0kMJFVHWUW0uDAyoGjDP6Mkup4ps2VJN4
INSTALL_LOCK   = true
SECRET_KEY     = 2crAW4UANgvLipDS6U5obRcFosjSJHQANll6MNfX7P0G3se3fKcCwwK3szPyGcbo
[database]
DB_TYPE  = sqlite3
[log]
MODE = console
LEVEL = Trace
REDIRECT_MACARON_LOG = true
MACARON = ,
ROUTER = ,`)
	cfg.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Run khulnasoft migrate", p)
	err = exec.Command(p, "migrate").Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Run khulnasoft admin", p)
	err = exec.Command(p, "admin", "create-user", "--username=test01", "--password=test01", "--email=test01@khulnasoft.io", "--admin=true", "--must-change-password=false", "--access-token").Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start Khulnasoft", p)
	return os.StartProcess(filepath.Base(p), []string{}, &os.ProcAttr{
		Dir: khulnasoftDir,
	})
}

func TestMain(m *testing.M) {
	if enableRunKhulnasoft() {
		p, err := runKhulnasoft()
		if err != nil {
			log.Fatal(err)
			return
		}
		defer func() {
			if err := p.Kill(); err != nil {
				log.Fatal(err)
			}
		}()
	}
	log.Printf("testing with %v, %v, %v\n", getKhulnasoftURL(), getKhulnasoftUsername(), getKhulnasoftPassword())
	exitCode := m.Run()
	os.Exit(exitCode)
}
