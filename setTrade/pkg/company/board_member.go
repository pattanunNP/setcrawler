package company

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BoardMember struct {
	No       int    `json:"no"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

func ParseBoardMembers(doc *goquery.Document) []BoardMember {
	var boardMembers []BoardMember

	// Handle both types of tables
	doc.Find("table[role='table']").Eq(1).Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		var boardMember BoardMember
		boardMember.No = i + 1
		boardMember.Name = strings.TrimSpace(s.Find("td").Eq(1).Text())
		boardMember.Position = strings.TrimSpace(s.Find("td").Eq(2).Text())
		boardMembers = append(boardMembers, boardMember)
	})

	doc.Find("div[data-element='element_text_editor'] table").Eq(1).Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		var boardMember BoardMember
		boardMember.No = i + 1
		boardMember.Name = strings.TrimSpace(s.Find("td").Eq(1).Text())
		boardMember.Position = strings.TrimSpace(s.Find("td").Eq(2).Text())
		boardMembers = append(boardMembers, boardMember)
	})

	if len(boardMembers) == 0 {
		log.Println("No board members found")
	}

	return boardMembers
}
