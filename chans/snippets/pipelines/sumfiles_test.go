package pipelines

import (
	"fmt"
	"testing"
)

func TestSumFiles(t *testing.T) {
	main := func() {
		m, err := md5All(".")
		if err != nil {
			t.Fatal(err)
		}

		paths, bytes := sortRes(m)

		for i := range paths {
			fmt.Printf("%20s	%16x\n", paths[i], bytes[i])
		}
	}
	main()
}
