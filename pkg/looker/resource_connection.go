package looker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/looker-open-source/sdk-codegen/go/rtl"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v3"
)

const (
	connectionsCreationURL = "%s/api/3.1/connections"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,
		Exists: resourceConnectionExists,
		Importer: &schema.ResourceImporter{
			State: resourceConnectionImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // todo: the ID is the name of the connection so if it changes i think it would require a new object be created.  I should verify this
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
				ValidateFunc: validation.StringDoesNotContainAny(" "),
			},
			"dialect_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "443",
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"jdbc_additional_params": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"db_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)

	params := apiclient.WriteDBConnection{}
	connectionName := d.Get("name").(string)
	params.Name = &connectionName

	dialectName := d.Get("dialect_name").(string)
	params.DialectName = &dialectName

	host := d.Get("host").(string)
	params.Host = &host

	port := d.Get("port").(string)
	params.Port = &port

	database := d.Get("database").(string)
	params.Database = &database

	userName := d.Get("username").(string)
	params.Username = &userName

	password := d.Get("password").(string)
	params.Password = &password

	schema := d.Get("schema").(string)
	params.Schema = &schema

	jdbcAdditionalParams := d.Get("jdbc_additional_params").(string)
	params.JdbcAdditionalParams = &jdbcAdditionalParams

	ssl := d.Get("ssl").(bool)
	params.Ssl = &ssl

	if dbTimezone, ok := d.GetOk("db_timezone"); ok {
		dbTimezone := dbTimezone.(string)
		params.DbTimezone = &dbTimezone
	}

	if queryTimezone, ok := d.GetOk("query_timezone"); ok {
		queryTimezone := queryTimezone.(string)
		params.QueryTimezone = &queryTimezone
	}

	body, _ := JSONMarshal(params)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(connectionsCreationURL, session.Config.BaseUrl), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	session.Authenticate(req)
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		bodyResult, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return err
		}

		log.Printf("%s", bodyResult)
		return err
	}

	d.SetId(*params.Name)

	return resourceConnectionRead(d, m)
}

func resourceConnectionRead(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	result, err := client.Connection(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", *result.Name)
	d.Set("dialect_name", *result.DialectName)
	d.Set("host", *result.Host)
	d.Set("port", *result.Port)
	d.Set("database", *result.Database)
	d.Set("username", *result.Username)
	d.Set("password", d.Get("password").(string))
	d.Set("schema", *result.Schema)
	d.Set("jdbc_additional_params", *result.JdbcAdditionalParams)
	d.Set("ssl", *result.Ssl)

	if result.DbTimezone != nil {
		d.Set("db_timezone", *result.DbTimezone)
	}
	if result.QueryTimezone != nil {
		d.Set("query_timezone", *result.QueryTimezone)
	}

	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	params := apiclient.WriteDBConnection{}
	connectionName := d.Get("name").(string)
	params.Name = &connectionName

	dialectName := d.Get("dialect_name").(string)
	params.DialectName = &dialectName

	host := d.Get("host").(string)
	params.Host = &host

	port := d.Get("port").(string)
	params.Port = &port

	database := d.Get("database").(string)
	params.Database = &database

	userName := d.Get("username").(string)
	params.Username = &userName

	password := d.Get("password").(string)
	params.Password = &password

	schema := d.Get("schema").(string)
	params.Schema = &schema

	jdbcAdditionalParams := d.Get("jdbc_additional_params").(string)
	params.JdbcAdditionalParams = &jdbcAdditionalParams

	ssl := d.Get("ssl").(bool)
	params.Ssl = &ssl

	if dbTimezone, ok := d.GetOk("db_timezone"); ok {
		dbTimezone := dbTimezone.(string)
		params.DbTimezone = &dbTimezone
	}

	if queryTimezone, ok := d.GetOk("query_timezone"); ok {
		queryTimezone := queryTimezone.(string)
		params.QueryTimezone = &queryTimezone
	}

	_, err := client.UpdateConnection(d.Id(), params, nil)
	if err != nil {
		return err
	}

	return resourceConnectionRead(d, m)
}

func resourceConnectionDelete(d *schema.ResourceData, m interface{}) error {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	_, err := client.DeleteConnection(d.Id(), nil)
	if err != nil {
		return err
	}

	return nil
}

func resourceConnectionExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	session := m.(*rtl.AuthSession)
	client := apiclient.NewLookerSDK(session)

	_, err := client.Connection(d.Id(), "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceConnectionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceConnectionRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
