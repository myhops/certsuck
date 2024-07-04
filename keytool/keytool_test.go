package keytool

import (
	"context"
	"os"
	"testing"
)

func TestImportCert(t *testing.T) {
	cmd := &ImportCert{
		Alias:     "alias1",
		Keystore:  "testtruststore.jks",
		File:      "/home/peza/DevProjects/certsuck/testresources/www.google.com-00.der",
		Keypass:   "changeme",
		Storepass: "changeme",
	}
	err := cmd.RunContext(context.Background())
	if err != nil {
		t.Errorf("error running keytool: %s", err.Error())
	}
}

func TestImportCerts(t *testing.T) {
	const store = "testtruststore-imports.jks"
	os.Remove(store)
	cmd := &ImportCerts{
		Keystore: store,
		Files: []string{
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-00.der",
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-01.der",
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-02.der",
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-03.der",
		},
		Keypass:   "changeme",
		Storepass: "changeme",
	}
	err := cmd.RunContext(context.Background())
	if err != nil {
		t.Errorf("error running keytool: %s", err.Error())
	}
}

func TestImportCertsWithAliases(t *testing.T) {
	const store = "testtruststore-imports-aliases.jks"
	os.Remove(store)
	cmd := &ImportCerts{
		Keystore: store,
		Aliases: []string{
			"a0",
			"a1",
			"a2",
			"a3",
			"a4",
			"a4",
		},
		Files: []string{
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-00.der",
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-01.der",
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-02.der",
			"/home/peza/DevProjects/certsuck/testresources/www.google.com-03.der",
		},
		Keypass:   "changeme",
		Storepass: "changeme",
	}
	err := cmd.RunContext(context.Background())
	if err != nil {
		t.Errorf("error running keytool: %s", err.Error())
	}
}
