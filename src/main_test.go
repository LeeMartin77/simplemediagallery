package main

import (
	"io/fs"
	"testing"
	"time"
)

type MockStatResponse struct {
	fileInfo fs.FileInfo
	err      error
}
type MockReadDirResponse struct {
	fileInfo []fs.DirEntry
	err      error
}

type MockOs struct {
	statInvocations    []string
	statResponses      []MockStatResponse
	readDirInvocations []string
	readDirResponse    []MockReadDirResponse
}

func (mk MockOs) Stat(inv string) (fs.FileInfo, error) {
	res := mk.statResponses[len(mk.statInvocations)]
	mk.statInvocations = append(mk.statInvocations, inv)
	return res.fileInfo, res.err
}

func (mk MockOs) ReadDir(inv string) ([]fs.DirEntry, error) {
	res := mk.readDirResponse[len(mk.statInvocations)]
	mk.readDirInvocations = append(mk.statInvocations, inv)
	return res.fileInfo, res.err
}

type FakeFile struct {
	name  string
	isDir bool
}

func (ff FakeFile) Name() string       { return ff.name }
func (ff FakeFile) IsDir() bool        { return ff.isDir }
func (ff FakeFile) Size() int64        { return 0 }
func (ff FakeFile) Mode() fs.FileMode  { panic("Unimplemented") }
func (ff FakeFile) ModTime() time.Time { panic("Unimplemented") }
func (ff FakeFile) Sys() any           { panic("Unimplemented") }

type FakeReadDir struct {
	name  string
	isDir bool
}

func (frd FakeReadDir) Name() string               { return frd.name }
func (frd FakeReadDir) IsDir() bool                { return frd.isDir }
func (frd FakeReadDir) Type() fs.FileMode          { panic("Unimplemented") }
func (frd FakeReadDir) Info() (fs.FileInfo, error) { panic("Unimplemented") }

func TestMinimumHappyPathDoesntError(t *testing.T) {
	mockOs := MockOs{
		statResponses: append([]MockStatResponse{}, MockStatResponse{
			fileInfo: FakeFile{name: "", isDir: true},
			err:      nil,
		}),
		readDirResponse: append([]MockReadDirResponse{}, MockReadDirResponse{
			fileInfo: append([]fs.DirEntry{}, FakeReadDir{
				name:  "testfilename.jpg",
				isDir: false,
			}),
			err: nil,
		}),
	}
	testHandler := RequestHandlers{
		MediaDirectory: "/testing/root",
		Stat:           mockOs.Stat,
		ReadDir:        mockOs.ReadDir,
	}

	testHandler.getPageData("/")
	// if actualBody != expectedBody {
	// 	t.Errorf("Expected body %s, but got %s", expectedBody, actualBody)
	// }
}
