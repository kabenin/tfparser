package tfparser

import (
	"os"
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

func TestParseFile(t *testing.T) {
	testFileName := "testdata/tf/main.tf"
	_, err := os.Stat(testFileName)
	if os.IsNotExist(err) {
		t.Skipf("testing terraform file is not found. Not testing ParseFile")
	}
	config, err := ParseFile(testFileName)
	if err != nil {
		t.Fatalf("ParseFile returned an error, %v", err)
	}
	assertTestingConfig1(config, t)
}

func TestParseDir(t *testing.T) {
	testDirName := "testdata/tf/"
	_, err := os.Stat(testDirName)
	if os.IsNotExist(err) {
		t.Skipf("testing terraform dir %#q is not found. Not testing ParseDir", testDirName)
	}
	config, err := ParseDir(testDirName)
	if err != nil {
		t.Fatalf("ParseDir returned an error, %v", err)
	}
	assertTestingConfig1(config, t)
}

func assertTestingConfig1(config *TFconfig, t *testing.T) {
	if len(config.Modules) != 2 {
		t.Fatalf("Unexpected number of moduels fetched. Expected 2 got %v", len(config.Modules))
	}
	m1, exists := config.Modules["module1"]
	if !exists {
		t.Fatalf("Module 'module1' is not found, but expected")
	}
	m1sp := "../../modules/module_name"
	if m1.SourcePath != m1sp {
		t.Fatalf("Unexpceted 'module1' source: %#q, expected %#q", m1.SourcePath, m1sp)
	}
	if len(m1.Providers) != 2 {
		t.Fatalf("Unexpected number of providers for 'module1' %v, expected 2", len(m1.Providers))
	}
	p1, exists := m1.Providers["aws.alice"]
	if !exists {
		t.Fatal("provider 'aws.alice' was not found in 'module1")
	}
	if p1 != "aws.us-east-1" {
		t.Fatalf("provider 'aws.alice' alias is %#q, expected 'aws.us-east-1'", p1)
	}
}
