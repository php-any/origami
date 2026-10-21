package php

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type OpenSSLX509ParseFunction struct{}

func (*OpenSSLX509ParseFunction) GetName() string { return "openssl_x509_parse" }
func (*OpenSSLX509ParseFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "certificate", 0, nil, nil), node.NewParameter(nil, "short_names", 1, data.NewBoolValue(true), nil)}
}
func (*OpenSSLX509ParseFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "certificate", 0, nil), node.NewVariable(nil, "short_names", 1, nil)}
}
func (*OpenSSLX509ParseFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	if value == nil {
		return data.NewBoolValue(false), nil
	}
	encoded := value.AsString()
	if strings.HasPrefix(encoded, "file://") {
		b, err := data.ReadPHPFile(ctx.GetVM(), strings.TrimPrefix(encoded, "file://"))
		if err != nil {
			return data.NewBoolValue(false), nil
		}
		encoded = string(b)
	}
	block, _ := pem.Decode([]byte(encoded))
	if block == nil || block.Type != "CERTIFICATE" {
		return data.NewBoolValue(false), nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return data.NewBoolValue(false), nil
	}
	short := true
	if v, ok := ctx.GetIndexValue(1); ok {
		if b, ok := v.(data.AsBool); ok {
			short, _ = b.AsBool()
		}
	}
	out := data.NewArrayValueFromSlots(nil)
	out.SetStringKey("name", data.NewStringValue(cert.Subject.String()))
	out.SetStringKey("subject", x509Name(cert.Subject, short))
	out.SetStringKey("issuer", x509Name(cert.Issuer, short))
	out.SetStringKey("version", data.NewIntValue(cert.Version-1))
	out.SetStringKey("serialNumber", data.NewStringValue(cert.SerialNumber.String()))
	out.SetStringKey("serialNumberHex", data.NewStringValue(strings.ToUpper(cert.SerialNumber.Text(16))))
	out.SetStringKey("validFrom", data.NewStringValue(cert.NotBefore.UTC().Format("060102150405Z")))
	out.SetStringKey("validTo", data.NewStringValue(cert.NotAfter.UTC().Format("060102150405Z")))
	out.SetStringKey("validFrom_time_t", data.NewIntValue(int(cert.NotBefore.Unix())))
	out.SetStringKey("validTo_time_t", data.NewIntValue(int(cert.NotAfter.Unix())))
	extensions := data.NewArrayValueFromSlots(nil)
	var san []string
	for _, name := range cert.DNSNames {
		san = append(san, "DNS:"+name)
	}
	for _, ip := range cert.IPAddresses {
		san = append(san, "IP Address:"+ip.String())
	}
	for _, email := range cert.EmailAddresses {
		san = append(san, "email:"+email)
	}
	for _, uri := range cert.URIs {
		san = append(san, "URI:"+uri.String())
	}
	if len(san) > 0 {
		extensions.SetStringKey("subjectAltName", data.NewStringValue(strings.Join(san, ", ")))
	}
	out.SetStringKey("extensions", extensions)
	return out, nil
}
func x509Name(name pkix.Name, short bool) *data.ArrayValue {
	out := data.NewArrayValueFromSlots(nil)
	fields := []struct {
		short, long string
		values      []string
	}{
		{"C", "countryName", name.Country}, {"ST", "stateOrProvinceName", name.Province},
		{"L", "localityName", name.Locality}, {"O", "organizationName", name.Organization},
		{"OU", "organizationalUnitName", name.OrganizationalUnit},
	}
	for _, f := range fields {
		if len(f.values) == 0 {
			continue
		}
		key := f.short
		if !short {
			key = f.long
		}
		if len(f.values) == 1 {
			out.SetStringKey(key, data.NewStringValue(f.values[0]))
		} else {
			var values []data.Value
			for _, v := range f.values {
				values = append(values, data.NewStringValue(v))
			}
			out.SetStringKey(key, data.NewArrayValue(values))
		}
	}
	if name.CommonName != "" {
		key := "CN"
		if !short {
			key = "commonName"
		}
		out.SetStringKey(key, data.NewStringValue(name.CommonName))
	}
	return out
}
