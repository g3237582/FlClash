import 'dart:io';

import 'package:path/path.dart' as p;
import 'package:setup_hooks/src/easytier_overlay.dart';
import 'package:test/test.dart';

void main() {
  test('applyEasyTierOverlay copies files except SOURCE.txt', () {
    final root = Directory.systemTemp.createTempSync('easytier_overlay_');
    addTearDown(() => root.deleteSync(recursive: true));

    final overlay = Directory(p.join(root.path, easyTierOverlayDir))
      ..createSync(recursive: true);
    File(p.join(overlay.path, 'SOURCE.txt')).writeAsStringSync('meta');
    File(
      p.join(overlay.path, 'adapter', 'outbound', 'easytier.go'),
    ).createSync(recursive: true);
    File(
      p.join(overlay.path, 'adapter', 'outbound', 'easytier.go'),
    ).writeAsStringSync('package outbound\n');

    final dest = Directory(p.join(root.path, 'core', 'Clash.Meta'))
      ..createSync(recursive: true);
    applyEasyTierOverlay(rootDir: root.path, coreDir: 'core');

    final copied = File(
      p.join(dest.path, 'adapter', 'outbound', 'easytier.go'),
    );
    expect(copied.readAsStringSync(), 'package outbound\n');
    expect(File(p.join(dest.path, 'SOURCE.txt')).existsSync(), isFalse);
  });

  test('applyEasyTierOverlay is a no-op without an overlay directory', () {
    final root = Directory.systemTemp.createTempSync('easytier_overlay_empty_');
    addTearDown(() => root.deleteSync(recursive: true));
    applyEasyTierOverlay(rootDir: root.path, coreDir: 'core');
  });
}
