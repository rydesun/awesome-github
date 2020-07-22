package fakeg

import (
	"fmt"
	"io/ioutil"
	"path"
)

const apiJsonUserPath = "./api/users/%s.json"
const htmlAwesomeReadmePath = "./html/README.html"
const apiJsonRepoPath = "./api/repos/%s_%s.json"

type DataHolder struct {
	dir string
}

func NewDataHolder(dir string) DataHolder {
	return DataHolder{
		dir: dir,
	}
}

func (h *DataHolder) GetUser() ([]byte, error) {
	fpath := path.Join(h.dir, fmt.Sprintf(apiJsonUserPath, "tester"))
	return ioutil.ReadFile(fpath)
}

func (h *DataHolder) GetHtmlAwesomeReadme() ([]byte, error) {
	fpath := path.Join(h.dir, htmlAwesomeReadmePath)
	return ioutil.ReadFile(fpath)
}

func (h *DataHolder) GetJsonRepo(user string, name string) ([]byte, error) {
	fpath := path.Join(h.dir, fmt.Sprintf(apiJsonRepoPath, user, name))
	return ioutil.ReadFile(fpath)
}
