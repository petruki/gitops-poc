package services

import (
	"encoding/json"
	"slices"

	"github.com/petruki/gitops-poc/src/model"
)

type DiffTypeV2 string
type DiffResult string

const (
	NEWV2     DiffTypeV2 = "NEW"
	CHANGEDV2 DiffTypeV2 = "CHANGED"
	DELETEDV2 DiffTypeV2 = "DELETED"

	GROUP          DiffResult = "GROUP"
	CONFIG         DiffResult = "CONFIG"
	STRATEGY       DiffResult = "STRATEGY"
	STRATEGY_VALUE DiffResult = "STRATEGY_VALUE"
	COMPONENT      DiffResult = "COMPONENT"
)

func NewJsonStringFromSnapshotV2(snapshot model.DiffResult) string {
	json, _ := json.MarshalIndent(snapshot, "", "  ")
	return string(json)
}

func CheckSnapshotDiffV2(left model.Snapshot, right model.Snapshot, diffType DiffTypeV2) model.DiffResult {
	diffResult := model.DiffResult{}
	return checkGroupDiffV2(left, right, diffType, diffResult)
}

func checkGroupDiffV2(left model.Snapshot, right model.Snapshot, diffType DiffTypeV2, diffResult model.DiffResult) model.DiffResult {
	for _, leftGroup := range left.Domain.Group {
		if !slices.Contains(model.GroupNames(right.Domain.Group), leftGroup.Name) {
			if diffType == NEWV2 || diffType == DELETEDV2 {
				appendDiffResults(string(diffType), string(GROUP), []string{}, leftGroup, &diffResult)
			}
		} else {
			rightGroup := model.GetGroupByName(right.Domain.Group, leftGroup.Name)
			modelDiffFound := model.Group{}

			diffFound := false
			diffFound = compareAndUpdateBoolV2(leftGroup.Activated, rightGroup.Activated, diffFound, &modelDiffFound.Activated)
			diffFound = compareAndUpdateStringV2(leftGroup.Description, rightGroup.Description, diffFound, &modelDiffFound.Description)
			checkConfigDiffV2(leftGroup, rightGroup, &diffResult, diffType)

			if diffFound {
				appendDiffResults(string(diffType), string(GROUP), []string{}, modelDiffFound, &diffResult)
			}
		}
	}

	return diffResult
}

func checkConfigDiffV2(leftGroup model.Group, rightGroup model.Group, diffResult *model.DiffResult, diffType DiffTypeV2) {
	if len(leftGroup.Config) == 0 {
		return
	}

	for _, leftConfig := range leftGroup.Config {
		if !slices.Contains(model.ConfigKeys(rightGroup.Config), leftConfig.Key) {
			if diffType == NEWV2 || diffType == DELETEDV2 {
				appendDiffResults(string(diffType), string(CONFIG), []string{leftGroup.Name}, leftConfig, diffResult)
			}
		} else {
			rightConfig := model.GetConfigByKey(rightGroup.Config, leftConfig.Key)
			modelDiffFound := model.Config{}

			diffFound := false

			if diffType == CHANGEDV2 {
				diffFound = compareAndUpdateBoolV2(leftConfig.Activated, rightConfig.Activated, diffFound, &modelDiffFound.Activated)
				diffFound = compareAndUpdateStringV2(leftConfig.Description, rightConfig.Description, diffFound, &modelDiffFound.Description)
			}

			checkStrategyDiffV2(leftConfig, rightConfig, leftGroup, diffResult, diffType)
			checkComponentsDiffV2(leftConfig, rightConfig, leftGroup, diffResult, diffType)

			if diffFound {
				appendDiffResults(string(diffType), string(CONFIG), []string{leftGroup.Name}, modelDiffFound, diffResult)
			}
		}
	}
}

func checkStrategyDiffV2(leftConfig model.Config, rightConfig model.Config, leftGroup model.Group, diffResult *model.DiffResult, diffType DiffTypeV2) {
	if len(leftConfig.Strategies) == 0 {
		return
	}

	for _, leftStrategy := range leftConfig.Strategies {
		if !slices.Contains(model.StrategyNames(rightConfig.Strategies), leftStrategy.Strategy) {
			if diffType == NEWV2 || diffType == DELETEDV2 {
				appendDiffResults(string(diffType), string(STRATEGY), []string{leftGroup.Name, leftConfig.Key}, leftStrategy, diffResult)
			}
		} else {
			rightStrategy := model.GetStrategyByName(rightConfig.Strategies, leftStrategy.Strategy)
			modelDiffFound := model.Strategy{}

			diffFound := false

			if diffType == CHANGEDV2 {
				diffFound = compareAndUpdateBoolV2(leftStrategy.Activated, rightStrategy.Activated, diffFound, &modelDiffFound.Activated)
				diffFound = compareAndUpdateStringV2(leftStrategy.Operation, rightStrategy.Operation, diffFound, &modelDiffFound.Operation)
			}

			checkValuesDiffV2(leftStrategy, rightStrategy, leftGroup, leftConfig, diffResult, diffType)

			if diffFound {
				appendDiffResults(string(diffType), string(STRATEGY),
					[]string{leftGroup.Name, leftConfig.Key}, modelDiffFound, diffResult)
			}
		}
	}
}

func checkValuesDiffV2(leftStrategy model.Strategy, rightStrategy model.Strategy, leftGroup model.Group, leftConfig model.Config,
	diffResult *model.DiffResult, diffType DiffTypeV2) {

	if len(leftStrategy.Values) == 0 {
		return
	}

	var diff []string
	for _, leftValue := range leftStrategy.Values {
		if (diffType == NEWV2 || diffType == DELETEDV2) && !slices.Contains(rightStrategy.Values, leftValue) {
			diff = append(diff, leftValue)
		}
	}

	if len(diff) > 0 {
		appendDiffResults(string(diffType), string(STRATEGY_VALUE),
			[]string{leftGroup.Name, leftConfig.Key, leftStrategy.Strategy}, diff, diffResult)
	}
}

func checkComponentsDiffV2(leftConfig model.Config, rightConfig model.Config, leftGroup model.Group,
	diffResult *model.DiffResult, diffType DiffTypeV2) {

	if len(leftConfig.Components) == 0 {
		return
	}

	var diff []string
	for _, leftComponent := range leftConfig.Components {
		if (diffType == NEWV2 || diffType == DELETEDV2) && !slices.Contains(rightConfig.Components, leftComponent) {
			diff = append(diff, leftComponent)
		}
	}

	if len(diff) > 0 {
		appendDiffResults(string(diffType), string(COMPONENT), []string{leftGroup.Name, leftConfig.Key}, diff, diffResult)
	}
}

func compareAndUpdateBoolV2(left *bool, right *bool, diffFound bool, modelDiffFound **bool) bool {
	if *left != *right {
		diffFound = true
		*modelDiffFound = right
	}
	return diffFound
}

func compareAndUpdateStringV2(left string, right string, diffFound bool, modelDiffFound *string) bool {
	if left != right {
		diffFound = true
		*modelDiffFound = right
	}
	return diffFound
}

func appendDiffResults(action string, diff string, path []string, content any, diffResult *model.DiffResult) {
	diffResult.Changes = append(diffResult.Changes, model.DiffDetails{
		Action:  action,
		Diff:    diff,
		Path:    path,
		Content: content,
	})
}
