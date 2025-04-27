package process

import (
	"os"
	"path/filepath"

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

var _ = Describe("collectFrom", func() {
	var tmpDir string

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()
	})

	It("collects only valid processes and parses name and state correctly", func() {
		// Valider Prozess mit Status `Running`
		createProcEntry(tmpDir, "42", "42 (foo) R 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0")

		// Invalider Eintrag. Kein numerisches Verzeichnis
		Expect(os.Mkdir(filepath.Join(tmpDir, "bar"), dirPerm)).To(Succeed())

		// Valider Prozess mit unbekanntem Namen und Status `Zombie`
		createProcEntry(tmpDir, "99", "99 (?) Z 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0")

		got := collectFrom(tmpDir)
		want := []models.Process{
			{PID: 42, Name: "foo", State: "running"},
			{PID: 99, Name: "unknown", State: "zombie"},
		}
		Expect(got).To(Equal(want))
	})
})
