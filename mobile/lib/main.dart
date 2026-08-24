import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'screens/title_screen.dart';
import 'services/audio_service.dart';
import 'services/api_service.dart';
import 'services/session_service.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  const apiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:8080');
  final api = ApiService(baseUrl: apiBaseUrl);
  final audio = AudioService();
  await audio.initialize();
  unawaited(api.initialize());
  unawaited(audio.resumeMusic());
  runApp(DetectiveGameApp(api: api, audio: audio));
}

class DetectiveGameApp extends StatelessWidget {
  final ApiService api;
  final AudioService audio;

  DetectiveGameApp({super.key, required this.api, AudioService? audio}) : audio = audio ?? AudioService();

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: api),
        Provider.value(value: audio),
        ChangeNotifierProvider(create: (_) => SessionService(api)),
      ],
      child: MaterialApp(
        title: 'ДетектИИв',
        debugShowCheckedModeBanner: false,
        // Android renders Flutter apps edge-to-edge. Keep every route and
        // modal sheet above the system navigation area; AppBar still handles
        // the top system inset itself.
        builder: (context, child) => SafeArea(top: false, child: child ?? const SizedBox.shrink()),
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.amber, brightness: Brightness.dark),
          useMaterial3: true,
        ),
        home: const TitleScreen(),
      ),
    );
  }
}
