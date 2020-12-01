package kong

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kevholditch/gokong"
)

func readStringArrayPtrFromResource(d *schema.ResourceData, key string) []*string {

	if attr, ok := d.GetOk(key); ok {
		var array []string
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(string)
			array = append(array, item)
		}

		return gokong.StringSlice(array)
	}

	return nil
}

func readIpPortArrayFromResource(d *schema.ResourceData, key string) []*gokong.IpPort {
	if attr, ok := d.GetOk(key); ok {
		set := attr.(*schema.Set)
		results := make([]*gokong.IpPort, 0)
		for _, item := range set.List() {
			m := item.(map[string]interface{})
			ipPort := &gokong.IpPort{}
			if port, ok := m["port"].(int); ok && port != 0 {
				ipPort.Port = &port
			}
			if ip, ok := m["ip"].(string); ok && ip != "" {
				ipPort.Ip = &ip
			}
			results = append(results, ipPort)
		}

		return results
	}

	return nil
}

func readArrayFromResource(d *schema.ResourceData, key string) []interface{} {
	if attr, ok := d.GetOk(key); ok {
		return attr.([]interface{})
	}

	return nil
}

func readStringFromResource(d *schema.ResourceData, key string) string {
	if value, ok := d.GetOk(key); ok {
		return value.(string)
	}
	return ""
}

func readIdPtrFromResource(d *schema.ResourceData, key string) *gokong.Id {
	if value, ok := d.GetOk(key); ok {
		id := gokong.Id(value.(string))
		return &id
	}
	return nil
}

func readStringPtrFromResource(d *schema.ResourceData, key string) *string {
	if value, ok := d.GetOkExists(key); ok {
		return gokong.String(value.(string))
	}
	return nil
}

func readBoolPtrFromResource(d *schema.ResourceData, key string) *bool {
	return gokong.Bool(d.Get(key).(bool))
}

func readIntFromResource(d *schema.ResourceData, key string) int {
	if value, ok := d.GetOk(key); ok {
		return value.(int)
	}
	return 0
}

func readIntPtrFromResource(d *schema.ResourceData, key string) *int {
	if value, ok := d.GetOkExists(key); ok {
		return gokong.Int(value.(int))
	}
	return nil
}
