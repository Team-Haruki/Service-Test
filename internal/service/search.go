package service

import (
	"fmt"
	"log/slog"

	"Haruki-Service-API/pkg/masterdata"
)

// CardSearchService 负责卡牌搜索逻辑。
type CardSearchService struct {
	repo   CardDataSource
	parser *CardParser
}

// NewCardSearchService 创建卡牌搜索服务。
func NewCardSearchService(repo CardDataSource, parser *CardParser) *CardSearchService {
	return &CardSearchService{
		repo:   repo,
		parser: parser,
	}
}

func (s *CardSearchService) CloneWithRepo(repo CardDataSource) *CardSearchService {
	if s == nil {
		return nil
	}
	return &CardSearchService{
		repo:   repo,
		parser: s.parser,
	}
}

// Search 根据查询字符串搜索单张卡牌。
func (s *CardSearchService) Search(query string) (*masterdata.Card, error) {
	info, err := s.parser.Parse(query)
	if err == nil && info != nil {
		switch info.Type {
		case QueryTypeID:
			card, getErr := s.repo.GetCardByID(info.Value)
			if getErr != nil {
				return nil, getErr
			}
			return card, nil
		case QueryTypeSeq:
			return s.repo.GetCardByCharacterAndSeq(info.CharacterID, info.Sequence)
		case QueryTypeFilter:
			filtered, filterErr := s.repo.FilterCards(info)
			if filterErr != nil {
				return nil, filterErr
			}
			if len(filtered) == 0 {
				return nil, fmt.Errorf("card not found (filter): %s", query)
			}
			return filtered[len(filtered)-1], nil
		}
	}

	slog.Debug("failed to parse card search query", "query", query)
	return nil, fmt.Errorf("无法解析的指令: %s", query)
}

// SearchList 根据查询字符串搜索卡牌列表。
func (s *CardSearchService) SearchList(query string) ([]*masterdata.Card, error) {
	info, err := s.parser.Parse(query)
	if err == nil && info != nil {
		switch info.Type {
		case QueryTypeFilter:
			filtered, filterErr := s.repo.FilterCards(info)
			if filterErr != nil {
				return nil, filterErr
			}
			if len(filtered) == 0 {
				return nil, fmt.Errorf("no cards found for filter: %s", query)
			}
			return filtered, nil
		case QueryTypeID:
			card, getErr := s.repo.GetCardByID(info.Value)
			if getErr != nil {
				return nil, getErr
			}
			return []*masterdata.Card{card}, nil
		case QueryTypeSeq:
			card, getErr := s.repo.GetCardByCharacterAndSeq(info.CharacterID, info.Sequence)
			if getErr != nil {
				return nil, getErr
			}
			return []*masterdata.Card{card}, nil
		}
	}

	slog.Debug("failed to parse card list query", "query", query)
	return nil, fmt.Errorf("无法解析的列表查询指令: %s", query)
}
