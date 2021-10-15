# us-east-1 - ap-southeast-2 routing

# this is a silly comment
module "module1" {
  source                 = "../../modules/module_name"
  alice_vpc_name         = "Development VPC"
  bob_vpc_name           = "Development VPC"
  boolean_value = true
  numeric_value = 12

  providers = {
    /* This provider is not needed...
    aws.test = aws.unused
    */
    aws.alice = aws.us-east-1
    aws.bob   = aws.ap-southeast-2
    // aws.bob = aws.us-east-1 # - wrong value
  }
}

/*
module "one_test" {
  source = "../../modules/testing"
  parameter1 = "value1"
}
*/

# use1.genhtcc.com <-> aps2.g1.genhtcc.com
module "module2" {
  source                 = "../../modules/vpc-routing"
  alice_vpc_name         = "Development VPC"
  bob_vpc_name           = "some.other.name"
  // Do we still need all those providers?
  providers = {
    aws.alice = aws.us-east-1
    aws.bob   = aws.ap-southeast-2
  }
}
