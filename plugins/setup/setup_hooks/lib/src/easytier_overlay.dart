import 'dart:io';

import 'package:logging/logging.dart';
import 'package:path/path.dart' as p;

import 'fingerprint.dart';
import 'util.dart';

final _log = Logger('easytier_overlay');

const easyTierOverlayDir = 'core/easytier_overlay';

void applyEasyTierOverlay({required String rootDir, required String coreDir}) {
  final overlay = Directory(p.join(rootDir, easyTierOverlayDir));
  if (!overlay.existsSync()) {
    return;
  }
  final destRoot = p.join(rootDir, coreDir, 'Clash.Meta');
  if (!Directory(destRoot).existsSync()) {
    throw StateError(
      'Clash.Meta is missing at $destRoot; run '
      'git submodule update --init --recursive',
    );
  }
  var count = 0;
  for (final entity in overlay.listSync(recursive: true)) {
    if (entity is! File) continue;
    if (p.basename(entity.path) == 'SOURCE.txt') continue;
    final rel = p.relative(entity.path, from: overlay.path);
    copyFile(entity.path, p.join(destRoot, rel));
    count++;
  }
  _log.info('Applied EasyTier overlay ($count files) to $destRoot');
}

List<String> easyTierOverlayInputs(String rootDir) {
  final overlay = Directory(p.join(rootDir, easyTierOverlayDir));
  if (!overlay.existsSync()) {
    return const [];
  }
  return collectFiles(overlay.path);
}
