package services

import (
	"encoding/json"

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
		if !containsValueV2(model.GroupNames(right.Domain.Group), leftGroup.Name) {
			if diffType == NEWV2 || diffType == DELETEDV2 {
				diffDetails := model.DiffDetails{
					Action:  string(diffType),
					Diff:    string(GROUP),
					Path:    []string{},
					Content: leftGroup,
				}
				diffResult.Changes = append(diffResult.Changes, diffDetails)
			}
		} else {
			rightGroup := model.GetGroupByName(right.Domain.Group, leftGroup.Name)
			modelDiffFound := model.Group{}

			diffFound := false
			diffFound = compareAndUpdateBoolV2(leftGroup.Activated, rightGroup.Activated, diffFound, &modelDiffFound.Activated)
			diffFound = compareAndUpdateStringV2(leftGroup.Description, rightGroup.Description, diffFound, &modelDiffFound.Description)
			checkConfigDiffV2(leftGroup, rightGroup, &diffResult, diffType)

			if diffFound {
				diffDetails := model.DiffDetails{
					Action:  string(diffType),
					Diff:    string(GROUP),
					Path:    []string{leftGroup.Name},
					Content: modelDiffFound,
				}
				diffResult.Changes = append(diffResult.Changes, diffDetails)
			}
		}
	}

	return diffResult
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

func checkConfigDiffV2(leftGroup model.Group, rightGroup model.Group, diffResult *model.DiffResult, diffType DiffTypeV2) {
	if len(leftGroup.Config) == 0 {
		return
	}

	for _, leftConfig := range leftGroup.Config {
		if !containsValueV2(model.ConfigKeys(rightGroup.Config), leftConfig.Key) {
			if diffType == NEWV2 || diffType == DELETEDV2 {
				diffDetails := model.DiffDetails{
					Action:  string(diffType),
					Diff:    string(CONFIG),
					Path:    []string{leftGroup.Name},
					Content: leftConfig,
				}
				diffResult.Changes = append(diffResult.Changes, diffDetails)
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
				diffDetails := model.DiffDetails{
					Action:  string(diffType),
					Diff:    string(CONFIG),
					Path:    []string{leftGroup.Name, leftConfig.Key},
					Content: modelDiffFound,
				}
				diffResult.Changes = append(diffResult.Changes, diffDetails)
			}
		}
	}
}

func checkStrategyDiffV2(leftConfig model.Config, rightConfig model.Config, leftGroup model.Group, diffResult *model.DiffResult, diffType DiffTypeV2) {
	if len(leftConfig.Strategies) == 0 {
		return
	}

	for _, leftStrategy := range leftConfig.Strategies {
		if !containsValueV2(model.StrategyNames(rightConfig.Strategies), leftStrategy.Strategy) {
			if diffType == NEWV2 || diffType == DELETEDV2 {
				diffDetails := model.DiffDetails{
					Action:  string(diffType),
					Diff:    string(STRATEGY),
					Path:    []string{leftGroup.Name, leftConfig.Key},
					Content: leftStrategy,
				}
				diffResult.Changes = append(diffResult.Changes, diffDetails)
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
				diffDetails := model.DiffDetails{
					Action:  string(diffType),
					Diff:    string(STRATEGY),
					Path:    []string{leftGroup.Name, leftConfig.Key, leftStrategy.Strategy},
					Content: modelDiffFound,
				}
				diffResult.Changes = append(diffResult.Changes, diffDetails)
			}
		}
	}
}

func checkValuesDiffV2(leftStrategy model.Strategy, rightStrategy model.Strategy, leftGroup model.Group, leftConfig model.Config,
	diffResult *model.DiffResult, diffType DiffTypeV2) {
	var diff []string

	if len(leftStrategy.Values) == 0 {
		return
	}

	for _, leftValue := range leftStrategy.Values {
		if (diffType == NEWV2 || diffType == DELETEDV2) && !containsValueV2(rightStrategy.Values, leftValue) {
			diff = append(diff, leftValue)
		}
	}

	if len(diff) > 0 {
		diffDetails := model.DiffDetails{
			Action:  string(diffType),
			Diff:    string(STRATEGY_VALUE),
			Path:    []string{leftGroup.Name, leftConfig.Key, leftStrategy.Strategy},
			Content: diff,
		}
		diffResult.Changes = append(diffResult.Changes, diffDetails)
	}
}

func checkComponentsDiffV2(leftConfig model.Config, rightConfig model.Config, leftGroup model.Group, diffResult *model.DiffResult, diffType DiffTypeV2) {
	var diff []string

	if len(leftConfig.Components) == 0 {
		return
	}

	for _, leftComponent := range leftConfig.Components {
		if (diffType == NEWV2 || diffType == DELETEDV2) && !containsValueV2(rightConfig.Components, leftComponent) {
			diff = append(diff, leftComponent)
		}
	}

	if len(diff) > 0 {
		diffDetails := model.DiffDetails{
			Action:  string(diffType),
			Diff:    string(COMPONENT),
			Path:    []string{leftGroup.Name, leftConfig.Key},
			Content: diff,
		}
		diffResult.Changes = append(diffResult.Changes, diffDetails)
	}
}

func containsValueV2(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
