/**
 * FlClash overwrite script: EasyTier FFI (primary) + Tailscale.
 * Replace every CHANGE_ME_* placeholder. Do not commit real mesh secrets.
 *
 * ffi-library is the copy FlClash stages into the Core home directory:
 *   {filesDir}/libeasytier_ffi.so
 */

function main(config) {
  const proxies = Array.isArray(config.proxies) ? config.proxies.slice() : [];
  const rules = Array.isArray(config.rules) ? config.rules.slice() : [];

  upsertProxy(proxies, {
    name: 'easytier',
    type: 'easytier',
    'ffi-library': './libeasytier_ffi.so',
    'instance-name': 'CHANGE_ME_INSTANCE',
    config: [
      'instance_name = "CHANGE_ME_INSTANCE"',
      'ipv4 = "10.77.0.10"',
      '',
      '[network_identity]',
      'network_name = "CHANGE_ME_NETWORK"',
      'network_secret = "CHANGE_ME_SECRET"',
      '',
      '[[peer]]',
      'uri = "tcp://CHANGE_ME_PEER_HOST:11010"',
      '',
      '[flags]',
      'no_tun = true',
    ].join('\n'),
    udp: true,
  });

  upsertProxy(proxies, {
    name: 'tailscale',
    type: 'tailscale',
    hostname: 'CHANGE_ME_TAILSCALE_HOSTNAME',
    'auth-key': 'CHANGE_ME_TAILSCALE_AUTHKEY',
    udp: true,
  });

  const prefix = [
    'IP-CIDR,10.77.0.0/24,easytier,no-resolve',
    'IP-CIDR,100.64.0.0/10,tailscale,no-resolve',
    'IP-CIDR,10.0.0.0/8,DIRECT,no-resolve',
    'IP-CIDR,172.16.0.0/12,DIRECT,no-resolve',
    'IP-CIDR,192.168.0.0/16,DIRECT,no-resolve',
    'IP-CIDR,127.0.0.0/8,DIRECT,no-resolve',
  ];

  config.proxies = proxies;
  config.rules = prefix.concat(rules.filter((rule) => !isOverlayOrLanRule(rule)));
  return config;
}

function upsertProxy(proxies, proxy) {
  const index = proxies.findIndex((item) => item && item.name === proxy.name);
  if (index >= 0) {
    proxies[index] = proxy;
    return;
  }
  proxies.push(proxy);
}

function isOverlayOrLanRule(rule) {
  if (typeof rule !== 'string') {
    return false;
  }
  return (
    rule.includes('10.77.0.0/24') ||
    rule.includes('100.64.0.0/10') ||
    rule.includes('10.0.0.0/8,DIRECT') ||
    rule.includes('172.16.0.0/12,DIRECT') ||
    rule.includes('192.168.0.0/16,DIRECT') ||
    rule.includes('127.0.0.0/8,DIRECT')
  );
}
