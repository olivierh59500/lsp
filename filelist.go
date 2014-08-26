// filelist.go contains Filelist []Fileinfo definition
// and methods to do the two tasks: filter, sort and
// represent as appropriate to the running mode (flags)
package main

import "os"

// FileList maintains current file items to show
var FileList []FileInfo

// FileListUpdate typed channels contain results of Fileinfo elements in FileList resolved asynchroniously
type FileListUpdate struct {
	i    int // index of update element
	item *FileInfo
	done bool // don't wait for more updates
}

type byType []FileInfo

func (a byType) Len() int      { return len(a) }
func (a byType) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byType) Less(i, j int) bool {
	if a[i].f.IsDir() && !a[j].f.IsDir() {
		return true
	}
	return false
}

func researchFileList(files []os.FileInfo) []FileInfo {
	fileList := make([]FileInfo, len(files))
	results := make(chan FileListUpdate)
	for i, f := range files {
		fileList[i].f = f
		go fileList[i].InvestigateFile(i, results)
	}

	setTimeoutTimer()

	leftNotUpdated := len(files)

	for leftNotUpdated > 0 {
		select {
		case u := <-results:
			if u.done {
				leftNotUpdated--
			}
			fileList[u.i] = *u.item
		case <-timeout:
			leftNotUpdated = 0
		}
	}
	return fileList
}
