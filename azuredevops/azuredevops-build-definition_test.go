package azuredevops

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"terraform-provider-azuredevops/azuredevops/utils"
	"testing"
)

func Test_buildDefinitionCheck(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testBuildDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testBuildDefCheckBasicMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testBuildDefExist("azuredevops_build_definition.test1"),
				),
			},
		},
	})
}

func testBuildDefCheckBasicMinimal() string {
	return fmt.Sprintf(
		`

  resource "azuredevops_project" "test1" {
			name  = "project-%v"
  }
   resource "azuredevops_build_definition" "test1" {
			name  = "build-def-%v"
			project_id ="${azuredevops_project.test1.name}"
			buildnumber_format = "1.0$(rev:.r)"
			badge_enabled = false
			description ="It's buiild definition created by Terraform"
			repository {
				name = "${azuredevops_project.test1.name}"
				type = "TfsGit"
                branch = "master"
			}
			variables {
				variable {
					name ="test"
                	value ="ok"
				}
			}
			designer_phase {
				name = "phase1"
				step {
					display_name="teststep"
					task_id = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"
					task_version ="2.*"
					enabled = true
					always_run = true
					condition = "always()"
					inputs = {
							failOnStderr= "false"
							script= "echo Write your commands here↵↵echo Use the environment variables input below to pass secret variables to this script"
							workingDirectory= ""
					}
					
				}

				step {
					display_name="teststep2"
					task_id = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"
					task_version ="2.*"
					enabled = true
					always_run = false
					continue_on_error = true
					timeout_in_minutes = 50
					inputs = {
							failOnStderr= "false"
							script= "echo Write your commands here↵↵echo Use the environment variables input below to pass secret variables to this script"
							workingDirectory= ""
					}
					environment_variables = {
							var1 = "key1"
							var2 = "key2"
					}
					
				}
				
			}

			designer_phase {
				name = "phase2"
				step {
					display_name="teststep3"
					task_id = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"
					task_version ="2.*"
					enabled = false
					continue_on_error = false
					reference_name ="testrefname"
					inputs = {
							failOnStderr= "false"
							script= "echo Write your commands here↵↵echo Use the environment variables input below to pass secret variables to this script"
							workingDirectory= ""
					}
					
				}

				
				
			}

			queue {
            	pool_name ="Hosted VS2017"
			}
  }
`, utils.StringRandom(5), utils.StringRandom(3))
}

func testBuildDefinitionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AzureDevOpsClient)
	client := testAccProvider.Meta().(*AzureDevOpsClient).buildClient
	ctx := testAccProvider.Meta().(*AzureDevOpsClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuredevops_check" {
			continue
		}

		defid, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		projectid := rs.Primary.Attributes["project_id"]
		build, err := client.GetDefinition(ctx, conn.organization, projectid, int32(defid), nil, nil, "", nil)

		if &build != nil {
			return fmt.Errorf("Bad: Build Definition %q still exists", rs.Primary.ID)
		}

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func testBuildDefExist(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		c := testAccProvider.Meta().(*AzureDevOpsClient)
		client := testAccProvider.Meta().(*AzureDevOpsClient).buildClient
		ctx := testAccProvider.Meta().(*AzureDevOpsClient).StopContext

		projectid := rs.Primary.Attributes["project_id"]
		defid, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		build, err := client.GetDefinition(ctx, c.organization, projectid, int32(defid), nil, nil, "", nil)

		if err != nil {
			return err
		}

		if &build == nil {
			return err
		}
		return nil
	}
}
