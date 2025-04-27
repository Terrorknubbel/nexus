package process

import (
	"os"
	"path/filepath"

	"nexus/internal/config"
	"nexus/pkg/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	dirPerm  = 0755 // rwxr-xr-x
	filePerm = 0644 // rw-r--r--
)

func createProcEntry(root, name, statContent string) {
	d := filepath.Join(root, name)
	Expect(os.Mkdir(d, dirPerm)).To(Succeed())
	path := filepath.Join(d, "stat")
	Expect(os.WriteFile(path, []byte(statContent), filePerm)).To(Succeed())
}

var _ = Describe("Collect", func() {
	var (
		tmpDir    string
		collector *BasicCollector
	)

	tmpDir = GinkgoT().TempDir()
	cfg := &config.Config{ProcRoot: tmpDir}
	collector = NewBasicCollector(cfg)

	It("should collect only valid processes and parse name and state correctly", func() {
		// Valider Prozess mit Status `Running`
		createProcEntry(tmpDir, "42", "42 (foo) R 0")

		// Invalider Eintrag. Kein numerisches Verzeichnis
		Expect(os.Mkdir(filepath.Join(tmpDir, "bar"), dirPerm)).To(Succeed())

		// Valider Prozess mit unbekanntem Namen und Status `Zombie`
		createProcEntry(tmpDir, "99", "99 (?) Z 0")

		procs, err := collector.Collect(nil)
		Expect(err).ToNot(HaveOccurred())

		want := []models.Process{
			{PID: 42, Name: "foo", State: "running"},
			{PID: 99, Name: "unknown", State: "zombie"},
		}
		Expect(procs).To(Equal(want))
	})
})
