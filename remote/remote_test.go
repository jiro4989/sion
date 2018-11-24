package remote

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupId(t *testing.T) {
	const userFile = "../testdata/passwd"
	const groupFile = "../testdata/group"
	type TD struct {
		fp     string
		id     string
		expect string
	}
	tds := []TD{
		TD{fp: userFile, id: "0", expect: "root"},
		TD{fp: userFile, id: "1001", expect: "test02"},
		TD{fp: userFile, id: "9999", expect: ""},
		TD{fp: groupFile, id: "0", expect: "root"},
		TD{fp: groupFile, id: "1002", expect: "test03"},
		TD{fp: groupFile, id: "9999", expect: ""},
	}
	for _, v := range tds {
		func() {
			r, err := os.Open(v.fp)
			assert.NoError(t, err)
			defer r.Close()

			s, err := lookupId(r, v.id)
			assert.Equal(t, v.expect, s)
			assert.NoError(t, err)
		}()
	}
}
