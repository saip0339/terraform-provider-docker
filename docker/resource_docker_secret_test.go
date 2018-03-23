package docker

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDockerSecret_basicNotUpdatable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckDockerSecretDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: `
				resource "docker_secret" "foo" {
					name = "foo-secret"
					data = "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("docker_secret.foo", "name", "foo-secret"),
					resource.TestCheckResourceAttr("docker_secret.foo", "updatable", "false"),
					resource.TestCheckResourceAttr("docker_secret.foo", "data", "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="),
				),
			},
		},
	})
}
func TestAccDockerSecret_basicUpdateble(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckDockerSecretShouldStillExist,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: `
				resource "docker_secret" "foo" {
					name 			 = "tftest-mysecret-${replace(timestamp(),":", ".")}"
					data 			 = "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="
					updatable = true

					lifecycle {
						ignore_changes = ["name"]
					}
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr("docker_secret.foo", "name", "foo"),
					resource.TestCheckResourceAttr("docker_secret.foo", "updatable", "true"),
					resource.TestCheckResourceAttr("docker_secret.foo", "data", "Ymxhc2RzYmxhYmxhMTI0ZHNkd2VzZA=="),
				),
			},
		},
	})
}

/////////////
// Helpers
/////////////
func testCheckDockerSecretDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).DockerClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "secrets" {
			continue
		}

		id := rs.Primary.Attributes["id"]
		secret, err := client.InspectSecret(id)

		if err == nil || secret != nil {
			return fmt.Errorf("Secret with id '%s' still exists", id)
		}
		return nil
	}
	return nil
}

func testCheckDockerSecretShouldStillExist(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).DockerClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "secrets" {
			continue
		}

		id := rs.Primary.Attributes["id"]
		secret, err := client.InspectSecret(id)

		if err != nil || secret == nil {
			return fmt.Errorf("Secret with id '%s' is destroyed but it should exist", id)
		}
		return nil
	}
	return nil
}
