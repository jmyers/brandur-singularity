package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	singularity "github.com/jmyers/brandur-singularity"
	"github.com/jmyers/brandur-singularity/pool"
	assert "github.com/stretchr/testify/require"
)

func TestEnsureSymlink(t *testing.T) {
	dir, err := os.MkdirTemp("", "symlink")
	assert.NoError(t, err)

	source := path.Join(dir, "source")
	err = os.WriteFile(source, []byte("source"), 0755)
	assert.NoError(t, err)

	dest := path.Join(dir, "symlink-dest")

	//
	// Case 1: Symlink does not exist
	//

	err = ensureSymlink(source, dest)
	assert.NoError(t, err)

	actual, _ := os.Readlink(dest)
	assert.Equal(t, source, actual)

	//
	// Case 2: Symlink does exist
	//
	// Consists solely of re-running the previous test case.
	//

	err = ensureSymlink(source, dest)
	assert.NoError(t, err)

	actual, _ = os.Readlink(dest)
	assert.Equal(t, source, actual)

	//
	// Case 3: Symlink file exists, but source doesn't
	//

	err = os.RemoveAll(dest)
	assert.NoError(t, err)

	source = path.Join(dir, "source")
	err = os.WriteFile(source, []byte("source"), 0755)
	assert.NoError(t, err)

	err = ensureSymlink(source, dest)
	assert.NoError(t, err)

	actual, _ = os.Readlink(dest)
	assert.Equal(t, source, actual)
}

func TestGetLocals(t *testing.T) {
	locals := getLocals("Title", map[string]interface{}{
		"Foo": "Bar",
	})

	assert.Equal(t, "Bar", locals["Foo"])
	assert.Equal(t, singularity.Release, locals["Release"])
	assert.Equal(t, "Title", locals["Title"])
}

func TestIsHidden(t *testing.T) {
	assert.Equal(t, true, isHidden(".gitkeep"))
	assert.Equal(t, false, isHidden("article"))
}

func TestRunTasks(t *testing.T) {
	conf.Concurrency = 3

	//
	// Success case
	//

	tasks := []*pool.Task{
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return nil }),
	}
	assert.Equal(t, true, runTasks(tasks))

	//
	// Failure case (1 error)
	//

	tasks = []*pool.Task{
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
	}
	assert.Equal(t, false, runTasks(tasks))

	//
	// Failure case (11 errors)
	//
	// Here we'll exit with a "too many errors" message.
	//

	tasks = []*pool.Task{
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
	}
	assert.Equal(t, false, runTasks(tasks))
}
