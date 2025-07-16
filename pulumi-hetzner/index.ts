import * as pulumi from "@pulumi/pulumi";
import * as hcloud from "@pulumi/hcloud";
import * as local from "@pulumi/local";
import * as command from "@pulumi/command";
import * as path from "path";
import * as forge from "node-forge";

// 1. Retrieve configuration values and secrets
const config = new pulumi.Config();
const serverRegion = config.require("serverRegion");
const wireguardServerPrivateIp = config.require("wireguardServerPrivateIp");
const wireguardServerListenPort = config.requireNumber("wireguardServerListenPort");
const wireguardClientPublicKey = config.requireSecret("wireguardClientPublicKey");
const wireguardClientIp = config.require("wireguardClientIp");

// 2. Read the external cloud-init script content
const cloudInitScript = local.getFile({
  filename: path.join(__dirname, "wireguard.yaml"),
});

// 3. Prepare the cloud-init user_data with interpolated secrets and configuration
const userData = pulumi.all([
  cloudInitScript,
  wireguardClientPublicKey,
]).apply(([script, publicKey]) => {
  return script.content
    .replace(/WIREGUARD_SERVER_PRIVATE_IP_PLACEHOLDER/g, wireguardServerPrivateIp)
    .replace(/WIREGUARD_SERVER_LISTEN_PORT_PLACEHOLDER/g, wireguardServerListenPort.toString())
    .replace(/WIREGUARD_CLIENT_PUBLIC_KEY_PLACEHOLDER/g, publicKey)
    .replace(/WIREGUARD_CLIENT_IP_PLACEHOLDER/g, wireguardClientIp);
});

// 4. Generate SSH key pair for the server
const keypair = forge.pki.rsa.generateKeyPair(4096);
const sshPrivateKey = forge.ssh.privateKeyToOpenSSH(keypair.privateKey);
const sshPublicKey = forge.ssh.publicKeyToOpenSSH(keypair.publicKey, "root@host");

const sshKey = new hcloud.SshKey("wireguard-ssh-key", {
  name: `wireguard-server-ssh-key`,
  publicKey: sshPublicKey,
});

// 5. Provision the Hetzner Cloud Server
const wireguardServer = new hcloud.Server("wireguard-server", {
  name: "wireguard-server",
  sshKeys: [sshKey.id],
  serverType: "cx22",
  image: "debian-12",
  location: serverRegion,
  userData: userData
});

// 6. Define the Hetzner Cloud Firewall
const wireguardFirewall = new hcloud.Firewall("wireguard-firewall", {
  name: "wireguard-server-firewall",
  rules: [
    {
      direction: "in",
      protocol: "tcp",
      port: "22",
      sourceIps: ["0.0.0.0/0", "::/0"],
      description: "Allow SSH access",
    },
    {
      direction: "in",
      protocol: "udp",
      port: wireguardServerListenPort.toString(),
      sourceIps: ["0.0.0.0/0", "::/0"],
      description: "Allow WireGuard VPN traffic",
    }
  ]
});

// 7. Attach the Firewall to the Server
new hcloud.FirewallAttachment("wireguard-firewall-attachment", {
  firewallId: wireguardFirewall.id.apply(id => parseInt(id, 10)),
  serverIds: [wireguardServer.id.apply(id => parseInt(id, 10))],
});

// 8. Wait for cloud-init to complete
const waitForCloudInit = new command.remote.Command("wait-for-cloud-init", {
  connection: {
    host: wireguardServer.ipv4Address,
    user: "root",
    privateKey: sshPrivateKey,
  },
  create: "cloud-init status --wait || true",
}, {dependsOn: [wireguardServer]});


// 9. Export relevant outputs
export const serverPublicIp = wireguardServer.ipv4Address;

const getWireguardPublicKey = new command.remote.Command("get-wireguard-public-key", {
  connection: {
    host: wireguardServer.ipv4Address,
    user: "root",
    privateKey: sshPrivateKey
  },
  create: "cat /etc/wireguard/server_public.key",
}, {dependsOn: [waitForCloudInit]});

export const wireguardServerPublicKey = getWireguardPublicKey.stdout;


