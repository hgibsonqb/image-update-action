package main

import (
	"encoding/json"
	"flag"
	image_update "github.com/fluxcd/image-automation-controller/pkg/update"
	imagev1_reflect "github.com/fluxcd/image-reflector-controller/api/v1beta1"
)

func main() {
	logger := NewLogger()

	var path *string = flag.String("path", ".", "Path to manifests to update, relative to working directory")
	var policy_list_json *string = flag.String("policy-list", "{}", "Image update policy configuration, in json format")
	flag.Parse()

	logger.Info("Parsing flags", "path", *path)
	logger.Info("Parsing flags", "policy-list", *policy_list_json)

	var policy_list imagev1_reflect.ImagePolicyList
	json.Unmarshal([]byte(*policy_list_json), &policy_list)

	result, err := image_update.UpdateWithSetters(logger, *path, *path, policy_list.Items)
	if err != nil {
		logger.Error(err, "error")
	}
	logger.Info("Files updated", "result", result)
}
