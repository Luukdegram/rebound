package importers

import (
	"fmt"

	"github.com/qmuntal/gltf"
)

//LoadGltfModel imports a GLTF file into a model
func LoadGltfModel() {
	doc, err := gltf.Open("./test.gltf")
	if err != nil {
		panic(err)
	}
	fmt.Print(doc.Asset)
}
