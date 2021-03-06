package prismacloud

import (
	"encoding/json"
	"log"
	"time"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
	"github.com/paloaltonetworks/prisma-cloud-go/cloud/account"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCloudAccount() *schema.Resource {
	return &schema.Resource{
		Create: createCloudAccount,
		Read:   readCloudAccount,
		Update: updateCloudAccount,
		Delete: deleteCloudAccount,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// AWS type.
			account.TypeAws: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AWS account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAzure,
					account.TypeGcp,
					account.TypeAlibaba,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS account ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS account external ID",
							Sensitive:   true,
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"role_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier for an AWS resource (ARN)",
						},
					},
				},
			},

			// Azure type.
			account.TypeAzure: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Azure account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeGcp,
					account.TypeAlibaba,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Azure account ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID registered with Active Directory",
						},
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Application ID key",
						},
						"monitor_flow_logs": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Automatically ingest flow logs",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Active Directory ID associated with Azure",
						},
						"service_principal_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID of the service principal object associated with the Prisma Cloud application that you create",
						},
					},
				},
			},

			// GCP type.
			account.TypeGcp: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "GCP account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeAzure,
					account.TypeAlibaba,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "GCP project ID",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether or not the account is enabled",
							Default:     true,
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"compression_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable flow log compression",
						},
						"dataflow_enabled_project": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GCP project for flow log compression",
						},
						"flow_log_storage_bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GCP flow logs storage bucket",
						},
						// Use a json string until this feature is added:
						// https://github.com/hashicorp/terraform-plugin-sdk/issues/248
						"credentials_json": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Content of the JSON credentials file",
							Sensitive:        true,
							DiffSuppressFunc: gcpCredentialsMatch,
						},
					},
				},
			},

			// Alibaba type.
			account.TypeAlibaba: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alibaba account type",
				MaxItems:    1,
				ConflictsWith: []string{
					account.TypeAws,
					account.TypeAzure,
					account.TypeGcp,
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Alibaba account ID",
						},
						"group_ids": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of account IDs to which you are assigning this account",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name to be used for the account on the Prisma Cloud platform (must be unique)",
						},
						"ram_arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier for an Alibaba RAM role resource",
						},
					},
				},
			},
		},
	}
}

func gcpCredentialsMatch(k, old, new string, d *schema.ResourceData) bool {
	var (
		err       error
		prev, cur account.GcpCredentials
	)

	if err = json.Unmarshal([]byte(old), &prev); err != nil {
		return false
	}

	if err = json.Unmarshal([]byte(new), &cur); err != nil {
		return false
	}

	return (prev.Type == cur.Type &&
		prev.ProjectId == cur.ProjectId &&
		prev.PrivateKeyId == cur.PrivateKeyId &&
		prev.PrivateKey == cur.PrivateKey &&
		prev.ClientEmail == cur.ClientEmail &&
		prev.ClientId == cur.ClientId &&
		prev.AuthUri == cur.AuthUri &&
		prev.TokenUri == cur.TokenUri &&
		prev.ProviderCertUrl == cur.ProviderCertUrl &&
		prev.ClientCertUrl == cur.ClientCertUrl)
}

func parseCloudAccount(d *schema.ResourceData, id string) (string, string, interface{}) {
	if x := ResourceDataInterfaceMap(d, account.TypeAws); len(x) != 0 {
		return account.TypeAws, x["name"].(string), account.Aws{
			AccountId:  id,
			Enabled:    x["enabled"].(bool),
			ExternalId: x["external_id"].(string),
			GroupIds:   ListToStringSlice(x["group_ids"].([]interface{})),
			Name:       x["name"].(string),
			RoleArn:    x["role_arn"].(string),
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeAzure); len(x) != 0 {
		return account.TypeAzure, x["name"].(string), account.Azure{
			Account: account.CloudAccount{
				AccountId: id,
				Enabled:   x["enabled"].(bool),
				GroupIds:  ListToStringSlice(x["group_ids"].([]interface{})),
				Name:      x["name"].(string),
			},
			ClientId:           x["client_id"].(string),
			Key:                x["key"].(string),
			MonitorFlowLogs:    x["monitor_flow_logs"].(bool),
			TenantId:           x["tenant_id"].(string),
			ServicePrincipalId: x["service_principal_id"].(string),
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeGcp); len(x) != 0 {
		var creds account.GcpCredentials
		_ = json.Unmarshal([]byte(x["credentials_json"].(string)), &creds)

		return account.TypeGcp, x["name"].(string), account.Gcp{
			Account: account.CloudAccount{
				AccountId: id,
				Enabled:   x["enabled"].(bool),
				GroupIds:  ListToStringSlice(x["group_ids"].([]interface{})),
				Name:      x["name"].(string),
			},
			CompressionEnabled:     x["compression_enabled"].(bool),
			DataflowEnabledProject: x["dataflow_enabled_project"].(string),
			FlowLogStorageBucket:   x["flow_log_storage_bucket"].(string),
			Credentials:            creds,
		}
	} else if x := ResourceDataInterfaceMap(d, account.TypeAlibaba); len(x) != 0 {
		return account.TypeAlibaba, x["name"].(string), account.Alibaba{
			AccountId: id,
			GroupIds:  ListToStringSlice(x["group_ids"].([]interface{})),
			Name:      x["name"].(string),
			RamArn:    x["ram_arn"].(string),
		}
	}

	return "", "", nil
}

func saveCloudAccount(d *schema.ResourceData, dest string, obj interface{}) {
	var val map[string]interface{}

	switch v := obj.(type) {
	case account.Aws:
		val = map[string]interface{}{
			"account_id":  v.AccountId,
			"enabled":     v.Enabled,
			"external_id": v.ExternalId,
			"group_ids":   v.GroupIds,
			"name":        v.Name,
			"role_arn":    v.RoleArn,
		}
	case account.Azure:
		val = map[string]interface{}{
			"account_id":           v.Account.AccountId,
			"enabled":              v.Account.Enabled,
			"group_ids":            v.Account.GroupIds,
			"name":                 v.Account.Name,
			"client_id":            v.ClientId,
			"key":                  v.Key,
			"monitor_flow_logs":    v.MonitorFlowLogs,
			"tenant_id":            v.TenantId,
			"service_principal_id": v.ServicePrincipalId,
		}
	case account.Gcp:
		b, _ := json.Marshal(v.Credentials)
		val = map[string]interface{}{
			"account_id":               v.Account.AccountId,
			"enabled":                  v.Account.Enabled,
			"group_ids":                v.Account.GroupIds,
			"name":                     v.Account.Name,
			"compression_enabled":      v.CompressionEnabled,
			"dataflow_enabled_project": v.DataflowEnabledProject,
			"flow_log_storage_bucket":  v.FlowLogStorageBucket,
			"credentials_json":         string(b),
		}
	case account.Alibaba:
		val = map[string]interface{}{
			"account_id": v.AccountId,
			"group_ids":  v.GroupIds,
			"name":       v.Name,
			"ram_arn":    v.RamArn,
		}
	}

	for _, key := range []string{account.TypeAws, account.TypeAzure, account.TypeGcp, account.TypeAlibaba} {
		if key != dest {
			d.Set(key, nil)
			continue
		}

		if err := d.Set(key, []interface{}{val}); err != nil {
			log.Printf("[WARN] Error setting %q field for %q: %s", key, d.Id(), err)
		}
	}
}

func createCloudAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	cloudType, name, obj := parseCloudAccount(d, "")

	if err := account.Create(client, obj); err != nil {
		return err
	}

	id, err := account.Identify(client, cloudType, name)
	if err != nil {
		return err
	}

	d.SetId(TwoStringsToId(cloudType, id))
	return readCloudAccount(d, meta)
}

func readCloudAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())

	obj, err := account.Get(client, cloudType, id)
	if err != nil {
		if err == pc.ObjectNotFoundError {
			d.SetId("")
			return nil
		}
		return err
	}

	saveCloudAccount(d, cloudType, obj)

	return nil
}

func updateCloudAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)

	_, id := IdToTwoStrings(d.Id())
	_, _, obj := parseCloudAccount(d, id)

	if err := account.Update(client, obj); err != nil {
		return err
	}

	return readCloudAccount(d, meta)
}

func deleteCloudAccount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pc.Client)
	cloudType, id := IdToTwoStrings(d.Id())

	err := account.Delete(client, cloudType, id)
	if err != nil {
		if err != pc.ObjectNotFoundError {
			return err
		}
	}

	d.SetId("")
	return nil
}
