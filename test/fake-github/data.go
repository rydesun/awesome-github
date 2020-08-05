package fakeg

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const apiJsonUserPath = "./api/users/%s.json"
const htmlAwesomeReadmePath = "./html/README.html"
const apiJsonRepoPath = "./api/repos/%s_%s.json"
const outputPath = "./output/data.json"

type DataHolder struct {
	dir string
}

func NewDataHolder(dir string) DataHolder {
	return DataHolder{
		dir: dir,
	}
}

func (h *DataHolder) GetUser() ([]byte, error) {
	fpath := filepath.Join(h.dir, fmt.Sprintf(apiJsonUserPath, "tester"))
	return ioutil.ReadFile(fpath)
}

func (h *DataHolder) GetHtmlAwesomeReadme() ([]byte, error) {
	fpath := filepath.Join(h.dir, htmlAwesomeReadmePath)
	return ioutil.ReadFile(fpath)
}

func (h *DataHolder) GetJsonRepo(user string, name string) ([]byte, error) {
	fpath := filepath.Join(h.dir, fmt.Sprintf(apiJsonRepoPath, user, name))
	return ioutil.ReadFile(fpath)
}

func (h *DataHolder) GetOutput() ([]byte, error) {
	fpath := filepath.Join(h.dir, outputPath)
	return ioutil.ReadFile(fpath)
}
