package services

import (
	"encoding/json"
	"os"

	"github.com/petruki/gitops-poc/src/model"
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

func CheckSnapshotDiff(left model.Snapshot, right model.Snapshot) model.Snapshot {
	var diff model.Snapshot
	diff.Domain = model.Domain{}
	return checkGroupDiff(left, right, diff)
}

func checkGroupDiff(left model.Snapshot, right model.Snapshot, diff model.Snapshot) model.Snapshot {
	for _, leftGroup := range left.Domain.Group {
		if !containsGroup(right.Domain.Group, leftGroup) {
			diff.Domain.Group = append(diff.Domain.Group, leftGroup)
		} else {
			rightGroup := getGroupByName(right.Domain.Group, leftGroup.Name)
			modelDiffFound := model.Group{
				Name: leftGroup.Name,
			}

			diffFound := false
			diffFound = compareAndUpdateGroupActivated(leftGroup, rightGroup, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateGroupDescription(leftGroup, rightGroup, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateGroupConfig(leftGroup, rightGroup, diffFound, &modelDiffFound)

			if diffFound {
				diff.Domain.Group = append(diff.Domain.Group, modelDiffFound)
			}
		}
	}

	return diff
}

func compareAndUpdateGroupActivated(leftGroup model.Group,
	rightGroup model.Group,
	diffFound bool,
	modelDiffFound *model.Group) bool {
	if *leftGroup.Activated != *rightGroup.Activated {
		diffFound = true
		modelDiffFound.Activated = rightGroup.Activated
	}
	return diffFound
}

func compareAndUpdateGroupDescription(leftGroup model.Group,
	rightGroup model.Group,
	diffFound bool,
	modelDiffFound *model.Group) bool {
	if leftGroup.Description != rightGroup.Description {
		diffFound = true
		modelDiffFound.Description = rightGroup.Description
	}
	return diffFound
}

func compareAndUpdateGroupConfig(leftGroup model.Group,
	rightGroup model.Group,
	diffFound bool,
	modelDiffFound *model.Group) bool {
	if len(leftGroup.Config) > 0 {
		modelDiffFound.Config = checkConfigDiff(leftGroup, rightGroup)
		if len(modelDiffFound.Config) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkConfigDiff(leftGroup model.Group, rightGroup model.Group) []model.Config {
	var diff []model.Config

	for _, leftConfig := range leftGroup.Config {
		if !containsConfig(rightGroup.Config, leftConfig) {
			diff = append(diff, leftConfig)
		} else {
			rightConfig := getConfigByKey(rightGroup.Config, leftConfig.Key)
			modelDiffFound := model.Config{
				Key: leftConfig.Key,
			}

			diffFound := false
			diffFound = compareAndUpdateConfigActivated(leftConfig, rightConfig, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateConfigDescription(leftConfig, rightConfig, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateConfigStrategies(leftConfig, rightConfig, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateConfigComponents(leftConfig, rightConfig, diffFound, &modelDiffFound)

			if diffFound {
				diff = append(diff, modelDiffFound)
			}
		}
	}

	return diff
}

func compareAndUpdateConfigActivated(leftConfig model.Config,
	rightConfig model.Config,
	diffFound bool,
	modelDiffFound *model.Config) bool {
	if *leftConfig.Activated != *rightConfig.Activated {
		diffFound = true
		modelDiffFound.Activated = rightConfig.Activated
	}
	return diffFound
}

func compareAndUpdateConfigDescription(leftConfig model.Config,
	rightConfig model.Config,
	diffFound bool,
	modelDiffFound *model.Config) bool {
	if leftConfig.Description != rightConfig.Description {
		diffFound = true
		modelDiffFound.Description = rightConfig.Description
	}
	return diffFound
}

func compareAndUpdateConfigStrategies(leftConfig model.Config,
	rightConfig model.Config,
	diffFound bool,
	modelDiffFound *model.Config) bool {
	if len(leftConfig.Strategies) > 0 {
		modelDiffFound.Strategies = checkStrategyDiff(leftConfig, rightConfig)
		if len(modelDiffFound.Strategies) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkStrategyDiff(leftConfig model.Config, rightConfig model.Config) []model.Strategy {
	var diff []model.Strategy

	for _, leftStrategy := range leftConfig.Strategies {
		if !containsStrategy(rightConfig.Strategies, leftStrategy) {
			diff = append(diff, leftStrategy)
		} else {
			rightStrategy := getStrategyByName(rightConfig.Strategies, leftStrategy.Strategy)
			modelDiffFound := model.Strategy{
				Strategy: leftStrategy.Strategy,
			}

			diffFound := false
			diffFound = compareAndUpdateStrategyActivated(leftStrategy, rightStrategy, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateStrategyOperation(leftStrategy, rightStrategy, diffFound, &modelDiffFound)
			diffFound = compareAndUpdateStrategyValues(leftStrategy, rightStrategy, diffFound, &modelDiffFound)

			if diffFound {
				diff = append(diff, modelDiffFound)
			}
		}
	}

	return diff
}

func compareAndUpdateStrategyActivated(leftStrategy model.Strategy,
	rightStrategy model.Strategy,
	diffFound bool,
	modelDiffFound *model.Strategy) bool {
	if *leftStrategy.Activated != *rightStrategy.Activated {
		diffFound = true
		modelDiffFound.Activated = rightStrategy.Activated
	}
	return diffFound
}

func compareAndUpdateStrategyOperation(leftStrategy model.Strategy,
	rightStrategy model.Strategy,
	diffFound bool,
	modelDiffFound *model.Strategy) bool {
	if leftStrategy.Operation != rightStrategy.Operation {
		diffFound = true
		modelDiffFound.Operation = rightStrategy.Operation
	}
	return diffFound
}

func compareAndUpdateStrategyValues(leftStrategy model.Strategy,
	rightStrategy model.Strategy,
	diffFound bool,
	modelDiffFound *model.Strategy) bool {
	if len(leftStrategy.Values) > 0 {
		modelDiffFound.Values = checkValuesDiff(leftStrategy, rightStrategy)
		if len(modelDiffFound.Values) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkValuesDiff(leftStrategy model.Strategy, rightStrategy model.Strategy) []string {
	var diff []string

	for _, leftValue := range leftStrategy.Values {
		if !containsValue(rightStrategy.Values, leftValue) {
			diff = append(diff, leftValue)
		}
	}

	return diff
}

func compareAndUpdateConfigComponents(leftConfig model.Config,
	rightConfig model.Config,
	diffFound bool,
	modelDiffFound *model.Config) bool {
	if len(leftConfig.Components) > 0 {
		modelDiffFound.Components = checkComponentsDiff(leftConfig, rightConfig)
		if len(modelDiffFound.Components) > 0 {
			diffFound = true
		}
	}
	return diffFound
}

func checkComponentsDiff(leftConfig model.Config, rightConfig model.Config) []string {
	var diff []string

	for _, leftComponent := range rightConfig.Components {
		if !containsValue(rightConfig.Components, leftComponent) {
			diff = append(diff, leftComponent)
		}
	}

	return diff
}

func containsStrategy(strategies []model.Strategy, strategy model.Strategy) bool {
	for _, s := range strategies {
		if s.Strategy == strategy.Strategy {
			return true
		}
	}
	return false
}

func getStrategyByName(strategies []model.Strategy, name string) model.Strategy {
	for _, s := range strategies {
		if s.Strategy == name {
			return s
		}
	}
	return model.Strategy{}
}

func containsValue(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func getConfigByKey(configs []model.Config, key string) model.Config {
	for _, c := range configs {
		if c.Key == key {
			return c
		}
	}
	return model.Config{}
}

func getGroupByName(groups []model.Group, name string) model.Group {
	for _, g := range groups {
		if g.Name == name {
			return g
		}
	}
	return model.Group{}
}

func containsGroup(groups []model.Group, group model.Group) bool {
	for _, g := range groups {
		if g.Name == group.Name {
			return true
		}
	}
	return false
}

func containsConfig(configs []model.Config, config model.Config) bool {
	for _, c := range configs {
		if c.Key == config.Key {
			return true
		}
	}
	return false
}
