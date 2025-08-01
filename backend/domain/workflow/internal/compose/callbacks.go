/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package compose

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/selector"
)

type selectorCallbackField struct {
	Key   string      `json:"key"`
	Type  vo.DataType `json:"type"`
	Value any         `json:"value"`
}

type selectorCondition struct {
	Left     selectorCallbackField  `json:"left"`
	Operator vo.OperatorType        `json:"operator"`
	Right    *selectorCallbackField `json:"right,omitempty"`
}

type selectorBranch struct {
	Conditions []*selectorCondition `json:"conditions"`
	Logic      vo.LogicType         `json:"logic"`
	Name       string               `json:"name"`
}

func (s *NodeSchema) toSelectorCallbackInput(sc *WorkflowSchema) func(_ context.Context, in map[string]any) (map[string]any, error) {
	return func(_ context.Context, in map[string]any) (map[string]any, error) {
		config := mustGetKey[[]*selector.OneClauseSchema]("Clauses", s.Configs)
		count := len(config)

		output := make([]*selectorBranch, count)

		for _, source := range s.InputSources {
			targetPath := source.Path
			if len(targetPath) == 2 {
				indexStr := targetPath[0]
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					return nil, err
				}

				branch := output[index]
				if branch == nil {
					output[index] = &selectorBranch{
						Conditions: []*selectorCondition{
							{
								Operator: config[index].Single.ToCanvasOperatorType(),
							},
						},
						Logic: selector.ClauseRelationAND.ToVOLogicType(),
					}
				}

				if targetPath[1] == selector.LeftKey {
					leftV, ok := nodes.TakeMapValue(in, targetPath)
					if !ok {
						return nil, fmt.Errorf("failed to take left value of %s", targetPath)
					}
					if source.Source.Ref.VariableType != nil {
						if *source.Source.Ref.VariableType == vo.ParentIntermediate {
							parentNodeKey, ok := sc.Hierarchy[s.Key]
							if !ok {
								return nil, fmt.Errorf("failed to find parent node key of %s", s.Key)
							}
							parentNode := sc.GetNode(parentNodeKey)
							output[index].Conditions[0].Left = selectorCallbackField{
								Key:   parentNode.Name + "." + strings.Join(source.Source.Ref.FromPath, "."),
								Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Type,
								Value: leftV,
							}
						} else {
							output[index].Conditions[0].Left = selectorCallbackField{
								Key:   "",
								Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Type,
								Value: leftV,
							}
						}
					} else {
						output[index].Conditions[0].Left = selectorCallbackField{
							Key:   sc.GetNode(source.Source.Ref.FromNodeKey).Name + "." + strings.Join(source.Source.Ref.FromPath, "."),
							Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Type,
							Value: leftV,
						}
					}
				} else if targetPath[1] == selector.RightKey {
					rightV, ok := nodes.TakeMapValue(in, targetPath)
					if !ok {
						return nil, fmt.Errorf("failed to take right value of %s", targetPath)
					}
					output[index].Conditions[0].Right = &selectorCallbackField{
						Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Type,
						Value: rightV,
					}
				}
			} else if len(targetPath) == 3 {
				indexStr := targetPath[0]
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					return nil, err
				}

				multi := config[index].Multi

				branch := output[index]
				if branch == nil {
					output[index] = &selectorBranch{
						Conditions: make([]*selectorCondition, len(multi.Clauses)),
						Logic:      multi.Relation.ToVOLogicType(),
					}
				}

				clauseIndexStr := targetPath[1]
				clauseIndex, err := strconv.Atoi(clauseIndexStr)
				if err != nil {
					return nil, err
				}

				clause := multi.Clauses[clauseIndex]

				if output[index].Conditions[clauseIndex] == nil {
					output[index].Conditions[clauseIndex] = &selectorCondition{
						Operator: clause.ToCanvasOperatorType(),
					}
				}

				if targetPath[2] == selector.LeftKey {
					leftV, ok := nodes.TakeMapValue(in, targetPath)
					if !ok {
						return nil, fmt.Errorf("failed to take left value of %s", targetPath)
					}
					if source.Source.Ref.VariableType != nil {
						if *source.Source.Ref.VariableType == vo.ParentIntermediate {
							parentNodeKey, ok := sc.Hierarchy[s.Key]
							if !ok {
								return nil, fmt.Errorf("failed to find parent node key of %s", s.Key)
							}
							parentNode := sc.GetNode(parentNodeKey)
							output[index].Conditions[clauseIndex].Left = selectorCallbackField{
								Key:   parentNode.Name + "." + strings.Join(source.Source.Ref.FromPath, "."),
								Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Properties[targetPath[2]].Type,
								Value: leftV,
							}
						} else {
							output[index].Conditions[clauseIndex].Left = selectorCallbackField{
								Key:   "",
								Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Properties[targetPath[2]].Type,
								Value: leftV,
							}
						}
					} else {
						output[index].Conditions[clauseIndex].Left = selectorCallbackField{
							Key:   sc.GetNode(source.Source.Ref.FromNodeKey).Name + "." + strings.Join(source.Source.Ref.FromPath, "."),
							Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Properties[targetPath[2]].Type,
							Value: leftV,
						}
					}
				} else if targetPath[2] == selector.RightKey {
					rightV, ok := nodes.TakeMapValue(in, targetPath)
					if !ok {
						return nil, fmt.Errorf("failed to take right value of %s", targetPath)
					}
					output[index].Conditions[clauseIndex].Right = &selectorCallbackField{
						Type:  s.InputTypes[targetPath[0]].Properties[targetPath[1]].Properties[targetPath[2]].Type,
						Value: rightV,
					}
				}
			}
		}

		return map[string]any{"branches": output}, nil
	}
}
