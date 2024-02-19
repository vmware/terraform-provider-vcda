// Copyright (c) 2023-2024 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vcda

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVcdaRemoteServicesThumbprint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVcdaRemoteServicesThumbprintRead,
		Schema: map[string]*schema.Schema{
			"address": {
				Type: schema.TypeString,
				Description: "The address of the remote appliance/service. " +
					"**NOTE:** this method produces a thumbprint that is not verified nor safe for use.",
				Optional:      true,
				ConflictsWith: []string{"pem_file"},
				RequiredWith:  []string{"port"},
			},
			"port": {
				Type:        schema.TypeString,
				Description: "The port of the remote appliance/service. Use only with `address`.",
				Optional:    true,
			},
			"pem_file": {
				Type: schema.TypeString,
				Description: "The name of the file that contains the last certificate " +
					"in the chain (end entity cert) of the remote appliance/service in PEM format. " +
					"On creation, include either `pem_file` or `address`.",
				Optional:      true,
				ConflictsWith: []string{"address"},
			},
		},
	}
}

func dataSourceVcdaRemoteServicesThumbprintRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	address := d.Get("address").(string)
	port := d.Get("port").(string)
	pemFile := d.Get("pem_file").(string)

	if address != "" {
		err := computeHostThumbprint(address, port, d)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	if pemFile != "" {
		err := computeThumbprintFromFile(pemFile, d)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func computeHostThumbprint(address string, port string, d *schema.ResourceData) error {
	config := &tls.Config{}
	config.InsecureSkipVerify = true

	conn, err := tls.Dial("tcp", address+":"+port, config)
	if err != nil {
		return err
	}
	cert := conn.ConnectionState().PeerCertificates[0]

	fingerprint := sha256.Sum256(cert.Raw)

	thumbprint := formatFingerprint(fingerprint)
	d.SetId(thumbprint)

	return nil
}

func computeThumbprintFromFile(pemFile string, d *schema.ResourceData) error {
	cert, err := os.ReadFile(pemFile)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(cert)

	if block == nil {
		return fmt.Errorf("failed to decode PEM file - invalid PEM format")
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse PEM file")
	}

	fingerprint := sha256.Sum256(caCert.Raw)

	thumbprint := formatFingerprint(fingerprint)
	d.SetId(thumbprint)

	return nil
}

func formatFingerprint(fingerprint [32]byte) string {
	var buf bytes.Buffer

	buf.WriteString("SHA-256:")
	for i, f := range fingerprint {
		if i > 0 {
			_, _ = fmt.Fprintf(&buf, ":")
		}
		_, _ = fmt.Fprintf(&buf, "%02X", f)
	}

	return buf.String()
}
