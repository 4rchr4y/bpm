package bundle

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/4rchr4y/bpm/constant"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/go-version"
)

const (
	VersionDateFormat   = "20060102150405"
	VersionShortHashLen = 12
)

const (
	versionRegexStr = `^(v\d+\.\d+\.\d+)\+(\d{14})-(\w+)$`
)

var (
	versionRegex = regexp.MustCompile(`^(v\d+\.\d+\.\d+)\+(\d{14})-(\w+)$`)
)

type VersionExpr struct {
	// TODO: create Tag type
	Tag       *version.Version // semantic tag if available, or pseudo version
	Timestamp string           // commit timestamp
	Hash      string           // commit hash
}

func NewVersionExpr(commit *object.Commit, tag *version.Version) *VersionExpr {
	return &VersionExpr{
		Tag:       tag,
		Timestamp: commit.Committer.When.UTC().Format(VersionDateFormat),
		Hash:      commit.Hash.String()[:VersionShortHashLen],
	}
}

func (v *VersionExpr) Major() int { return v.Tag.Segments()[0] }
func (v *VersionExpr) Minor() int { return v.Tag.Segments()[1] }
func (v *VersionExpr) Path() int  { return v.Tag.Segments()[2] }

func (v *VersionExpr) Equal(o *VersionExpr) bool       { return v.Tag.Equal(o.Tag) }
func (v *VersionExpr) GreaterThan(o *VersionExpr) bool { return v.Tag.GreaterThan(o.Tag) }

func (v *VersionExpr) IsPseudo() bool {
	return v.Tag.Original() == constant.BundlePseudoVersion && v.Hash != "" && v.Timestamp != ""
}

func (v *VersionExpr) String() string {
	if v == nil {
		return "latest"
	}

	if v.Tag != nil && v.Tag.Original() != constant.BundlePseudoVersion {
		return v.Tag.Original()
	}

	return fmt.Sprintf("%s+%s-%s", constant.BundlePseudoVersion, v.Timestamp, v.Hash)
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

		return &VersionExpr{Tag: v}, nil

	default:
		matches := versionRegex.FindStringSubmatch(versionStr)
		if matches == nil || len(matches) != 4 {
			return nil, fmt.Errorf("invalid version format")
		}

		v, err := version.NewVersion(matches[1])
		if err != nil {
			return nil, err
		}

		return &VersionExpr{
			Tag:       v,
			Timestamp: matches[2],
			Hash:      matches[3],
		}, nil
	}
}
