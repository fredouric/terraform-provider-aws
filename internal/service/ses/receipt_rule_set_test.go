// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ses_test

import (
	"context"
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfses "github.com/hashicorp/terraform-provider-aws/internal/service/ses"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSESReceiptRuleSet_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ses_receipt_rule_set.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t); testAccPreCheckReceiptRule(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckReceiptRuleSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReceiptRuleSetConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReceiptRuleSetExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "rule_set_name", rName),
					acctest.CheckResourceAttrRegionalARN(ctx, resourceName, names.AttrARN, "ses", fmt.Sprintf("receipt-rule-set/%s", rName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSESReceiptRuleSet_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ses_receipt_rule_set.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t); testAccPreCheckReceiptRule(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SESServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckReceiptRuleSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReceiptRuleSetConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReceiptRuleSetExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfses.ResourceReceiptRuleSet(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckReceiptRuleSetDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).SESClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_ses_receipt_rule_set" {
				continue
			}

			_, err := tfses.FindReceiptRuleSetByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("SES Receipt Rule Set %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckReceiptRuleSetExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SESClient(ctx)

		_, err := tfses.FindReceiptRuleSetByName(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccReceiptRuleSetConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_ses_receipt_rule_set" "test" {
  rule_set_name = %[1]q
}
`, rName)
}
