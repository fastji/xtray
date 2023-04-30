package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	futils "github.com/moqsien/free/pkgs/utils"
)

type KoanfJSON struct{}

func NewJsonParser() *KoanfJSON {
	return &KoanfJSON{}
}

// Unmarshal parses the given JSON bytes.
func (p *KoanfJSON) Unmarshal(b []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Marshal marshals the given config map to JSON bytes.
func (p *KoanfJSON) Marshal(o map[string]interface{}) ([]byte, error) {
	return json.MarshalIndent(o, "", "    ")
}

type Koanfer struct {
	k      *koanf.Koanf
	parser *KoanfJSON
	fpath  string // file path
}

func NewKoanfer(path string) (r *Koanfer) {
	r = &Koanfer{
		k:      koanf.New("."),
		parser: &KoanfJSON{},
		fpath:  path,
	}
	r.initDirs()
	return
}

func (that *Koanfer) initDirs() {
	pDir := filepath.Dir(that.fpath)
	if ok, _ := futils.PathIsExist(pDir); !ok {
		if err := os.MkdirAll(pDir, os.ModePerm); err != nil {
			fmt.Println("Make dir failed: ", err)
		}
	}
}

func (that *Koanfer) Save(obj interface{}) {
	that.k.Load(structs.Provider(obj, "koanf"), nil)
	if b, err := that.k.Marshal(that.parser); err == nil && len(b) > 0 {
		os.WriteFile(that.fpath, b, 0666)
	} else {
		fmt.Println("[Save File Failed]", err)
	}
}

func (that *Koanfer) Load(obj interface{}) {
	err := that.k.Load(file.Provider(that.fpath), that.parser)
	if err != nil {
		fmt.Println("[Load File Failed] ", err)
		return
	}
	that.k.UnmarshalWithConf("", obj, koanf.UnmarshalConf{Tag: "koanf"})
}
