package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkFlow(t *testing.T) {
	db, err:= NewBadgerDB(t.TempDir(), true)
	assert.NoError(t, err)
	defer db.Close()

	wf := NewBadgerWorkFlowStore(db)

	data := map[string]any{"edge": true}

	err = wf.Save("test", data)

	assert.NoError(t, err)

	workflows, err := wf.Load()
	t.Log("workflows: ",workflows)

	assert.NoError(t, err)
	assert.Len(t, workflows, 1)
	if len(workflows) > 0 {
		assert.Equal(t, data, workflows[0])
	}
}