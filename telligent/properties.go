package telligent // import "github.com/BenLubar/wtdwtf-science/telligent"

import (
	"database/sql"
	"strconv"
	"strings"
)

type properties struct {
	keys    sql.NullString
	strings sql.NullString
	bytes   []byte
	parsed  map[string]sql.NullString
}

func (p *properties) parse() {
	if p.parsed != nil {
		return
	}

	p.parsed = make(map[string]sql.NullString)
	if !p.keys.Valid || p.keys.String == "" {
		return
	}
	keys := strings.Split(p.keys.String, ":")
	if len(keys)%4 != 1 || keys[len(keys)-1] != "" {
		panic("unexpected key array length: " + strconv.Itoa(len(keys)) + "\n" + strconv.Quote(p.keys.String))
	}
	keys = keys[:len(keys)-1]
	values := []rune(p.strings.String)

	for i := 0; i < len(keys); i += 4 {
		start, err := strconv.Atoi(keys[i+2])
		if err != nil {
			panic(err)
		}
		length, err := strconv.Atoi(keys[i+3])
		if err != nil {
			panic(err)
		}
		if length == -1 {
			p.parsed[keys[i]] = sql.NullString{Valid: false}
			continue
		}
		switch keys[i+1] {
		case "S":
			p.parsed[keys[i]] = sql.NullString{Valid: true, String: string(values[start : start+length])}
		case "B":
			p.parsed[keys[i]] = sql.NullString{Valid: true, String: string(p.bytes[start : start+length])}
		default:
			panic("unexpected key type " + strconv.Quote(keys[i+1]))
		}
	}
}

func (p *properties) String(key string) (s sql.NullString, ok bool) {
	p.parse()
	s, ok = p.parsed[key]
	return
}

func (p *properties) StringOrDefault(key, def string) string {
	s, ok := p.String(key)
	if !ok || !s.Valid {
		return def
	}
	return s.String
}

func (p *properties) Bytes(key string) (b []byte, ok bool) {
	p.parse()
	s, ok := p.parsed[key]
	if s.Valid {
		b = []byte(s.String)
	}
	return
}

func (p *properties) BytesOrDefault(key string, def []byte) []byte {
	b, ok := p.Bytes(key)
	if !ok || b == nil {
		return def
	}
	return b
}
