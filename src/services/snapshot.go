package services

import (
	"encoding/json"
	"os"

	"github.com/petruki/gitops-poc/src/model"
)

type DiffType int

const (
	NEW     DiffType = iota
	CHANGED DiffType = iota
	DELETED DiffType = iota
)

func ReadJsonFromFile(path string) string {
	file, _ := os.Open(path)
	defer file.Close()

	stat, _ := file.Stat()
	bs := make([]byte, stat.Size())
	file.Read(bs)
	return string(bs)
}

func NewJsonStringFromSnapshot(snapshot model.Snapshot) string {
	json, _ := json.MarshalIndent(snapshot, "", "  ")
	return string(json)
}

func NewSnapshotFromJson(jsonData []byte) model.Snapshot {
	var snapshot model.Snapshot
	json.Unmarshal(jsonData, &snapshot)
	return snapshot
}

func CheckSnapshotDiff(left model.Snapshot, right model.Snapshot, diffType DiffType) model.Snapshot {
	var diff model.Snapshot
	diff.Domain = model.Domain{}
	return checkGroupDiff(left, right, diffType, diff)
}

func checkGroupDiff(left model.Snapshot, right model.Snapshot, diffType DiffType, diff model.Snapshot) model.Snapshot {
	for _, leftGroup := range left.Domain.Group {
		if !containsValue(model.GroupNames(right.Domain.Group), leftGroup.Name) {
			if diffType == NEW || diffType == DELETED {
				diff.Domain.Group = append(diff.Domain.Group, leftGroup)
			}
		} else {
			rightGroup := model.GetGroupByName(right.Domain.Group, leftGroup.Name)
			modelDiffFound := model.Group{
				Name: leftGroup.Name,
			}

			diffFound := false
			diffFound = compareAndUpdateBool(leftGroup.Activated, rightGroup.Activated, diffFound, &modelDiffFound.Activated)
			diffFound = compareAndUpdateString(leftGroup.Description, rightGroup.Description, diffFound, &modelDiffFound.Description)
			diffFound = compareAndUpdateGroupConfig(leftGroup, rightGroup, diffType, diffFound, &modelDiffFound)

			if diffFound {
				diff.Domain.Group = append(diff.Domain.Group, modelDiffFound)
			}
		}
	}

	return diff
}

func compareAndUpdateBool(left *bool, right *bool, diffFound bool, modelDiffFound **bool) bool {
	if *left != *right {
		diffFound = true
		*modelDiffFound = right
	}
	return diffFound
}

func compareAndUpdateString(left string, right string, diffFound bool, modelDiffFound *string) bool {
	if left != right {
		diffFound = true
		*modelDiffFound = right
	}
	return diffFound
}

func compareAndUpdateGroupConfig(leftGroup model.Group,
	rightGroup model.Group,
	diffType DiffType,
	diffFound bool,
	modelDiffFound *model.Group) bool {
	if len(leftGroup.Config) > 0 {
		modelDiffFound.Config = checkConfigDiff(leftGroup, rightGroup, diffType)
		if len(modelDiffFound.Config) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkConfigDiff(leftGroup model.Group, rightGroup model.Group, diffType DiffType) []model.Config {
	var diff []model.Config

	for _, leftConfig := range leftGroup.Config {
		if !containsValue(model.ConfigKeys(rightGroup.Config), leftConfig.Key) {
			if diffType == NEW || diffType == DELETED {
				diff = append(diff, leftConfig)
			}
		} else {
			rightConfig := model.GetConfigByKey(rightGroup.Config, leftConfig.Key)
			modelDiffFound := model.Config{
				Key: leftConfig.Key,
			}

			diffFound := false

			if diffType == CHANGED {
				diffFound = compareAndUpdateBool(leftConfig.Activated, rightConfig.Activated, diffFound, &modelDiffFound.Activated)
				diffFound = compareAndUpdateString(leftConfig.Description, rightConfig.Description, diffFound, &modelDiffFound.Description)
			}

			diffFound = compareAndUpdateConfigStrategies(leftConfig, rightConfig, diffType, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateConfigComponents(leftConfig, rightConfig, diffType, diffFound, &modelDiffFound)

			if diffFound {
				diff = append(diff, modelDiffFound)
			}
		}
	}

	return diff
}

func compareAndUpdateConfigStrategies(leftConfig model.Config,
	rightConfig model.Config,
	diffType DiffType,
	diffFound bool,
	modelDiffFound *model.Config) bool {
	if len(leftConfig.Strategies) > 0 {
		modelDiffFound.Strategies = checkStrategyDiff(leftConfig, rightConfig, diffType)
		if len(modelDiffFound.Strategies) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkStrategyDiff(leftConfig model.Config, rightConfig model.Config, diffType DiffType) []model.Strategy {
	var diff []model.Strategy

	for _, leftStrategy := range leftConfig.Strategies {
		if !containsValue(model.StrategyNames(rightConfig.Strategies), leftStrategy.Strategy) {
			if diffType == NEW || diffType == DELETED {
				diff = append(diff, leftStrategy)
			}
		} else {
			rightStrategy := model.GetStrategyByName(rightConfig.Strategies, leftStrategy.Strategy)
			modelDiffFound := model.Strategy{
				Strategy: leftStrategy.Strategy,
			}

			diffFound := false

			if diffType == CHANGED {
				diffFound = compareAndUpdateBool(leftStrategy.Activated, rightStrategy.Activated, diffFound, &modelDiffFound.Activated)
				diffFound = compareAndUpdateString(leftStrategy.Operation, rightStrategy.Operation, diffFound, &modelDiffFound.Operation)
			}

			diffFound = compareAndUpdateStrategyValues(leftStrategy, rightStrategy, diffType, diffFound, &modelDiffFound)

			if diffFound {
				diff = append(diff, modelDiffFound)
			}
		}
	}

	return diff
}

func compareAndUpdateStrategyValues(leftStrategy model.Strategy,
	rightStrategy model.Strategy,
	diffType DiffType,
	diffFound bool,
	modelDiffFound *model.Strategy) bool {
	if len(leftStrategy.Values) > 0 {
		modelDiffFound.Values = checkValuesDiff(leftStrategy, rightStrategy, diffType)
		if len(modelDiffFound.Values) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkValuesDiff(leftStrategy model.Strategy, rightStrategy model.Strategy, diffType DiffType) []string {
	var diff []string

	for _, leftValue := range leftStrategy.Values {
		if (diffType == NEW || diffType == DELETED) && !containsValue(rightStrategy.Values, leftValue) {
			diff = append(diff, leftValue)
		}
	}

	return diff
}

func compareAndUpdateConfigComponents(leftConfig model.Config,
	rightConfig model.Config,
	diffType DiffType,
	diffFound bool,
	modelDiffFound *model.Config) bool {
	if len(leftConfig.Components) > 0 {
		modelDiffFound.Components = checkComponentsDiff(leftConfig, rightConfig, diffType)
		if len(modelDiffFound.Components) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkComponentsDiff(leftConfig model.Config, rightConfig model.Config, diffType DiffType) []string {
	var diff []string

	for _, leftComponent := range leftConfig.Components {
		if (diffType == NEW || diffType == DELETED) && !containsValue(rightConfig.Components, leftComponent) {
			diff = append(diff, leftComponent)
		}
	}

	return diff
}

func containsValue(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
