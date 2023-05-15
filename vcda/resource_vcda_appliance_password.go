/* Copyright 2023 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vcda

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVcdaAppliancePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppliancePasswordCreate,
		ReadContext:   resourceAppliancePasswordRead,
		UpdateContext: resourceAppliancePasswordUpdate,
		DeleteContext: resourceAppliancePasswordDelete,
		Schema: map[string]*schema.Schema{
			"current_password": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Description: "The current password of the appliance.",
				Required:    true,
			},
			"new_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Description: "The new password of the appliance. Note: This value is never returned on read. " +
					"On creation, include either `new_password` or `password_file`.",
				Optional:      true,
				ConflictsWith: []string{"password_file"},
			},
			"password_file": {
				Type: schema.TypeString,
				Description: "The name of a file containing the appliance password. " +
					"On creation, include either `password_file` or `new_password`.",
				Optional:      true,
				ConflictsWith: []string{"new_password"},
			},
			"service_cert": {
				Type:        schema.TypeString,
				Description: "The service certificate.",
				Required:    true,
			},
			"appliance_ip": {
				Type:        schema.TypeString,
				Description: "The IP address of the appliance.",
				Required:    true,
			},
			//computed
			"root_password_expired": {
				Type:        schema.TypeBool,
				Description: "Flag indicating whether the **root** user password is already expired.",
				Computed:    true,
			},
			"seconds_until_expiration": {
				Type:        schema.TypeInt,
				Description: "Seconds until the **root** user password expires.",
				Computed:    true,
			},
		},
	}

}

func resourceAppliancePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	currentPassword := d.Get("current_password").(string)
	newPassword := d.Get("new_password").(string)
	passwordFile := d.Get("password_file").(string)
	applianceIP := d.Get("appliance_ip").(string)
	serviceCert := d.Get("service_cert").(string)

	if currentPassword == "" {
		return diag.Errorf(`"current_password" cannot be empty`)
	}

	newPass, err := getNewPasswordInput(newPassword, passwordFile)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.changePassword(applianceIP, currentPassword, *newPass, serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceAppliancePasswordRead(ctx, d, m)
}

func resourceAppliancePasswordRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	applianceIP := d.Get("appliance_ip").(string)
	serviceCert := d.Get("service_cert").(string)

	passExpiration, err := c.checkPasswordExpired(applianceIP, serviceCert)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setPasswordData(d, passExpiration); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAppliancePasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*Client)

	if d.HasChange("new_password") || d.HasChange("password_file") {
		newPassword := d.Get("new_password").(string)
		passwordFile := d.Get("password_file").(string)
		currentPassword := d.Get("current_password").(string)
		applianceIP := d.Get("appliance_ip").(string)
		serviceCert := d.Get("service_cert").(string)

		if currentPassword == "" {
			return diag.Errorf(`"current_password" cannot be empty`)
		}

		newPass, err := getNewPasswordInput(newPassword, passwordFile)
		if err != nil {
			return diag.FromErr(err)
		}

		err = c.changePassword(applianceIP, currentPassword, *newPass, serviceCert)

		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

		return resourceAppliancePasswordRead(ctx, d, m)
	}
	return diags
}

func resourceAppliancePasswordDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}

func setPasswordData(d *schema.ResourceData, passExpiration *PasswordExpiration) error {
	if err := d.Set("root_password_expired", passExpiration.RootPasswordExpired); err != nil {
		return fmt.Errorf("error setting root_password_expired field: %s", err)
	}

	if err := d.Set("seconds_until_expiration", passExpiration.SecondsUntilExpiration); err != nil {
		return fmt.Errorf("error setting seconds_until_expiration field: %s", err)
	}

	return nil
}

func getNewPasswordInput(newPassword string, passwordFile string) (*string, error) {
	var newPass string
	if newPassword != "" {
		newPass = newPassword
	}

	if newPassword != "" && passwordFile != "" {
		return nil, fmt.Errorf(`either "new_password" or "password_file" should be given, but not both`)
	}

	if passwordFile != "" {
		passwordBytes, err := os.ReadFile(filepath.Clean(passwordFile))
		if err != nil {
			return nil, err
		}
		passwordStr := strings.TrimSpace(string(passwordBytes))
		if passwordStr != "" {
			newPass = passwordStr
		}
	}

	return &newPass, nil
}
