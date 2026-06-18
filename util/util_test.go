package util

import (
	"strings"
	"testing"

	"github.com/ngoduykhanh/wireguard-ui/model"
)

func TestBuildBackupClientConfigOnlyChangesEndpoint(t *testing.T) {
	client := model.Client{
		PrivateKey:   "client-private",
		PresharedKey: "client-psk",
		AllocatedIPs: []string{"11.13.115.2/32"},
		AllowedIPs:   []string{"0.0.0.0/0"},
		UseServerDNS: true,
	}
	server := model.Server{
		KeyPair: &model.ServerKeypair{PublicKey: "server-public"},
		Interface: &model.ServerInterface{
			ListenPort: 51823,
		},
	}
	setting := model.GlobalSetting{
		EndpointAddress:       "primary.example.org",
		BackupEndpointAddress: "wg2.int.automesh.org",
		DNSServers:            []string{"1.1.1.1"},
		MTU:                   1392,
		PersistentKeepalive:   60,
	}

	primary := BuildClientConfig(client, server, setting)
	backup := BuildBackupClientConfig(client, server, setting)

	if !strings.Contains(primary, "Endpoint = primary.example.org:51823\n") {
		t.Fatalf("primary endpoint missing from config:\n%s", primary)
	}
	if !strings.Contains(backup, "Endpoint = wg2.int.automesh.org:51823\n") {
		t.Fatalf("backup endpoint missing from config:\n%s", backup)
	}

	withoutEndpoint := func(config string) string {
		lines := strings.Split(config, "\n")
		kept := lines[:0]
		for _, line := range lines {
			if strings.HasPrefix(line, "Endpoint = ") {
				continue
			}
			kept = append(kept, line)
		}
		return strings.Join(kept, "\n")
	}

	if withoutEndpoint(primary) != withoutEndpoint(backup) {
		t.Fatalf("backup config changed fields other than Endpoint\nprimary:\n%s\nbackup:\n%s", primary, backup)
	}
}
