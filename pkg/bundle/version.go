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
	VersionRegex = regexp.MustCompile(`^(v\d+\.\d+\.\d+)\+(\d{14})-(\w+)$`)
)

type VersionExpr struct {
	SemTag    *version.Version // semantic tag if available, or pseudo semantic tag
	Timestamp time.Time        // commit timestamp
	Hash      string           // commit hash
}

func NewVersionExprFromCommit(commit *object.Commit, tag *version.Version) *VersionExpr {
	return &VersionExpr{
		SemTag:    tag,
		Timestamp: commit.Committer.When.UTC(),
		Hash:      commit.Hash.String()[:VersionShortHashLen],
	}
}

func (v *VersionExpr) IsPseudo() bool {
	return v.SemTag.Original() == PseudoSemTagStr && v.Hash != "" && !v.Timestamp.IsZero()
}

func (v *VersionExpr) Major() int { return v.SemTag.Segments()[0] }
func (v *VersionExpr) Minor() int { return v.SemTag.Segments()[1] }
func (v *VersionExpr) Path() int  { return v.SemTag.Segments()[2] }

func (v *VersionExpr) Equal(o *VersionExpr) bool {
	if v.IsPseudo() && o.IsPseudo() {
		return v.String() == o.String()
	}

	return v.SemTag.Equal(o.SemTag)
}

func (v *VersionExpr) GreaterThan(o *VersionExpr) bool {
	if v.IsPseudo() && o.IsPseudo() {
		return v.Timestamp.After(o.Timestamp)
	}

	return v.SemTag.GreaterThan(o.SemTag)
}

func (v *VersionExpr) String() string {
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

func ParseVersionExpr(versionStr string) (*VersionExpr, error) {
	switch {
	case versionStr == "":
		return nil, nil

	case !strings.Contains(versionStr, "+"):
		v, err := version.NewVersion(versionStr)
		if err != nil {
			return nil, err
		}

		return &VersionExpr{SemTag: v}, nil

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

		return &VersionExpr{
			SemTag:    v,
			Timestamp: timestamp,
			Hash:      matches[3],
		}, nil
	}
}

func FormatVersion(source string, version string) string {
	return source + "@" + version
}

func FormatVersionFromBundle(b *Bundle) string {
	return FormatVersion(b.Repository(), b.Version.String())
}
