package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// ReferencesMarker marks where references should be placed.
const ReferencesMarker = "// References:"

// citations is the regular expression for citations in comments.
var citations = regexp.MustCompile(`\[[[:alnum:]]{4,}\]`)

// Source represents a parsed source file with references.
type Source struct {
	Lines     []string
	InsertAt  int
	Citations map[string]bool
}

// Parse a source file.
func Parse(r io.Reader) (*Source, error) {
	s := &Source{
		InsertAt:  -1,
		Citations: map[string]bool{},
	}

	scanner := bufio.NewScanner(r)
	insideReferenceBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		// Is this the start of the reference block?
		if line == ReferencesMarker && s.InsertAt < 0 {
			s.InsertAt = len(s.Lines)
			insideReferenceBlock = true
			continue
		}

		if IsComment(line) {
			// Look for citations.
			keys := ParseCitations(line)
			for _, key := range keys {
				s.Citations[key] = true
			}
		} else if insideReferenceBlock {
			insideReferenceBlock = false
		}

		// Record the line.
		if !insideReferenceBlock {
			s.Lines = append(s.Lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return s, nil
}

// ParseFile parses a source file for citations and references.
func ParseFile(path string) (*Source, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s, err := Parse(f)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// IsComment returns whether the line is a comment.
func IsComment(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "//")
}

// ParseCitations parses citations from a line.
func ParseCitations(line string) []string {
	keys := []string{}
	matches := citations.FindAllString(line, -1)
	for _, match := range matches {
		keys = append(keys, match[1:len(match)-1])
	}
	return keys
}

// Validate the citations in the source.
func (s *Source) Validate(b *Bibliography) error {
	for key := range s.Citations {
		if b.Lookup(key) == nil {
			return fmt.Errorf("unknown reference '%s'", key)
		}
	}
	return nil
}

func (s *Source) Write(w io.Writer, b *Bibliography) error {
	for i, line := range s.Lines {
		// Write the reference block if we're at the insertion point.
		if i == s.InsertAt {
			if err := s.writeReferences(w, b); err != nil {
				return err
			}
		}

		// Write this line.
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func (s *Source) writeReferences(w io.Writer, b *Bibliography) error {
	fmt.Fprintf(w, "%s\n//\n", ReferencesMarker)
	for key := range s.Citations {
		_, err := fmt.Fprintf(w, "// %s\n", key)
		if err != nil {
			return err
		}
	}
	return nil
}
