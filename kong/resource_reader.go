package kong

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kevholditch/gokong"
)

func readStringArrayFromResource(d *schema.ResourceData, key string) []string {

	if attr, ok := d.GetOk(key); ok {
		var array []string
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(string)
			array = append(array, item)
		}

		return array
	}

	return nil
}

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

func readIntArrayFromResource(d *schema.ResourceData, key string) []int {

	if attr, ok := d.GetOk(key); ok {
		var array []int
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(int)
			array = append(array, item)
		}

		return array
	}

	return nil
}

func readStringFromResource(d *schema.ResourceData, key string) string {
	if attr, ok := d.GetOk(key); ok {
		return attr.(string)
	}
	return ""
}

func readStringPtrFromResource(d *schema.ResourceData, key string) *string {
	if attr, ok := d.GetOk(key); ok {
		return gokong.String(attr.(string))
	}
	return nil
}

func readBoolPtrFromResource(d *schema.ResourceData, key string) *bool {
	if attr, ok := d.GetOk(key); ok {
		return gokong.Bool(attr.(bool))
	}
	return nil
}

func readIntFromResource(d *schema.ResourceData, key string) int {
	if attr, ok := d.GetOk(key); ok {
		return attr.(int)
	}
	return 0
}

func readIntPtrFromResource(d *schema.ResourceData, key string) *int {
	if attr, ok := d.GetOk(key); ok {
		return gokong.Int(attr.(int))
	}
	return nil
}

func readMapFromResource(d *schema.ResourceData, key string) map[string]interface{} {

	if attr, ok := d.GetOk(key); ok {
		result := attr.(map[string]interface{})
		return result
	}

	return nil
}
