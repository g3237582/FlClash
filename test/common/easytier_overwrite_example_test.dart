import 'dart:io';

import 'package:flutter_test/flutter_test.dart';
import 'package:yaml/yaml.dart';

void main() {
  final script = File('examples/easytier/easytier-tailscale-overwrite.js');
  final yaml = File('examples/easytier/easytier-tailscale.yaml');

  test('overwrite examples stay placeholders and keep EasyTier primary', () {
    expect(script.existsSync(), isTrue);
    expect(yaml.existsSync(), isTrue);

    final js = script.readAsStringSync();
    expect(js, contains("type: 'easytier'"));
    expect(js, contains("'ffi-library': './libeasytier_ffi.so'"));
    expect(js, contains("'instance-name': 'CHANGE_ME_INSTANCE'"));
    expect(js, contains('CHANGE_ME_NETWORK'));
    expect(js, contains('CHANGE_ME_SECRET'));
    expect(js, contains("type: 'tailscale'"));
    expect(js, contains('CHANGE_ME_TAILSCALE_AUTHKEY'));
    expect(js, contains('10.77.0.0/24,easytier'));
    expect(js, contains('100.64.0.0/10,tailscale'));
    expect(js, contains('10.0.0.0/8,DIRECT'));
    expect(js, isNot(contains('glq')));
    expect(js, isNot(contains('tskey-auth-')));

    final parsed = loadYaml(yaml.readAsStringSync()) as YamlMap;
    final proxies = parsed['proxies'] as YamlList;
    expect(proxies[0]['type'], 'easytier');
    expect(proxies[0]['ffi-library'], './libeasytier_ffi.so');
    expect(proxies[0]['udp'], isTrue);
    expect(proxies[1]['type'], 'tailscale');
    final text = yaml.readAsStringSync();
    expect(text, contains('CHANGE_ME_SECRET'));
    expect(text, contains('no_tun = true'));
    expect(
      File('android/easytier/README.md').readAsStringSync(),
      contains('WireGuard Portal'),
    );
  });
}
