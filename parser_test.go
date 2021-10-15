package tfparser

import (
	"testing"
)

var testTFCode = `
# use1.genhtcc.com <-> aps2.genhtcc.com
module "legacy_use1_legacy_aps2_routing" {
  source                 = "../../modules/vpc-routing"
  alice_vpc_name         = "Development VPC"
	//commented_param         = value
  bob_vpc_name           = "aps2.g1.genhtcc.com"

  providers = {
    aws.alice = aws.us-east-1
    aws.bob   = aws.ap-southeast-2
  }
}
/*

multiline comment

*/

# use1.genhtcc.com <-> aps2.g1.genhtcc.com
module "legacy_use1_g1_aps2_routing" {
  source                 = "../../modules/vpc-routing"
  alice_region           = "us-east-1"
  alice_vpc_name         = "Development VPC"
  bob_vpc_name           = "aps2.g1.genhtcc.com"

  providers = {
    aws.alice = aws.us-east-1
    aws.bob   = aws.ap-southeast-2
  }
}

`
var config, parseErr = ParseString(testTFCode)

func TestParser1(t *testing.T) {
	if parseErr != nil {
		t.Fatalf("parser returned error %v", parseErr)
	}
}
func TestModulesLength(t *testing.T) {
	if len(config.Modules) != 2 {
		t.Fatalf("Unexpceted number of modules parsed: %v", len(config.Modules))
	}
}

func TestModule1ParamsNo(t *testing.T) {
	m, exists := config.Modules["legacy_use1_legacy_aps2_routing"]
	// test case when we check that it must exist is next tst
	if exists {
		if l := len(m.Parameters); l != 2 {
			t.Fatalf("Unexpected number of parameters found for Module 'legacy_use1_legacy_aps2_routing', %v instead of 2", l)
		}
	}
}

func TestModule1Params(t *testing.T) {
	m, exists := config.Modules["legacy_use1_legacy_aps2_routing"]
	if !exists {
		t.Fatal("Module 'legacy_use1_legacy_aps2_routing' was not found")
	}
	al, exists := m.Providers["aws.alice"]
	if !exists {
		t.Fatal("Module 'legacy_use1_legacy_aps2_routing' does not have provider alias 'aws.alice'")
	}
	if al != "aws.us-east-1" {
		t.Fatalf("Module 'legacy_use1_legacy_aps2_routing' provider alias 'aws.alice' unexpected value %#q, expected 'aws.us-east-1'", al)
	}
}

func TestModule2Params(t *testing.T) {
	m, exists := config.Modules["legacy_use1_g1_aps2_routing"]
	if !exists {
		t.Fatal("Module 'legacy_use1_g1_aps2_routing' was not found")
	}
	al, exists := m.Providers["aws.bob"]
	if !exists {
		t.Fatal("Module 'legacy_use1_g1_aps2_routing' does not have provider alias 'aws.bob'")
	}
	if al != "aws.ap-southeast-2" {
		t.Fatalf("Module 'legacy_use1_g1_aps2_routing' provider alias 'aws.alice' unexpected value %#q, expected 'aws.ap-southeast-2'", al)
	}
}

func TestModule2ParamsNo(t *testing.T) {
	m, exists := config.Modules["legacy_use1_g1_aps2_routing"]
	// test case when we check that it must exist is next tst
	if exists {
		if l := len(m.Parameters); l != 3 {
			t.Fatalf("Unexpected number of parameters found for Module 'legacy_use1_g1_aps2_routing', %v instead of 3", l)
		}
	}
}

func TestModule1Source(t *testing.T) {
	m, exists := config.Modules["legacy_use1_legacy_aps2_routing"]
	if exists {
		expSourcePath := "../../modules/vpc-routing"
		if m.SourcePath != expSourcePath {
			t.Fatalf("Unexpected source path for module 'legacy_use1_legacy_aps2_routing', found %#q, expected %#q", m.SourcePath, expSourcePath)
		}
	}
}
