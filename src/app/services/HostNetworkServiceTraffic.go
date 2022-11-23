// Package services stores business logic for each entity
package services

import (
	"rol/app/errors"
	"rol/app/mappers"
	"rol/domain"
	"rol/dtos"
)

//CreateTrafficRule Create netfilter traffic rule for specified table
//
//Params:
//	table - table to create a rule
//	rule - rule entity
//Return:
//	domain.NetworkTrafficRule - new traffic rule
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) CreateTrafficRule(table string, ruleDto dtos.HostNetworkTrafficRuleCreateDto) (dtos.HostNetworkTrafficRuleDto, error) {
	var rule domain.HostNetworkTrafficRule
	var out dtos.HostNetworkTrafficRuleDto
	mappers.MapHostNetworkTrafficRuleCreateDtoToEntity(ruleDto, &rule)
	newRule, err := h.manager.CreateTrafficRule(table, rule)
	if err != nil {
		return dtos.HostNetworkTrafficRuleDto{}, errors.Internal.Wrap(err, "host network manager failed to create traffic rule")
	}
	mappers.MapHostNetworkTrafficRuleEntityToDto(newRule, &out)
	return out, nil
}

//DeleteTrafficRule Delete netfilter traffic rule in specified table
//
//Params:
//	table - table to delete a rule
//	rule - rule entity
//Return:
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) DeleteTrafficRule(table string, ruleDto dtos.HostNetworkTrafficRuleDeleteDto) error {
	var rule domain.HostNetworkTrafficRule
	mappers.MapHostNetworkTrafficRuleDeleteDtoToEntity(ruleDto, &rule)
	return h.manager.DeleteTrafficRule(table, rule)
}

//GetChainRules Get selected netfilter chain rules at specified table
//
//Params:
//	table - table to get a rules
//	chain - chain where we get the rules
//Return:
//	[]domain.NetworkTrafficRule - slice of rules
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) GetChainRules(table string, chain string) ([]dtos.HostNetworkTrafficRuleDto, error) {
	var out []dtos.HostNetworkTrafficRuleDto
	rules, err := h.manager.GetChainRules(table, chain)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "host network manager failed to get chain rule")
	}
	for _, rule := range rules {
		var dto dtos.HostNetworkTrafficRuleDto
		mappers.MapHostNetworkTrafficRuleEntityToDto(rule, &dto)
		out = append(out, dto)
	}
	return out, nil
}

//GetTableRules Get specified netfilter table rules
//
//Params:
//	table - table to get a rules
//Return:
//	[]domain.NetworkTrafficRule - slice of rules
//	error - if an error occurs, otherwise nil
func (h *HostNetworkService) GetTableRules(table string) ([]dtos.HostNetworkTrafficRuleDto, error) {
	var out []dtos.HostNetworkTrafficRuleDto
	rules, err := h.manager.GetTableRules(table)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "host network manager failed to get table rule")
	}
	for _, rule := range rules {
		var dto dtos.HostNetworkTrafficRuleDto
		mappers.MapHostNetworkTrafficRuleEntityToDto(rule, &dto)
		out = append(out, dto)
	}
	return out, nil
}
