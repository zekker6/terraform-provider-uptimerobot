package uptimerobot

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/louy/terraform-provider-uptimerobot/uptimerobot/api"
)

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorCreate,
		Read:   resourceMonitorRead,
		Update: resourceMonitorUpdate,
		Delete: resourceMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"friendly_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(uptimerobotapi.MonitorType, false),
			},
			"sub_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(uptimerobotapi.MonitorSubType, false),
				// required for port monitoring
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				// required for port monitoring
			},
			"keyword_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(uptimerobotapi.MonitorKeywordType, false),
				// required for keyword monitoring
			},
			"keyword_value": {
				Type:     schema.TypeString,
				Optional: true,
				// required for keyword monitoring
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"http_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"http_success_codes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"http_down_codes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alert_contact": {
				Type:     schema.TypeList,
				Optional: true,
				// PromoteSingle: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"recurrence": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"custom_http_headers": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			// TODO - mwindows
			// TODO - ignore_ssl_errors
		},
	}
}

func resourceMonitorCreate(d *schema.ResourceData, m interface{}) error {
	req := uptimerobotapi.MonitorCreateRequest{
		FriendlyName: d.Get("friendly_name").(string),
		URL:          d.Get("url").(string),
		Type:         d.Get("type").(string),
	}

	switch req.Type {
	case "port":
		req.SubType = d.Get("sub_type").(string)
		req.Port = d.Get("port").(int)
		break
	case "keyword":
		req.KeywordType = d.Get("keyword_type").(string)
		req.KeywordValue = d.Get("keyword_value").(string)

		req.HTTPUsername = d.Get("http_username").(string)
		req.HTTPPassword = d.Get("http_password").(string)
		break
	case "http":
		req.HTTPUsername = d.Get("http_username").(string)
		req.HTTPPassword = d.Get("http_password").(string)
		break
	}

	req.HTTPSuccessCodes = make([]int, len(d.Get("http_success_codes").([]interface{})))
	for k, v := range d.Get("http_success_codes").([]interface{}) {
		req.HTTPSuccessCodes[k] = v.(int)
	}

	req.HTTPDownCodes = make([]int, len(d.Get("http_down_codes").([]interface{})))
	for k, v := range d.Get("http_down_codes").([]interface{}) {
		req.HTTPDownCodes[k] = v.(int)
	}

	// Add optional attributes
	req.Interval = d.Get("interval").(int)

	req.AlertContacts = make([]uptimerobotapi.MonitorRequestAlertContact, len(d.Get("alert_contact").([]interface{})))
	for k, v := range d.Get("alert_contact").([]interface{}) {
		req.AlertContacts[k] = uptimerobotapi.MonitorRequestAlertContact{
			ID:         v.(map[string]interface{})["id"].(string),
			Threshold:  v.(map[string]interface{})["threshold"].(int),
			Recurrence: v.(map[string]interface{})["recurrence"].(int),
		}
	}

	// custom_http_headers
	httpHeaderMap := d.Get("custom_http_headers").(map[string]interface{})
	req.CustomHTTPHeaders = make(map[string]string, len(httpHeaderMap))
	for k, v := range httpHeaderMap {
		req.CustomHTTPHeaders[k] = v.(string)
	}

	monitor, err := m.(uptimerobotapi.UptimeRobotApiClient).CreateMonitor(req)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", monitor.ID))
	updateMonitorResource(d, monitor)
	return nil
}

func resourceMonitorRead(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	monitor, err := m.(uptimerobotapi.UptimeRobotApiClient).GetMonitor(id)
	if err != nil {
		return err
	}

	updateMonitorResource(d, monitor)

	return nil
}

func resourceMonitorUpdate(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	req := uptimerobotapi.MonitorUpdateRequest{
		ID:           id,
		FriendlyName: d.Get("friendly_name").(string),
		URL:          d.Get("url").(string),
		Type:         d.Get("type").(string),
	}

	switch req.Type {
	case "port":
		req.SubType = d.Get("sub_type").(string)
		req.Port = d.Get("port").(int)
		break
	case "keyword":
		req.KeywordType = d.Get("keyword_type").(string)
		req.KeywordValue = d.Get("keyword_value").(string)

		req.HTTPUsername = d.Get("http_username").(string)
		req.HTTPPassword = d.Get("http_password").(string)
		break
	case "http":
		req.HTTPUsername = d.Get("http_username").(string)
		req.HTTPPassword = d.Get("http_password").(string)
		break
	}

	// Add optional attributes
	req.Interval = d.Get("interval").(int)

	req.AlertContacts = make([]uptimerobotapi.MonitorRequestAlertContact, len(d.Get("alert_contact").([]interface{})))
	for k, v := range d.Get("alert_contact").([]interface{}) {
		req.AlertContacts[k] = uptimerobotapi.MonitorRequestAlertContact{
			ID:         v.(map[string]interface{})["id"].(string),
			Threshold:  v.(map[string]interface{})["threshold"].(int),
			Recurrence: v.(map[string]interface{})["recurrence"].(int),
		}
	}

	req.HTTPSuccessCodes = make([]int, len(d.Get("http_success_codes").([]interface{})))
	for k, v := range d.Get("http_success_codes").([]interface{}) {
		req.HTTPSuccessCodes[k] = v.(int)
	}

	req.HTTPDownCodes = make([]int, len(d.Get("http_down_codes").([]interface{})))
	for k, v := range d.Get("http_down_codes").([]interface{}) {
		req.HTTPDownCodes[k] = v.(int)
	}

	// custom_http_headers
	httpHeaderMap := d.Get("custom_http_headers").(map[string]interface{})
	req.CustomHTTPHeaders = make(map[string]string, len(httpHeaderMap))
	for k, v := range httpHeaderMap {
		req.CustomHTTPHeaders[k] = v.(string)
	}

	monitor, err := m.(uptimerobotapi.UptimeRobotApiClient).UpdateMonitor(req)
	if err != nil {
		return err
	}

	updateMonitorResource(d, monitor)
	return nil
}

func resourceMonitorDelete(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = m.(uptimerobotapi.UptimeRobotApiClient).DeleteMonitor(id)
	if err != nil {
		return err
	}
	return nil
}

func updateMonitorResource(d *schema.ResourceData, m uptimerobotapi.Monitor) {
	d.Set("friendly_name", m.FriendlyName)
	d.Set("url", m.URL)
	d.Set("type", m.Type)
	d.Set("status", m.Status)
	d.Set("interval", m.Interval)

	d.Set("sub_type", m.SubType)
	d.Set("port", m.Port)

	d.Set("keyword_type", m.KeywordType)
	d.Set("keyword_value", m.KeywordValue)

	d.Set("http_username", m.HTTPUsername)
	d.Set("http_password", m.HTTPPassword)

	d.Set("custom_http_headers", m.CustomHTTPHeaders)

	d.Set("http_down_codes", m.HTTPDownCodes)
	d.Set("http_success_codes", m.HTTPSuccessCodes)
}
