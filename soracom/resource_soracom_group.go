package soracom

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/soracom/soracom-sdk-go"
)

func resourceSoracomGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSoracomGroupCreate,
		Read:   resourceSoracomGroupRead,
		Update: resourceSoracomGroupUpdate,
		Delete: resourceSoracomGroupDelete,

		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceSoracomGroupCreate(d *schema.ResourceData, meta interface{}) error {
	var tags = soracom.Tags{}
	for k, v := range d.Get("tags").(map[string]interface{}) {
		tags[k] = v.(string)
	}

	client := meta.(*soracom.APIClient)
	group, err := client.CreateGroup(tags)
	if err != nil {
		return err
	}

	d.SetId(group.GroupID)

	return resourceSoracomGroupRead(d, meta)
}

func resourceSoracomGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*soracom.APIClient)
	group, err := client.GetGroup(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	for k, v := range group.Tags {
		d.Set(k, v)
	}

	return nil
}

func resourceSoracomGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("tags") {
		client := meta.(*soracom.APIClient)
		group, err := client.GetGroup(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				d.SetId("")
				return nil
			}
			return err
		}

		for k := range group.Tags {
			client.DeleteGroupTag(d.Id(), k)
			if err != nil {
				if isNotFoundError(err) {
					d.SetId("")
					return nil
				}
				return err
			}
		}

		tags := []soracom.Tag{}
		for k, v := range d.Get("tags").(map[string]interface{}) {
			tags = append(tags, soracom.Tag{TagName: k, TagValue: v.(string)})
		}

		_, err = client.UpdateGroupTags(d.Id(), tags)
		if err != nil {
			if isNotFoundError(err) {
				d.SetId("")
				return nil
			}
			return err
		}
	}

	return resourceSoracomGroupRead(d, meta)
}

func resourceSoracomGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*soracom.APIClient)
	err := client.DeleteGroup(d.Id())
	if err != nil {
		if !isNotFoundError(err) {
			return err
		}
	}

	d.SetId("")

	return nil
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	if apiError, ok := err.(*soracom.APIError); ok {
		return apiError.HTTPStatusCode == http.StatusNotFound
	}

	return false
}
