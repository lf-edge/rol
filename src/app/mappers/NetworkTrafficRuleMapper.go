// Package mappers uses for entity <--> dto conversions
package mappers

import (
	"github.com/coreos/go-iptables/iptables"
	"rol/domain"
	"rol/dtos"
)

//MapStatToTrafficRule map iptables.Stat to domain.HostNetworkTrafficRule
func MapStatToTrafficRule(stat iptables.Stat, rule *domain.HostNetworkTrafficRule) {
	rule.Source = stat.Source.String()
	rule.Destination = stat.Destination.String()
	rule.Action = stat.Target
}

//MapHostNetworkTrafficRuleEntityToDto map HostNetworkTrafficRule entity to dto
func MapHostNetworkTrafficRuleEntityToDto(entity domain.HostNetworkTrafficRule, dto *dtos.HostNetworkTrafficRuleDto) {
	dto.Chain = entity.Chain
	dto.Action = entity.Action
	dto.Source = entity.Source
	dto.Destination = entity.Destination
}

//MapHostNetworkTrafficRuleCreateDtoToEntity map HostNetworkTrafficRuleCreateDto dto to entity
func MapHostNetworkTrafficRuleCreateDtoToEntity(dto dtos.HostNetworkTrafficRuleCreateDto, entity *domain.HostNetworkTrafficRule) {
	entity.Chain = dto.Chain
	entity.Action = dto.Action
	entity.Source = dto.Source
	entity.Destination = dto.Destination
}

//MapHostNetworkTrafficRuleDeleteDtoToEntity map HostNetworkTrafficRuleDeleteDto dto to entity
func MapHostNetworkTrafficRuleDeleteDtoToEntity(dto dtos.HostNetworkTrafficRuleDeleteDto, entity *domain.HostNetworkTrafficRule) {
	entity.Chain = dto.Chain
	entity.Action = dto.Action
	entity.Source = dto.Source
	entity.Destination = dto.Destination
}
