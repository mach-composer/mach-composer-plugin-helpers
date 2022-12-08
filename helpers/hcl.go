package helpers

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var regexVars = regexp.MustCompilePOSIX(`"\$\$\{([^\}]+)\}"`)

func SerializeToHCL(attributeName string, data any) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	rootBody.SetAttributeValue(attributeName, transformToCTY(data))
	return fixVariableReference(string(f.Bytes()))
}

func fixVariableReference(data string) string {
	matches := regexVars.FindAllStringSubmatch(data, -1)
	for _, match := range matches {
		replacement := match[1]

		// Unescape quotes. Required for secret references, e.g.:
		// 	data.sops_external.variables.data[\"my-key\"]
		// should become:
		// 	data.sops_external.variables.data["my-key"]
		replacement = strings.ReplaceAll(replacement, `\"`, `"`)
		data = strings.Replace(data, match[0], replacement, 1)
	}

	return data
}

func transformToCTY(source any) cty.Value {
	jsonBytes, err := json.Marshal(source)
	if err != nil {
		panic(err)
	}
	var ctyJsonVal ctyjson.SimpleJSONValue
	if err := ctyJsonVal.UnmarshalJSON(jsonBytes); err != nil {
		panic(err)
	}

	return ctyJsonVal.Value
}
