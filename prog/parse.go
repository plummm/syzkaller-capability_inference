// Copyright 2015 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package prog

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/google/syzkaller/pkg/log"
)

// LogEntry describes one program in execution log.
type LogEntry struct {
	P     *Prog
	Proc  int // index of parallel proc
	Start int // start offset in log
	End   int // end offset in log
	Index int
}

func (target *Target) ParseLog(data []byte) []*LogEntry {
	var entries []*LogEntry
	ent := &LogEntry{}
	var cur []byte
	faultCall, faultNth := -1, -1
	callIndex := 0
	for pos := 0; pos < len(data); {
		nl := bytes.IndexByte(data[pos:], '\n')
		if nl == -1 {
			nl = len(data) - 1
		} else {
			nl += pos
		}
		line := data[pos : nl+1]
		pos0 := pos
		pos = nl + 1

		if strings.HasPrefix(string(line), "MAGIC?!CALL") {
			callIndex += 1
		}

		if proc, ok := extractInt(line, "executing program "); ok {
			ent.Index = callIndex
			callIndex = 0
			if ent.P != nil && len(ent.P.Calls) != 0 {
				ent.End = pos0
				entries = append(entries, ent)
				faultCall, faultNth = -1, -1
			}
			ent = &LogEntry{
				Proc:  proc,
				Start: pos0,
			}
			// We no longer print it this way, but we still parse such fragments to preserve
			// the backward compatibility.
			if parsedFaultCall, ok := extractInt(line, "fault-call:"); ok {
				faultCall = parsedFaultCall
				faultNth, _ = extractInt(line, "fault-nth:")
			}
			cur = nil
			continue
		}
		if ent == nil {
			continue
		}
		tmp := append(cur, line...)

		p, err := target.Deserialize(tmp, NonStrict)
		if err != nil {
			continue
		}

		if faultCall >= 0 && faultCall < len(p.Calls) {
			// We add 1 because now the property is 1-based.
			p.Calls[faultCall].Props.FailNth = faultNth + 1
		}

		cur = tmp
		ent.P = p
	}
	if ent.P != nil && len(ent.P.Calls) != 0 {
		ent.Index = callIndex
		ent.End = len(data)
		entries = append(entries, ent)
	}

	for i, ent := range entries {
		log.Logf(1, "ent[%v] %v: %v", i, ent.Index, string(ent.P.Serialize()))
	}
	return entries
}

func extractInt(line []byte, prefix string) (int, bool) {
	pos := bytes.Index(line, []byte(prefix))
	if pos == -1 {
		return 0, false
	}
	pos += len(prefix)
	end := pos
	for end != len(line) && line[end] >= '0' && line[end] <= '9' {
		end++
	}
	v, _ := strconv.Atoi(string(line[pos:end]))
	return v, true
}
