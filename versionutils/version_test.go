package versionutils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/wx-chevalier/go-utils/versionutils"
)

var _ = Describe("Version", func() {

	getVersion := func(major, minor, patch int) *versionutils.Version {
		return &versionutils.Version{
			Major: major,
			Minor: minor,
			Patch: patch,
		}
	}

	var _ = Context("matchesRegex", func() {
		It("works", func() {
			Expect(versionutils.MatchesRegex("v0.1.2")).To(BeTrue())
			Expect(versionutils.MatchesRegex("v0.0.0")).To(BeTrue())
			Expect(versionutils.MatchesRegex("v0.0.1")).To(BeTrue())
			Expect(versionutils.MatchesRegex("v1.0.0")).To(BeTrue())
			Expect(versionutils.MatchesRegex("v1.0.0-rc1")).To(BeTrue())
			Expect(versionutils.MatchesRegex("v0.5.20-rc100")).To(BeTrue())
			Expect(versionutils.MatchesRegex("v0.0.0-rc1")).To(BeTrue())
			Expect(versionutils.MatchesRegex("0.1.2")).To(BeFalse())
			Expect(versionutils.MatchesRegex("v1.2")).To(BeFalse())
			Expect(versionutils.MatchesRegex("v1.1.2.5")).To(BeFalse())
			Expect(versionutils.MatchesRegex("vX.Y.2")).To(BeFalse())
			Expect(versionutils.MatchesRegex("v1.2.3-rc")).To(BeFalse())
			Expect(versionutils.MatchesRegex("v1.2.3-rc-1")).To(BeFalse())
			Expect(versionutils.MatchesRegex("v1.2.3-release1")).To(BeFalse())
			Expect(versionutils.MatchesRegex("v1.2.3+rc1")).To(BeFalse())
		})
	})

	var _ = Context("ParseVersion", func() {
		matches := func(tag string, major, minor, patch int) bool {
			parsed, err := versionutils.ParseVersion(tag)
			return err == nil && (*parsed == *getVersion(major, minor, patch))
		}

		It("works", func() {
			Expect(matches("v0.0.0", 0, 0, 0)).To(BeTrue())
			Expect(matches("v0.1.2", 0, 1, 2)).To(BeTrue())
			Expect(matches("v0.1.2", 0, 1, 3)).To(BeFalse())
		})

		It("errors when invalid semver version provided", func() {
			parsed, err := versionutils.ParseVersion("0.1.2")
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(Equal(versionutils.InvalidSemverVersionError("0.1.2").Error()))
			Expect(parsed).To(BeNil())
		})

	})

	var _ = Context("IsGreaterThanTag", func() {

		expectResult := func(greater, lesser string, worked bool, err string) {
			actualWorked, actualErr := versionutils.IsGreaterThanTag(greater, lesser)
			Expect(actualWorked).To(BeEquivalentTo(worked))
			greaterVersion, parseGreaterErr := versionutils.ParseVersion(greater)
			lesserVersion, parseLesserErr := versionutils.ParseVersion(lesser)
			gteResult, gteError := greaterVersion.IsGreaterThanOrEqualTo(lesserVersion)
			if err == "" {
				Expect(actualErr).To(BeNil())
				Expect(gteResult).To(BeEquivalentTo(worked))
				Expect(gteError).To(BeNil())
			} else {
				Expect(actualErr.Error()).To(BeEquivalentTo(err))
			}
			if parseGreaterErr == nil && parseLesserErr == nil {
				Expect(greaterVersion.IsGreaterThanOrEqualTo(lesserVersion)).To(BeEquivalentTo(worked))
			}
			if parseGreaterErr != nil {
				Expect(parseGreaterErr.Error()).To(BeEquivalentTo(err))
				Expect(gteError.Error()).To(BeEquivalentTo("cannot compare versions, greater version is nil"))
			}
			if parseLesserErr != nil {
				Expect(parseLesserErr.Error()).To(BeEquivalentTo(err))
				Expect(gteError.Error()).To(BeEquivalentTo("cannot compare versions, lesser version is nil"))
			}
		}

		It("works", func() {
			expectResult("v0.1.2", "v0.0.1", true, "")
			expectResult("v0.0.1", "v0.1.2", false, "")
			expectResult("v0.0.1", "v0.0.0", true, "")
			expectResult("0.0.2", "v0.0.1", false, versionutils.InvalidSemverVersionError("0.0.2").Error())
			expectResult("v0.0.2", "0.0.1", false, versionutils.InvalidSemverVersionError("0.0.1").Error())
			expectResult("v1.0.0", "v0.0.1-rc1", true, "")
			expectResult("v1.0.0-rc1", "v1.0.0-rc2", false, "")
			expectResult("v1.0.0-rc2", "v1.0.0-rc1", true, "")
			expectResult("v1.0.0-rc1", "v1.0.0", false, "")
			expectResult("v1.0.0", "v1.0.0-rc1", true, "")
		})
	})

	var _ = Context("GetVersionFromTag", func() {

		It("works", func() {
			Expect(versionutils.GetVersionFromTag("v0.1.2")).To(Equal("0.1.2"))
			Expect(versionutils.GetVersionFromTag("0.1.2")).To(Equal("0.1.2"))
		})
	})

	var _ = Context("IncrementVersion", func() {

		expectResult := func(start *versionutils.Version, breakingChange bool, newFeature bool, expected *versionutils.Version) {
			actualIncremented := start.IncrementVersion(breakingChange, newFeature)
			Expect(actualIncremented).To(BeEquivalentTo(expected))
		}

		It("works", func() {
			expectResult(getVersion(0, 0, 1), true, true, getVersion(0, 1, 0))
			expectResult(getVersion(0, 0, 1), true, false, getVersion(0, 1, 0))
			expectResult(getVersion(0, 1, 10), true, false, getVersion(0, 2, 0))
			expectResult(getVersion(0, 1, 10), true, true, getVersion(0, 2, 0))
			expectResult(getVersion(1, 1, 10), true, false, getVersion(2, 0, 0))
			expectResult(getVersion(1, 1, 10), true, true, getVersion(2, 0, 0))
			expectResult(getVersion(0, 0, 1), false, false, getVersion(0, 0, 2))
			expectResult(getVersion(0, 0, 1), false, true, getVersion(0, 0, 2))
			expectResult(getVersion(0, 1, 10), false, false, getVersion(0, 1, 11))
			expectResult(getVersion(0, 1, 10), false, true, getVersion(0, 1, 11))
			expectResult(getVersion(1, 1, 10), false, false, getVersion(1, 1, 11))
			expectResult(getVersion(1, 1, 10), false, true, getVersion(1, 2, 0))
		})
	})

})
