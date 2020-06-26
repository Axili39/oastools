package protobuf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Axili39/oastools/oasmodel"
)

func TestLoop(t *testing.T) {
	matches, _ := filepath.Glob("tests/*.yaml")
	for _, match := range matches {
		fmt.Println(match)
		info, err := os.Stat(match)
		if os.IsNotExist(err) {
			t.Errorf("File %s does not exists\n", match)
			break
		}

		oa := oasmodel.OpenAPI{}
		err = oa.Load(match)
		if err != nil {
			t.Errorf("error loading %s : %v", match, err)
		}
		output := &bytes.Buffer{}
		err = Components2Proto(&oa, output, "")
		if err != nil {
			t.Errorf("Error loading file %s : %v\n", info.Name(), err)
		}

		// Verify
		// Get Result Filename
		resultFile := strings.Replace(match, ".yaml", ".proto", 1)
		expected, err := ioutil.ReadFile(resultFile)
		if err != nil {
			t.Errorf("Error loading result file %s : %v", resultFile, err)
		}

		// compare
		if string(expected) != output.String() {
			t.Errorf("Result differ for %s \ngot:\n%s\nexpected:\n%s", match, output.String(), string(expected))
		}
	}
}
