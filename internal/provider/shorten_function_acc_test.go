package provider

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestShortenFunction(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
output "test" {
  value = provider::ag5::shorten("` + strings.Repeat("a", 100) + `", 64)
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"test",
							knownvalue.StringExact(strings.Repeat("a", 58)+"-28165"),
						),
					},
				},
			},
			{
				Config: `
output "test" {
  value = provider::ag5::shorten("` + strings.Repeat("a", 100) + `", 64, 12)
}
`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue(
							"test",
							knownvalue.StringExact(strings.Repeat("a", 51)+"-2816597888e4"),
						),
					},
				},
			},
			{
				Config: `
output "test" {
  value = provider::ag5::shorten("anything", 8)
}
`,
				ExpectError: errorPattern("max_length must be at least 9"),
			},
			{
				Config: `
output "test" {
  value = provider::ag5::shorten("anything", 64, 0)
}
`,
				ExpectError: errorPattern("hash_length must be between 1 and 64"),
			},
			{
				Config: `
output "test" {
  value = provider::ag5::shorten("anything", 64, 65)
}
`,
				ExpectError: errorPattern("hash_length must be between 1 and 64"),
			},
			{
				Config: `
output "test" {
  value = provider::ag5::shorten("anything", 64, 5, 8)
}
`,
				ExpectError: errorPattern("hash_length accepts at most one value"),
			},
		},
	})
}

// errorPattern matches message text regardless of where Terraform wraps it.
func errorPattern(message string) *regexp.Regexp {
	return regexp.MustCompile(strings.ReplaceAll(regexp.QuoteMeta(message), " ", `\s+`))
}
