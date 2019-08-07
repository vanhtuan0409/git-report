package gitreport

import (
	"errors"
	"fmt"
)

type Parser struct {
	JsonStr string
	Head    int
	Len     int
}

func strStr(haystack string, needle string) int {
	if len(needle) == 0 {
		return 0
	}
	nlen := len(needle)
	for i := 0; i <= len(haystack)-nlen; i++ {
		if haystack[i:i+nlen] == needle {
			return i
		}
	}
	return -1
}

func (p *Parser) nextToken(needle string) int {
	if p.Head > p.Len {
		return -1
	}
	return strStr((p.JsonStr)[p.Head:], needle)
}

func (p *Parser) initParser(data []byte) {
	p.Head = 0
	p.JsonStr = string(data)
	p.Len = len(data)

}

func parserUser(str string) (*User, error) {
	name_idx := strStr(str, "\"email\"")
	if name_idx == -1 {
		return nil, errors.New("Not found email")
	}
	name := str[12 : name_idx-2]

	email_idx := strStr(str[name_idx+7:], "\"date\"")
	if email_idx == -1 {
		return nil, errors.New("Not found date")
	}
	email_idx += name_idx
	email := str[name_idx+10 : email_idx]
	date := str[email_idx+16 : len(str)-3]
	//fmt.Printf("Name:%s\nEmail:%s\nDate:%s\n", name, email, date)

	return &User{
		Name:  name,
		Email: email,
		Date:  date,
	}, nil
}

func (p *Parser) parseStr(from int, to int) (*GitCommit, error) {
	str := p.JsonStr[from:to]
	commitStr_idx := strStr(str, "\"refs\"")
	if commitStr_idx == -1 {
		return nil, errors.New("Not found refs")
	}

	commitStr := str[9 : commitStr_idx-2]

	refs_idx := strStr(str[commitStr_idx:], "\"subject\"")
	if refs_idx == -1 {
		return nil, errors.New("Not found subject")
	}
	refs_idx += commitStr_idx
	refsStr := str[commitStr_idx+8 : refs_idx-3]

	subject_idx := strStr(str[refs_idx:], "\"body\"")
	if subject_idx == -1 {
		return nil, errors.New("Not found body")
	}
	subject_idx += refs_idx
	subjectStr := str[refs_idx+12 : subject_idx-2]

	body_idx := strStr(str[subject_idx:], "\"author\"")
	if body_idx == -1 {
		return nil, errors.New("Not found author")
	}
	body_idx += subject_idx
	bodyStr := str[subject_idx+9 : body_idx-2]

	author_idx := strStr(str[body_idx:], "\"commiter\"")
	if author_idx == -1 {
		return nil, errors.New("Not found commiter")
	}
	author_idx += body_idx

	//fmt.Printf("commit:%s\nrefs:%s\nsubject:%s\nbody:%s\n", commitStr, refsStr, subjectStr, bodyStr)
	u, err := parserUser(str[body_idx+8 : author_idx])
	if err != nil {
		return nil, err
	}

	commiter, err := parserUser(str[author_idx+10:])
	if err != nil {
		return nil, err
	}

	return &GitCommit{
		Hash:      commitStr,
		Refs:      refsStr,
		Subject:   subjectStr,
		Body:      bodyStr,
		Author:    u,
		Committer: commiter,
	}, nil
}

func (p *Parser) Unmarshal(data []byte) ([]*GitCommit, error) {
	var commits = make([]*GitCommit, 0)
	p.initParser(data)

	from_idx := p.nextToken("\"commit\"")
	for {
		p.Head += 8
		to_idx := p.nextToken("\"commit\"")
		if to_idx != -1 {
			to_idx += p.Head
			commit, err := p.parseStr(from_idx, to_idx)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			commits = append(commits, commit)
			from_idx = to_idx
			p.Head = to_idx + 8
		} else { // end of parser laster one
			commit, err := p.parseStr(from_idx, p.Len-1)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			commits = append(commits, commit)
			break
		}

	}
	return commits, nil
}
