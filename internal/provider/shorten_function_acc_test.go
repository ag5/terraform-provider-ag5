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
							knownvalue.StringExact(strings.Repeat("a", 55)+"-28165978"),
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
				ExpectError: regexp.MustCompile(`max_length must be at least 9`),
			},
		},
	})
}
