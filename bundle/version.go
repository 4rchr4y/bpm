package bundle

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/4rchr4y/godevkit/v3/must"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
)

const (
	PseudoSemTagStr     = "v0.0.0"
	versionLatestStr    = "latest"
	VersionDateFormat   = "20060102150405"
	VersionShortHashLen = 12
	VersionRegexStr     = `^(v\d+\.\d+\.\d+)\+(\d{14})-(\w+)$`
)

var (
	PseudoSemTag = must.Must(version.NewSemver(PseudoSemTagStr))
	VersionRegex = regexp.MustCompile(VersionRegexStr)
)

type VersionSpec struct {
	SemTag    *version.Version // semantic tag if available, or pseudo semantic tag
	Timestamp time.Time        // commit timestamp
	Hash      string           // commit hash
}

func NewVersionSpecFromCommit(commit *object.Commit, tag *version.Version) *VersionSpec {
	return &VersionSpec{
		SemTag:    tag,
		Timestamp: commit.Committer.When.UTC(),
		Hash:      commit.Hash.String()[:VersionShortHashLen],
	}
}

func (v *VersionSpec) IsPseudo() bool {
	return v.SemTag != nil &&
		v.SemTag.Original() == PseudoSemTagStr &&
		v.Hash != "" &&
		!v.Timestamp.IsZero()
}

func (v *VersionSpec) Major() int { return v.SemTag.Segments()[0] }
func (v *VersionSpec) Minor() int { return v.SemTag.Segments()[1] }
func (v *VersionSpec) Path() int  { return v.SemTag.Segments()[2] }

func (v *VersionSpec) Equal(o *VersionSpec) bool {
	if v.IsPseudo() && o.IsPseudo() {
		return v.String() == o.String()
	}

	return v.SemTag.Equal(o.SemTag)
}

func (v *VersionSpec) GreaterThan(o *VersionSpec) bool {
	if v.IsPseudo() && o.IsPseudo() {
		return v.Timestamp.After(o.Timestamp)
	}

	return v.SemTag.GreaterThan(o.SemTag)
}

func (v *VersionSpec) String() string {
	if v == nil {
		return versionLatestStr
	}

	if v.SemTag != nil && v.SemTag.Original() != PseudoSemTagStr {
		return v.SemTag.Original()
	}

	return fmt.Sprintf("%s+%s-%s",
		PseudoSemTagStr,
		v.Timestamp.Format(VersionDateFormat),
		v.Hash,
	)
}

func ParseVersionExpr(versionStr string) (*VersionSpec, error) {
	switch {
	case versionStr == "":
		return nil, nil

	case !strings.Contains(versionStr, "+"):
		v, err := version.NewVersion(versionStr)
		if err != nil {
			return nil, err
		}

		return &VersionSpec{SemTag: v}, nil

	default:
		matches := VersionRegex.FindStringSubmatch(versionStr)
		if matches == nil || len(matches) != 4 {
			return nil, fmt.Errorf("invalid version format")
		}

		v, err := version.NewVersion(matches[1])
		if err != nil {
			return nil, err
		}

		timestamp, err := time.Parse(VersionDateFormat, matches[2])
		if err != nil {
			return nil, err
		}

		return &VersionSpec{
			SemTag:    v,
			Timestamp: timestamp,
			Hash:      matches[3],
		}, nil
	}
}
