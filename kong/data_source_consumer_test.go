package kong

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceKongConsumer(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConsumerDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.kong_consumer.consumer_data_source", "username", "User777"),
					resource.TestCheckResourceAttr("data.kong_consumer.consumer_data_source", "custom_id", "123456"),
				),
			},
		},
	})
}

const testConsumerDataSourceConfig = `
resource "kong_consumer" "test_consumer" {
	username  = "User777"
	custom_id = "123456"
}

data "kong_consumer" "consumer_data_source" {
	filter {
		id 		  = "${kong_consumer.test_consumer.id}"
		username  = "User777"
		custom_id = "123456"
	}
}
`
